package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/34South/envr"
)

const waitGroupSize = 5

var httpClient = &http.Client{Timeout: 30 * time.Second}
var wg sync.WaitGroup

var api string
var apiAuth string
var apiMemberIDList string
var apiResourceIDList string
var apiModuleIDList string
var apiMembers string
var apiResources string
var apiModules string

var token string

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

type Member struct {
	Data map[string]interface{}
}

type Resource struct {
	Data interface{}
}

type Module struct {
	Data interface{}
}

// flags
var collections = flag.String("c", "none", "collections to sync - 'none' for check, 'all' for everything or 'only' followed by 'name:all' or 'name:id1,id2,id3'.\neg 'members:all' to sync member records, 'resources:1234' to sync resource with id=1234")
var backdays = flag.Int("b", 1, "how many days back to check for updated records")

// The filters string passed on the url
var urlQuery string

func init() {
	envr.New("mongrEnv", []string{
		"MAPPCPD_ADMIN_USER",
		"MAPPCPD_ADMIN_PASS",
		"MAPPCPD_API_URL",
	}).Auto()

	api = os.Getenv("MAPPCPD_API_URL")
	apiAuth = api + "/v1/auth/admin"
	apiMemberIDList = api + "/v1/a/idlist?t=member"
	apiResourceIDList = api + "/v1/a/idlist?t=ol_resource"
	apiModuleIDList = api + "/v1/a/idlist?t=ol_module"
	apiMembers = api + "/v1/a/members/"
	apiResources = api + "/v1/a/resources/"
	apiModules = api + "/v1/a/modules/"
}

const doSync = true

func main() {

	flag.Parse()

	switch *collections {
	case "none":
		none()
	case "all":
		all()
	case "only":
		only()
	default:
		fmt.Println("Not sure what to do, specify -c 'none', 'all' or 'only'")
	}
}

func none() {
	fmt.Println("No sync targets specified, will only test the api connection")
	test()
	auth()
}

func all() {
	fmt.Println("syncing all collections...")
	test()
	auth()
	filter()
	if doSync == true {
		syncData("members", "")
		syncData("modules", "")
		syncData("resources", "")
	}

}

func only() {
	fmt.Print("Selectively syncing ")
	test()
	auth()
	filter()

	// Need args to specify what to sync
	args := flag.Args()
	var col string
	if len(args) > 0 {
		col = args[0]
	}
	fmt.Println("Collection", col)
	var ids string
	if len(args) > 1 {
		ids = args[1]
	}

	switch col {
	case "members":
		fmt.Print(col)
		if doSync == true {
			syncData(col, ids)
		}
	case "modules":
		fmt.Print(col)
		if doSync == true {
			syncData(col, ids)
		}
	case "resources":
		fmt.Print(col)
		if doSync == true {
			syncData(col, ids)
		}
	default:
		fmt.Println(" ???? - need to specify 'members', 'modules' or 'resources'")
		os.Exit(1)
	}
}

func test() {

	fmt.Print("Test connection to API...")
	r, err := httpClient.Get(api)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()
	io.Copy(ioutil.Discard, r.Body)
	fmt.Println(r.Status)
	if r.StatusCode != 200 {
		fmt.Println("test() failed... exiting")
		os.Exit(1)
	}
}

func auth() {

	fmt.Print("Testing authentication... ")
	a := AuthRequest{}
	b := `{"login": "` + os.Getenv("MAPPCPD_ADMIN_USER") + `","password": "` + os.Getenv("MAPPCPD_ADMIN_PASS") + `"}`
	res, err := httpClient.Post(apiAuth, "application/json", strings.NewReader(b))
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&a)
	token = a.Data.Token
	if len(token) == 0 {
		fmt.Println("auth() failed... exiting")
		os.Exit(1)
	}
	fmt.Println("ok")
}

func filter() {

	// Time format we need here is MySQL timestamp
	t := time.Now().AddDate(0, 0, -*backdays).Format("2006-01-02 15:04:05")
	// URL encode the time string...
	ut := strings.Replace(t, " ", "%20", -1)
	// Build the URL query filter string...
	urlQuery = "&f=WHERE%20updated_at%20%3E%27" + ut + "%27"
	fmt.Println("Checking back", *backdays, "days - UTC", t)
	fmt.Println(urlQuery)
}

func syncData(entity, ids string) {

	// the ids to sync...
	xi := []int{}

	// the correct api and sync func
	api := ""
	var syncFunc func(int, *sync.WaitGroup)
	switch entity {
	case "members":
		api = apiMemberIDList
		syncFunc = getMember
	case "modules":
		api = apiModuleIDList
		syncFunc = getModule
	case "resources":
		api = apiResourceIDList
		syncFunc = getResource
	}

	// id list passed in as a comma-separated string, otherwise, fetch the resources to sync based on query...
	if len(ids) > 0 {
		fmt.Printf(" with ids %v\n", ids)
		// Convert strings to int, skip any duds
		xi = stringToSliceInt(ids)
	} else {
		fmt.Printf(" updated in the last %v days\n", *backdays)
		xi = getIDs(api + urlQuery)
	}

	if len(xi) > 0 {
		fmt.Println(len(xi), entity, "to sync...")
		syncGroup(xi, waitGroupSize, &wg, syncFunc)
	} else {
		fmt.Println("No", entity, "to sync")
	}
}

func getIDs(url string) []int {

	fmt.Println("Getting ids at", url)
	var j struct {
		Data []int `json:"data"`
	}

	// Set up the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	// Add auth header
	req.Header.Add("Authorization", "Bearer "+token)

	// Do the request with our global client
	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.NewDecoder(res.Body).Decode(&j)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	return j.Data
}

// syncGroup sets up a loop to sync and entity, with a wait group of a set size.
// That is - sync all of the entities with ids specified in 'ids', with 'set' number of concurrent go routines,
// using the function 'syncFunc'
func syncGroup(ids []int, set int, wg *sync.WaitGroup, syncFunc func(int, *sync.WaitGroup)) {

	i := 0
	for range ids {
		// ... sub loop to wait for processes to finish
		if len(ids)-i == 0 {
			fmt.Println("All done")
			break

		} else if len(ids)-i < set {
			set = len(ids) - i
		}

		fmt.Println("--- next set of", set)
		for l := 0; l < set; l++ {
			wg.Add(1)
			go syncFunc(ids[i], wg)
			i++
		}
		wg.Wait()
	}
}

// getMember fetches the member by id, from the MySQL database, via the API.
// This, in turn, triggers a check of the updatedAt field in MongoDB. If they are out
// of sync the API will update the document record with a fresh one assembled from the MySQL db.
func getMember(id int, wg *sync.WaitGroup) {

	defer wg.Done()

	url := apiMembers + strconv.Itoa(id)
	fmt.Println("GET", url)

	// Set up the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	// Add auth header
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	//m := Member{}
	//json.NewDecoder(res.Body).Decode(&m)
	//fmt.Println(m)
	// Body needs to be FULLY read before closing to
	// ensure that the http connection can be re-used
	// ref: http://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang
	io.Copy(ioutil.Discard, res.Body)
}

// getResource works identically to getMember, except for Resource records.
func getResource(id int, wg *sync.WaitGroup) {

	defer wg.Done()

	url := apiResources + strconv.Itoa(id)
	fmt.Println("GET", url)

	// Set up the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	// Add auth header
	req.Header.Add("Authorization", "Bearer "+token)

	//r := Resource{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	//json.NewDecoder(res.Body).Decode(&r)
	//fmt.Println(r)
	// Body needs to be FULLY read before closing to
	// ensure that the http connection can be re-used
	// ref: http://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang
	io.Copy(ioutil.Discard, res.Body)
}

// getModule works identically to getMember and getResource
func getModule(id int, wg *sync.WaitGroup) {

	defer wg.Done()

	url := apiModules + strconv.Itoa(id)
	fmt.Println("GET", url)

	// Set up the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	// Add auth header
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	io.Copy(ioutil.Discard, res.Body)
}

// stringToSliceInt converts a string like "1234,4567,8987,6543" to a slice of integers. this is used to sort out
// the ids passed as command-line args
func stringToSliceInt(s string) []int {

	xi := []int{}
	xs := strings.Split(s, ",")
	for _, v := range xs {
		fmt.Println(v)
		i, err := strconv.Atoi(v)
		if err == nil {
			xi = append(xi, i)
		} else {
			fmt.Println(err)
		}
	}

	return xi
}
