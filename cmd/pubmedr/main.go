package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"log"

	"io/ioutil"
	"net/http"
	"encoding/json"
	"encoding/xml"

	"github.com/34South/envr"
	"github.com/pkg/errors"
)

// Pointer to an http.client
var httpClient = &http.Client{Timeout: 90 * time.Second}

// Pubmed api endpoints
const pubmedSearch = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&retmode=json&retmax=%v&retstart=%v" +
	"&reldate=%v&datetype=pdat&term=%v"
const pubMedFetch = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&retmode=xml&rettype=abstract&id="

// The id of the resource type (ol_resource_type table) for journal articles
const resourceTypeID = 80

var api, apiAuth, apiResource string

type PubMedSearch struct {
	Header map[string]string  `json:"header"`
	Result PubMedSearchResult `json:"esearchresult"`
}
type PubMedSearchResult struct {
	Count string   `json:"count"`
	IDs   []string `json:"idlist"`
}

// pubmedQuery type represents a single pubmed query, ie one of the objects read in from pubmed.json
type pubmedQuery struct {
	Run            bool                   `json:"run"`
	Category       string                 `json:"category"`
	SearchTerm     string                 `json:"searchTerm"`
	RelDate        int                    `json:"relDate"`
	Attributes     map[string]interface{} `json:"attributes"`
	ResourceTypeID int                    `json:"resourceTypeId"`
}

// PubMedSummary Result - each result is indexed by the id of the record requested, even if there is only one.
// As we can pass multiple ids on the URL to save requests eg 521345,765663,121234,3124256
// use XMl to fetch this as easier to get to nested data
type PubMedArticleSet struct {
	Articles []PubMedArticle `xml:"PubmedArticle"`
}
type PubMedArticle struct {
	ID              int            `xml:"MedlineCitation>PMID"`
	Title           string         `xml:"MedlineCitation>Article>ArticleTitle"`
	Abstract        []AbstractText `xml:"MedlineCitation>Article>Abstract>AbstractText"`
	PublishYear     string         `xml:"PubmedData>History>PubMedPubDate>Year"`
	PublishMonth    string         `xml:"PubmedData>History>PubMedPubDate>Month"`
	PublishDay      string         `xml:"PubmedData>History>PubMedPubDate>Day"`
	ArticleIDList   []ArticleID    `xml:"PubmedData>ArticleIdList>ArticleId"`
	KeywordList     []string       `xml:"MedlineCitation>KeywordList>Keyword"`
	MeshHeadingList []string       `xml:"MedlineCitation>MeshHeadingList>MeshHeading>DescriptorName"`
	AuthorList      []Author       `xml:"MedlineCitation>Article>AuthorList>Author"`
}
type AbstractText struct {
	Key   string `xml:"label,attr"`
	Value string `xml:",chardata"`
}
type ArticleID struct {
	Key   string `xml:"IdType,attr"`
	Value string `xml:",chardata"`
}
type Author struct {
	Key      string `xml:"ValidYN,attr"`
	LastName string `xml:"LastName"`
	Initials string `xml:"Initials"`
}

type Resource struct {
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
	PubDate      PubDate                `json:"pubDate"`
	TypeID       int                    `json:"typeId"`
	Primary      bool                   `json:"primary"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Keywords     []string               `json:"keywords"`
	ResourceURL  string                 `json:"resourceUrl"`
	ShortURL     string                 `json:"shortUrl"`
	ThumbnailURL string                 `json:"thumbnailUrl"`
	Attributes   map[string]interface{} `json:"attributes"`
}

type PubDate struct {
	Date  time.Time `json:"date" bson:"date"`
	Year  int       `json:"year" bson:"year"`
	Month int       `json:"month" bson:"month"`
	Day   int       `json:"day" bson:"day"`
}

type AuthRequest struct {
	Status  int
	Message string
	Data    AuthData
}
type AuthData struct {
	Token     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// Universal token for accessing API
var token string

// Batch size - ie how many to process at a time
var batchSize int

// init the env vars
func init() {
	envr.New("algrEnv", []string{
		"MAPPCPD_ADMIN_USER",
		"MAPPCPD_ADMIN_PASS",
		"MAPPCPD_API_URL",
		"MAPPCPD_PUBMED_RETMAX",
		"MAPPCPD_PUBMED_BATCH_FILE",
	}).Auto()

	// Initialise the batchSize based on PUBMED_RETMAX
	var err error
	batchSize, err = strconv.Atoi(os.Getenv("MAPPCPD_PUBMED_RETMAX"))
	if err != nil {
		log.Fatalln(err)
	}

	// set api strings
	api = os.Getenv("MAPPCPD_API_URL")
	apiAuth = api + "/v1/auth/admin"
	apiResource = api + "/v1/a/batch/resources"

}

func main() {

	// Read in the batch file and set up the jobs
	pubmedBatch, err := setBatch()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Make sure the API is ok...
	testAPI()
	authAPI()

	// run the jobs...
	fmt.Println("Pubmed Jobs:")
	for i, v := range pubmedBatch {
		fmt.Println("\n####################################################################")
		fmt.Println("Job ", i, "- Category:", v.Category)
		if v.Run == false {
			fmt.Println("Run = false ...skipping job")
			continue
		} else {
			fmt.Printf("Doing pubmed search going back %v days...\n", v.RelDate)
			c := pubmedCount(v.SearchTerm, v.RelDate, 0)
			fmt.Printf("Set size: %v\n", c)
			for i := 0; i < c; i += batchSize {
				// new set each iteration
				pmSearch := PubMedSearch{}
				pmArticles := PubMedArticleSet{}

				fmt.Println("Fetch batch of IDs, starting at", i)
				pmSearch.fetchIDs(v.SearchTerm, v.RelDate, i)
				fmt.Println("Fetch the article summaries for this batch of ids...")
				fmt.Println(pmSearch.Result.IDs)
				pmSearch.getSummaries(&pmArticles)
				//pmArticles.inspect()
				// pass the standard attributes for this set as they get saved to the db as well...
				pmArticles.indexSummaries(v.Attributes)
			}
		}
	}
}

// setBatch reads the json batch file and sets up the batch of jobs stored in pubmedBatch
func setBatch() ([]pubmedQuery, error) {

	fmt.Println("setBatch()...")
	f := os.Getenv("MAPPCPD_PUBMED_BATCH_FILE")
	pq := []pubmedQuery{}

	// Decide if it is a local file or a url...
	if strings.Contains(f, "http") {
		fmt.Println("Fetching batch file at url: ", f)
		readURLBatchFile(f, &pq)
	} else {
		fmt.Println("Reading local batch file at: ", f)
		readLocalBatchFile(f, &pq)
	}

	return pq, nil
}

// readURLBatchFile sets the pubmed queries from a local json file
func readURLBatchFile(f string, pq *[]pubmedQuery) error {

	res, err := http.Get(f)
	if err != nil {
		msg := fmt.Sprintf("readURLBatchFile could not fetch json file at %v - %v", f, err)
		return errors.New(msg)
	}

	xb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		msg := fmt.Sprintf("readURLBatchFile could not read response body %v", err)
		return errors.New(msg)
	}

	err = json.Unmarshal(xb, &pq)
	if err != nil {
		msg := fmt.Sprintf("Error Unmarshaling into pubmedBatch - %v", err)
		return errors.New(msg)
	}

	return nil
}

// readLocalBatchFile sets the pubmed queries from a local json file
func readLocalBatchFile(f string, pq *[]pubmedQuery) error {

	xb, err := ioutil.ReadFile(f)
	if err != nil {
		msg := fmt.Sprintf("Error reading batch file: %v - %v", f, err)
		return errors.New(msg)
	}
	err = json.Unmarshal(xb, &pq)
	if err != nil {
		msg := fmt.Sprintf("Error Unmarshaling into pubmedBatch - %v", err)
		return errors.New(msg)
	}

	return nil
}

// testAPI pings to API to make sure we have a connection
func testAPI() {
	fmt.Print("Test API connection... ")
	fmt.Println(api)
	res, err := httpClient.Get(api)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatalln(res.Status)
	}

	fmt.Printf("%v\n", res.Status)
}

func authAPI() {
	fmt.Print("Authenticate and get token... ")
	a := AuthRequest{}
	b := `{"login": "` + os.Getenv("MAPPCPD_ADMIN_USER") + `","password": "` + os.Getenv("MAPPCPD_ADMIN_PASS") + `"}`
	res, err := httpClient.Post(apiAuth, "application/json", strings.NewReader(b))
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&a)
	token = a.Data.Token
	fmt.Println("ok")
}

// pubmedCount runs the pubmed query with rettype=count to get the number or articles
func pubmedCount(searchTerm string, relDate, startAt int) int {

	var c = struct {
		Result map[string]string `json:"esearchresult"`
	}{}

	url := fmt.Sprintf(pubmedSearch,
		os.Getenv("MAPPCPD_PUBMED_RETMAX"),
		startAt,
		relDate,
		searchTerm) + "&rettype=count"
	fmt.Println(url)
	r, err := httpClient.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		log.Fatalln(err)
	}
	count, err := strconv.Atoi(c.Result["count"])
	if err != nil {
		log.Fatalln(err)
	}

	return count
}

// fetchIDs queries the pubmed api and returns a maximum number of ids specified in env var PUBMED_RETMAX, based on the
// search term specified in PUBMED_TERM
func (ps *PubMedSearch) fetchIDs(searchTerm string, relDate, startAt int) {
	// Add the return max and search term to the api url...
	url := fmt.Sprintf(pubmedSearch,
		os.Getenv("MAPPCPD_PUBMED_RETMAX"),
		startAt,
		relDate,
		searchTerm)

	r, err := httpClient.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()

	//xb, _ := ioutil.ReadAll(r.Body)
	//fmt.Println(string(xb))

	err = json.NewDecoder(r.Body).Decode(&ps)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Here's the pubmed search url:")
	fmt.Println("########################################")
	fmt.Println(url)
	fmt.Println("########################################")
	fmt.Println("Returning max", os.Getenv("MAPPCPD_PUBMED_RETMAX"), "records from", ps.Result.Count, "total results")
}

// getSummaries fetches the article summary for each of its (pubmed) IDs and stores them in the PubMedArticleSet that is passed in.
// Fetching multiple articles means fewer calls to the api.
func (ps *PubMedSearch) getSummaries(pa *PubMedArticleSet) {

	idString := strings.Join(ps.Result.IDs, ",")
	// This might be too long for GET request so need to use POST!
	u := pubMedFetch + idString

	//emptyBody := bytes.Buffer{}
	r, err := httpClient.Get(u)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()

	//xb, _ := ioutil.ReadAll(r.Body)
	//fmt.Println(string(xb))

	err = xml.NewDecoder(r.Body).Decode(&pa)
	if err != nil {
		log.Fatalln(err)
	}
}

// indexSummaries adds the resources records to the MappCPD MySQL database via the API.
// These will subsequently be updated to MongoDB and then to Algolia indexes via mongr and algr services.
// Note that any double quotes INSIDE strings are replaced with single quotes so it will play with MySQL
// nicely... this seems better than escaping the double quotes!
func (pa PubMedArticleSet) indexSummaries(attributes map[string]interface{}) {

	// We post the articles as a batch - ie an array of article objects, so convert them to JSON here...
	var js string

	for i, v := range pa.Articles {

		var err error

		r := Resource{}
		r.CreatedAt = time.Now()
		r.UpdatedAt = time.Now()

		// Concat date string, then create time.Time value from the string format "2006-1-2"
		d := v.PublishYear + "-" + v.PublishMonth + "-" + v.PublishDay
		r.PubDate.Date, err = time.Parse("2006-1-2", d)
		if err != nil {
			fmt.Println(err)
		}
		r.PubDate.Year, err = strconv.Atoi(v.PublishYear)
		if err != nil {
			fmt.Println(err)
		}
		r.PubDate.Month, err = strconv.Atoi(v.PublishMonth)
		if err != nil {
			fmt.Println(err)
		}
		r.PubDate.Day, err = strconv.Atoi(v.PublishDay)
		if err != nil {
			fmt.Println(err)
		}
		r.Primary = true
		r.TypeID = resourceTypeID
		r.Name = strings.Replace(v.Title, `"`, `'`, -1)

		// Resource description will come from the <Abstract> node. This contains sub nodes, <AbstractText>
		// that may be of different types, distinguished by a "label" attribute with values like "BACKGROUND",
		// "METHODS", "RESULTS", "CONCLUSION", "CLINICAL TRIAL REGISTRATION". For now, take the first one
		// which is generally "BACKGROUND"... can get picky later on. It can also be empty!
		if len(v.Abstract) > 0 {
			r.Description = strings.Replace(v.Abstract[0].Value, `"`, `'`, -1)
		}

		// Do keywords
		// If we are using MEDLINE only we get MeshHeadings, otherwise KeyWords
		// So mash them together as one set of KeyWords...
		// Also, we will stick the Authors into the Keywords to assist with searches
		// and tack the PubmedID on the end as well - why not?
		v.KeywordList = append(v.KeywordList, v.MeshHeadingList...)
		// Concat Author LastName and Initials
		xa := []string{}
		for _, a := range v.AuthorList {
			n := a.LastName + " " + a.Initials
			xa = append(xa, n)
		}
		v.KeywordList = append(v.KeywordList, xa...)
		v.KeywordList = append(v.KeywordList, strconv.Itoa(v.ID))

		// Some of the keywords are actually phrases split with a comma like this:
		// Intubation, Intratracheal - ie that is a single mesh term
		// This stuffs up the unmarshaling at the API a bit so we will strip out the comma.
		// The only down side is a few repeated keywords
		for i, w := range v.KeywordList {
			v.KeywordList[i] = strings.TrimSpace(strings.Replace(w, ",", "", -1))
		}

		r.Keywords = v.KeywordList

		// Make DOI link
		for _, aid := range v.ArticleIDList {
			if aid.Key == "doi" {
				r.ResourceURL = fmt.Sprintf("https://doi.org/%s", aid.Value)
			}
		}
		// We cannot do short url as we don't have an ID for new resource record

		// Attributes is fixed (json string) for all pubmed articles in the same batch
		// This is saves to the db and used later to assist with search faceting
		r.Attributes = attributes

		j, err := json.Marshal(r)
		if err != nil {
			log.Fatalln(err)
		}

		js += string(j)

		// Add to data and add ,\n if we are not on last one
		if (i + 1) < len(pa.Articles) {
			js = js + ",\n"
		}
	}

	// data is now a string of JSON objects... just need the square brackets...
	js = `{"data": [` + js + `]}`
	//fmt.Println(js)

	// do it!
	addResources(js)
}

// inspect is a utility func to look at results and exit
func (pa *PubMedArticleSet) inspect() {

	fmt.Println("Inspecting Articles")

	// Range over the articles
	for i, a := range pa.Articles {
		fmt.Println("\npmArticles.Articles[", i, "]")
		fmt.Println("------------------------")
		fmt.Printf("PMID: %v\n", a.ID)
		fmt.Printf("Title: %s\n", a.Title)
		if len(a.Abstract) > 0 {
			fmt.Printf("Abstract: %s\n", a.Abstract)
		}

		fmt.Printf("Day: %s\n", a.PublishDay)
		fmt.Printf("Month: %s\n", a.PublishMonth)
		fmt.Printf("Year: %s\n", a.PublishYear)

		// Keywords...
		fmt.Printf("Keywords: %v\n", a.KeywordList)
		fmt.Printf("MeshHeadings: %v\n", a.MeshHeadingList)
		fmt.Printf("Authors: %s\n", a.AuthorList)

		xa := []string{}
		for _, v := range a.AuthorList {
			n := v.LastName + " " + v.Initials
			xa = append(xa, n)
		}
		a.KeywordList = append(a.KeywordList, a.MeshHeadingList...)
		a.KeywordList = append(a.KeywordList, xa...)
		a.KeywordList = append(a.KeywordList, strconv.Itoa(a.ID))

		// Some of the keywords are actually phrases split with a comma like this:
		// Intubation, Intratracheal - ie that is a single mesh term
		// This stuffs up the unmarshalling at the API a bit so we will strip out the comma.
		// The only down side is a few repeated keywords
		for i, v := range a.KeywordList {
			a.KeywordList[i] = strings.TrimSpace(strings.Replace(v, ",", "", -1))
		}

		fmt.Printf("All %v Keywords: %v\n", len(a.KeywordList), a.KeywordList)

		for _, aid := range a.ArticleIDList {
			if aid.Key == "doi" {
				fmt.Printf("Link: https://doi.org/%s\n", aid.Value)
			}
		}
	}
	os.Exit(1)
}

// addResources POSTs JSON-formatted resource record to MappCPD api
// j is fully formatted JSOn string {"data": [{...}, {...}]}
func addResources(j string) {

	fmt.Print("POST batch of resources to api...")
	req, err := http.NewRequest("POST", apiResource, strings.NewReader(j))
	if err != nil {
		log.Fatalln(err)
	}

	// Add headers
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln("Error posting http request", err)
	}
	defer res.Body.Close()

	// Response from API is JSON, so have a squiz...
	var jb = struct {
		Status  int
		Result  string
		Message string
		Meta    interface{}
		Data    interface{}
	}{}

	json.NewDecoder(res.Body).Decode(&jb)
	if jb.Result == "failed" {
		log.Fatalln(jb.Message)
	}
	//fmt.Println(jb)
}
