package main

import (
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

// backDate is a string date in format "2017-01-21T13:35:30+10:00" (RFC3339) so we can pass it to the API
var backDate string

func init() {
	envr.New("algrEnv", []string{
		"MAPPCPD_ADMIN_USER",
		"MAPPCPD_ADMIN_PASS",
		"MAPPCPD_ALGOLIA_APP_ID",
		"MAPPCPD_ALGOLIA_API_KEY",
		"MAPPCPD_ALGOLIA_BACK_DAYS",
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

	// set backDate from BACK_DAYS
	t := time.Now()
	d, err := strconv.Atoi(os.Getenv("MAPPCPD_ALGOLIA_BACK_DAYS"))
	if err != nil {
		log.Fatalln(err)
	}
	backDate = t.AddDate(0, 0, -d).Format(time.RFC3339)
}

func main() {
	log.Println("Running algr...")
	log.Println("Test connection to API...")
	test()
	log.Println("Authenticating...")
	auth()
	indexMembers()
	indexResources()
	indexModules()
}

func test() {

	t := TestRequest{}
	hc := &http.Client{Timeout: 90 * time.Second}
	r, err := hc.Get(api)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&t)
	log.Println("Response:", t.Status)
}

func auth() {

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
}

func indexMembers() {

	if strings.ToLower(membersIndex) == "off" {
		fmt.Println("Member index is set to OFF... nothing to do")
		return
	}

	fmt.Println("Fetch Member Docs updated since", backDate)
	xm := Docs{
		Index: membersIndex,
	}
	// Only members with a membership record
	// mongo shell query is: db.Members.find({"memberships.title": {$exists : true}})
	q := `{ "query": { "memberships.title": {"$exists": true}, "updatedAt": {"$gte": "` + backDate + `"} }}`
	fetchDocs(apiMembers, q, &xm)
	fmt.Println("Index member docs...")
	indexDocs(&xm)
}

func indexResources() {

	if strings.ToLower(resourcesIndex) == "off" {
		fmt.Println("Resource index is set to OFF... nothing to do")
		return
	}

	fmt.Println("Fetch Resource Docs updated since", backDate)
	xr := Docs{
		Index: resourcesIndex,
	}
	q := `{"find": {"primary": true, "updatedAt": {"$gte": "` + backDate + `"}}}`
	fetchDocs(apiResources, q, &xr)
	fmt.Println("Index resource docs...")
	indexDocs(&xr)
}

func indexModules() {

	if strings.ToLower(modulesIndex) == "off" {
		fmt.Println("Modules index is set to OFF... nothing to do")
		return
	}

	fmt.Println("Fetch Module Docs updated since", backDate)
	xm := Docs{
		Index: modulesIndex,
	}
	q := `{"find": {"current": true, "updatedAt": {"$gte": "` + backDate + `"}}}`
	fmt.Println(q)
	fetchDocs(apiModules, q, &xm)
	fmt.Println("Index module docs...")
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
	fmt.Println("Have ", len(docs.Data), "docs to add to index:", docs.Index)
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
