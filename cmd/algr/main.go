package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	"net/http"

	"github.com/34South/envr"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)

var api string
var apiAuth string
var apiMembers string
var membersIndex string
var apiResources string
var resourcesIndex string
var apiModules string
var modulesIndex string

var maxBatch int

var token string

type TestRequest struct {
	Status int
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

type Docs struct {
	Index string
	Data  []map[string]interface{}
}

// flags
var collections = flag.String("c", "none", "collections to sync - 'none' for check, 'all' for everything or 'only' followed by a collection name")
var backdays = flag.Int("b", 1, "how many days back to check for updated records")

// backDate is a string date in format "2017-01-21T13:35:30+10:00" (RFC3339) so we can pass it to the API
var backDate string

func init() {
	envr.New("algrEnv", []string{
		"MAPPCPD_ADMIN_USER",
		"MAPPCPD_ADMIN_PASS",
		"MAPPCPD_ALGOLIA_APP_ID",
		"MAPPCPD_ALGOLIA_API_KEY",
		"MAPPCPD_ALGOLIA_BATCH_SIZE",
		"MAPPCPD_API_URL",
		"MAPPCPD_ALGOLIA_MEMBERS_INDEX",
		"MAPPCPD_ALGOLIA_MODULES_INDEX",
		"MAPPCPD_ALGOLIA_RESOURCES_INDEX",
	}).Auto()

	api = os.Getenv("MAPPCPD_API_URL")
	apiAuth = api + "/v1/auth/admin"
	apiMembers = api + "/v1/a/members"
	membersIndex = os.Getenv("MAPPCPD_ALGOLIA_MEMBERS_INDEX")
	apiResources = api + "/v1/a/resources"
	resourcesIndex = os.Getenv("MAPPCPD_ALGOLIA_RESOURCES_INDEX")
	apiModules = api + "/v1/a/modules"
	modulesIndex = os.Getenv("MAPPCPD_ALGOLIA_MODULES_INDEX")

	// Don't shadow maxBatch!
	var err error
	maxBatch, err = strconv.Atoi(os.Getenv("MAPPCPD_ALGOLIA_BATCH_SIZE"))
	if err != nil {
		log.Fatalln("Could not set batch size:", err)
	}
}

func main() {

	// set backDate from -b flag
	flag.Parse()
	t := time.Now()
	backDate = t.AddDate(0, 0, -(*backdays)).Format(time.RFC3339)

	switch *collections {
	case "none":
		test()
		auth()
	case "all":
		test()
		auth()
		indexMembers()
		indexResources()
		indexModules()
	case "only":
		test()
		auth()
		only()
	default:
		fmt.Println("Not sure what to do, specify -c 'none', 'all' or 'only'")
	}
}

func test() {

	fmt.Print("Test connection to API... ")

	hc := &http.Client{Timeout: 90 * time.Second}
	r, err := hc.Get(api)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()

	fmt.Println(r.Status)
}

func auth() {

	fmt.Print("Authenticating with API... ")

	a := AuthRequest{}
	b := `{"login": "` + os.Getenv("MAPPCPD_ADMIN_USER") + `","password": "` + os.Getenv("MAPPCPD_ADMIN_PASS") + `"}`
	hc := &http.Client{Timeout: 90 * time.Second}
	res, err := hc.Post(apiAuth, "application/json", strings.NewReader(b))
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&a)
	token = a.Data.Token
	fmt.Println("ok")
}

// only runs algr on just one of the collections
func only() {

	args := flag.Args()
	var col string
	if len(args) > 0 {
		col = strings.ToLower(args[0])
	}

	switch col {
	case "members":
		indexMembers()
	case "resources":
		indexResources()
	case "modules":
		indexModules()
	default:
		log.Fatalln("'only' flag required a valid collection name")
	}
}

func indexMembers() {

	fmt.Println("Indexing member records... ")

	if strings.ToLower(membersIndex) == "off" {
		fmt.Println("... member index is set to 'OFF' - nothing to do")
		return
	}

	fmt.Println("... fetching member docs updated since", backDate)
	xm := Docs{
		Index: membersIndex,
	}
	// Only members with a membership record
	// mongo shell query is: db.Members.find({"memberships.title": {$exists : true}})
	q := `{ "query": { "memberships.title": {"$exists": true}, "contact.directory": true, "updatedAt": {"$gte": "` + backDate + `"} }}`
	fetchDocs(apiMembers, q, &xm)

	// reshape the Resources Docs for algolia
	fmt.Println("... reshaping member docs for Algolia index")
	xm.Data = reshapeMembers(xm.Data)

	fmt.Println("... updating Algolia index:", os.Getenv("MAPPCPD_ALGOLIA_MEMBERS_INDEX"))
	indexDocs(&xm)
}

// indexResources manages the resources index. Note the Mongo collection now has a bool field called 'active'
// which mirrors the flag (0/1) field used in MySQL for soft deletes. So any record with active=false should be
// removed from the index.
func indexResources() {

	fmt.Println("Indexing resource records... ")

	if strings.ToLower(resourcesIndex) == "off" {
		fmt.Println("... resource index is set to OFF - nothing to do")
		return
	}

	fmt.Println("... fetching resource docs updated since", backDate)
	xr := Docs{
		Index: resourcesIndex,
	}
	q := `{"find": {"primary": true, "updatedAt": {"$gte": "` + backDate + `"}}}`
	fetchDocs(apiResources, q, &xr)

	fmt.Println("... reshaping resource docs for Algolia index")
	xr.Data = reshapeResources(xr.Data)

	fmt.Println("... updating Algolia index:", os.Getenv("MAPPCPD_ALGOLIA_RESOURCES_INDEX"))
	indexDocs(&xr)

	// Remove inactive resources from index
	fmt.Println("... removing inactive resources")
	q = `{"find": {"active": false}}`
	fetchDocs(apiResources, q, &xr)
	var objectIDs []string
	for _, v := range xr.Data {
		objectIDs = append(objectIDs, v["_id"].(string))
	}
	if err := deleteObjects(objectIDs, xr.Index); err != nil {
		fmt.Println("... error deleting resource objects -", err)
	}
}

func indexModules() {

	fmt.Println("Indexing module records... ")

	if strings.ToLower(modulesIndex) == "off" {
		fmt.Println("... modules index is set to 'OFF' - nothing to do")
		return
	}

	fmt.Println("... fetching module docs updated since", backDate)
	xm := Docs{
		Index: modulesIndex,
	}
	q := `{"find": {"current": true, "updatedAt": {"$gte": "` + backDate + `"}}}`
	fetchDocs(apiModules, q, &xm)

	fmt.Println("... updating Algolia index:", os.Getenv("MAPPCPD_ALGOLIA_MODULES_INDEX"))
	indexDocs(&xm)

}

// fetchDocs does a request to the API for doc records
func fetchDocs(api string, query string, docs *Docs) {

	// Set up the request
	hc := &http.Client{Timeout: 90 * time.Second}
	req, err := http.NewRequest("POST", api, strings.NewReader(query))
	if err != nil {
		log.Fatalln(err)
	}
	// Add auth header
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := hc.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(docs)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("... have", len(docs.Data), "docs to add to index:", docs.Index)
}

func indexDocs(docs *Docs) {

	batchSize := maxBatch

	for i := 0; i < len(docs.Data); i++ {

		// Reset this or we accumulate!!!
		objects := []algoliasearch.Object{}

		// If remaining items is less than batch size...
		if len(docs.Data)-i < batchSize {
			batchSize = len(docs.Data) - i
		}

		fmt.Println("--- next batch of", batchSize)
		for c := 0; c < batchSize; c++ {
			// set algolia objectID to _id
			docs.Data[i]["objectID"] = docs.Data[i]["_id"]
			objects = append(objects, docs.Data[i])
			// Don't increment i on the last loop because the outer loop
			// also increments it so we end up missing a value
			if c < (batchSize - 1) {
				i++
			}
		}

		indexBatch(objects, docs.Index)
	}
}

// indexBatch adds a batch of items to an index. Note that this seems to time out at somewhere over 1000 objects
// regardless of batch size, timeout settings or anything else.
func indexBatch(xo []algoliasearch.Object, indexName string) {

	client := algoliasearch.NewClient(os.Getenv("MAPPCPD_ALGOLIA_APP_ID"), os.Getenv("MAPPCPD_ALGOLIA_API_KEY"))
	index := client.InitIndex(indexName)
	batch, err := index.AddObjects(xo)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Algolia taskID:", batch.TaskID, "completed indexing of", len(batch.ObjectIDs), "objects for index:", indexName)
}

// deleteObjects removes objects from an algolia index, identified by an array of ObjectID strings
func deleteObjects(objectIDs []string, indexName string) error {
	client := algoliasearch.NewClient(os.Getenv("MAPPCPD_ALGOLIA_APP_ID"), os.Getenv("MAPPCPD_ALGOLIA_API_KEY"))
	index := client.InitIndex(indexName)
	batch, err := index.DeleteObjects(objectIDs)
	if err != nil {
		return err
	}
	fmt.Println("Algolia taskID:", batch.TaskID, "completed removal of", len(batch.ObjectIDs), "objects for index:", indexName)
	return nil
}

// reshapeResources modifies the resource values into a more suitable format for the Algolia index
func reshapeResources(data []map[string]interface{}) []map[string]interface{} {

	var d []map[string]interface{}

	for _, v := range data {

		pubDate := v["pubDate"].(map[string]interface{})
		publishedAt := pubDate["date"]
		publishedAtTS := timeStampFromDate(pubDate["date"].(string))

		//timeStampFromDate(v["pubDate"]["date"])
		r := map[string]interface{}{
			"_id":                  v["_id"],
			"id":                   v["id"],
			"createdAt":            v["createdAt"],
			"updatedAt":            v["updatedAt"],
			"publishedAt":          publishedAt,
			"publishedAtTimestamp": publishedAtTS,
			"type":                 v["type"],
			"name":                 v["name"],
			"description":          v["description"],
			"keywords":             v["keywords"],
			"shortUrl":             v["shortUrl"],
			"resourceUrl":          v["resourceUrl"],
		}

		// Pubmed Attributes
		attributes := v["attributes"].(map[string]interface{})

		v, ok := attributes["sourceId"]
		if ok {
			r["sourceId"] = v
		}

		v, ok = attributes["sourceName"]
		if ok {
			r["sourceName"] = v
		}

		v, ok = attributes["sourceNameAbbrev"]
		if ok {
			r["sourceNameAbbrev"] = v
		}

		v, ok = attributes["sourcePubDate"]
		if ok {
			r["sourcePubDate"] = v
		}

		v, ok = attributes["sourceVolume"]
		if ok {
			r["sourceVolume"] = v
		}

		v, ok = attributes["sourceIssue"]
		if ok {
			r["sourceIssue"] = v
		}

		v, ok = attributes["sourcePages"]
		if ok {
			r["sourcePages"] = v
		}

		d = append(d, r)
	}

	return d
}

// reshapeMembers modifies the members values into a more suitable format for the Algolia index
func reshapeMembers(data []map[string]interface{}) []map[string]interface{} {

	var d []map[string]interface{}

	for _, dv := range data {

		// concat name fields
		name := fmt.Sprintf("%s %s %s", dv["title"], dv["firstName"], dv["lastName"])

		// personal contact details
		contact := dv["contact"].(map[string]interface{})
		email := contact["emailPrimary"]
		mobile := contact["mobile"]

		// only use location info from the directory contact record, and only the general locality
		var location string
		xl, ok := contact["locations"].([]interface{})
		if ok {
			for _, v := range xl {
				l := v.(map[string]interface{})
				if l["type"] == "Directory" {
					location = fmt.Sprintf("%s %s %s", l["city"], l["state"], l["postcode"])
				}
			}
		}

		// Membership title - dig into the memberships array even though there is only one for now.
		var memberships []string
		xm, ok := dv["memberships"].([]interface{})
		if ok {
			for _, v := range xm {
				m := v.(map[string]interface{})
				memberships = append(memberships, fmt.Sprintf("%s %s", m["orgCode"], m["title"]))
			}
		}

		// Specialities
		var specialities []string
		xs, ok := dv["specialities"].([]interface{})
		if ok {
			for _, v := range xs {
				s := v.(map[string]interface{})
				specialities = append(specialities, s["name"].(string))
			}
		}

		// Affiliations are positions / affiliations with certain groups. Only include positions
		// where the end date is not set, or is in the future
		var affiliations []string
		xa, ok := dv["positions"].([]interface{})
		if ok {
			for _, v := range xa {
				a := v.(map[string]interface{})
				endDate, err := time.Parse("2006-01-02", a["end"].(string))
				if err != nil || endDate.After(time.Now()) {
					affiliations = append(affiliations, a["orgName"].(string))
				}
			}
		}

		r := map[string]interface{}{
			"_id":          dv["_id"],
			"id":           dv["id"],
			"name":         name,
			"email":        email,
			"mobile":       mobile,
			"location":     location,
			"membership":   memberships,
			"affiliations": affiliations,
			"specialities": specialities,
		}

		d = append(d, r)
	}

	fmt.Println(d)

	return d
}

func timeStampFromDate(date string) int64 {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		fmt.Println("Error parsing date string", err)
	}
	return t.Unix()
}
