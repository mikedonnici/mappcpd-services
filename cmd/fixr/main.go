package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/34South/envr"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/resource"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type link struct {
	id        int
	title     string
	shortPath string
	shortURL  string
	longURL   string
}

type resourceData struct {
	ID         int
	Name       string
	Keywords   string
	PubDate    string
	PubYear    string
	PubMonth   string
	PubDay     string
	Attributes attributes
}

// ol_resource.attributes holds a JSON string which should map to this:
type attributes struct {
	Category         string `json:"category"`
	Free             bool   `json:"free"`
	Public           bool   `json:"public"`
	Source           string `json:"source"`
	SourceID         string `json:"sourceId"`
	SourceName       string `json:"sourceName"`
	SourceNameAbbrev string `json:"sourceNameAbbrev"`
	SourcePubDate    string `json:"sourcePubDate"`
	SourceVolume     string `json:"sourceVolume"`
	SourceIssue      string `json:"sourceIssue"`
	SourcePages      string `json:"sourcePages"`
}

// backdays specifies how far back to include records in whatever task is being performed
var backdays int

// tasksFlag flag is used to specify specific functions to run, comma-separated
var tasksFlag string

var validTasks = []string{"fixResources", "pubmedData"}

var DS datastore.Datastore

func init() {
	envr.New("fixrEnv", []string{
		"MAPPCPD_SHORT_LINK_URL",
		"MAPPCPD_SHORT_LINK_PREFIX",
		"MAPPCPD_MONGO_DBNAME",
		"MAPPCPD_MONGO_DESC",
		"MAPPCPD_MONGO_URL",
		"MAPPCPD_MYSQL_DESC",
		"MAPPCPD_MYSQL_URL",
	}).Auto()

	// flags
	flag.IntVar(&backdays, "b", 1, "Specify backdays as an integer > 0")
	flag.StringVar(&tasksFlag, "t", "", "Specify comma-separated list of tasks to run, 'task1, task2, task3'")

	var err error
	DS, err = datastore.FromEnv()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {

	// Flag check
	flag.Parse()
	if backdays == 1 {
		fmt.Println("Backdays not specified with -b flag, defaulting to 1")
	} else {
		fmt.Println("Checking records updated within the last", backdays, "days")
	}

	if tasksFlag == "" {
		fmt.Println("No tasks specified")
		os.Exit(0)
	}
	t := strings.Replace(tasksFlag, " ", "", -1)
	tasks := strings.Split(t, ",")
	if err := verifyTasks(tasks); err != nil {
		fmt.Println(errors.Cause(err))
		os.Exit(1)
	}

	for _, v := range tasks {

		if v == "fixResources" {
			fmt.Println("Running task:", v)
			if err := checkShortLinks(); err != nil {
				fmt.Println(errors.Cause(err))
				os.Exit(1)
			}
			if err := syncActiveFlag(); err != nil {
				fmt.Println(errors.Cause(err))
				os.Exit(1)
			}
			fmt.Println("--- done")
		}

		if v == "pubmedData" {
			fmt.Println("Running task:", v)
			updatePubmedData()
		}
	}
}

// verifyTasks checks that the task list contains tasks that are valid
func verifyTasks(tasks []string) error {

	for _, t := range tasks {

		f := false
		for _, vt := range validTasks {
			if t == vt {
				f = true
			}
		}

		if f == false {
			return fmt.Errorf("Invalid task: '%s'", t)
		}
	}

	return nil
}

// checkShortLinks checks all the Resource records for a short link, and if incorrect or not found, fixes them.
func checkShortLinks() error {

	// Select resources that start with 'http%' so don't break relative URLs
	// Note the short_url value from the primary record can be NULL.
	// When this is the case the .Scan method below bombs out - hence use of COALESCE.
	query := "SELECT id, name, COALESCE(short_url, ''), resource_url FROM ol_resource " +
		"WHERE `active` = 1 AND `primary` = 1 AND resource_url LIKE 'http%' " +
		"AND updated_at >= NOW() - INTERVAL " + strconv.Itoa(backdays) + " DAY"

	rows, err := DS.MySQL.Session.Query(query)
	if err != nil {
		return err
	}

	for rows.Next() {
		l := link{}
		err := rows.Scan(&l.id, &l.title, &l.shortURL, &l.longURL)
		if err != nil {
			return err
		}

		// Work out what we expect, or need to set up for a short link
		// Custom short URL based on id of resourceData with no padding, eg r12, r3435
		l.shortPath = fmt.Sprintf("%v%v", os.Getenv("MAPPCPD_SHORT_LINK_PREFIX"), l.id)
		// The short_url should be
		expectedShortURL := fmt.Sprintf("%v/%v", os.Getenv("MAPPCPD_SHORT_LINK_URL"), l.shortPath)
		// fmt.Printf("/%s -> %s", l.shortPath, expectedShortURL)

		// There are two scenarios:
		// 1. short_url already has the expected value in the primary store
		// Look for any differences between the primary record and the Links doc, if found then sync changes.
		//
		// 2. Missing, invalid or does not match expected based on config
		// Need to set the full short URL in the primary record, and then sync as above

		// Set in primary record, if not the expected value...
		if l.shortURL != expectedShortURL {
			fmt.Printf("/%s -> %s", l.shortPath, expectedShortURL)
			fmt.Println("...no short url - will create one and then sync")
			l.shortURL = expectedShortURL
			err := setShortURL(l.id, l.shortURL)
			if err != nil {
				return err
			}
		}

		// Check for changes and sync if required...
		err = checkSync(l)
		if err != nil {
			return err
		}
	}

	return nil
}

// setShortURL sets the value of the short_url field in the ol_resource (primary) record. It also updates the
// updated_at value to ensure this record will be picked up for sync later on (by mongr)
func setShortURL(id int, shortURL string) error {

	query := `UPDATE ol_resource SET short_url = "%v", updated_at = NOW() WHERE id = %v LIMIT 1`
	query = fmt.Sprintf(query, shortURL, id)
	_, err := DS.MySQL.Session.Exec(query)

	return errors.Wrap(err, "sql updated failed")
}

// checkSync looks for differences between the primary record and the Links doc, and
// does an upsert of needed. Note that the full short url is only stored in the ol_resource
// record. In the Links doc only the shortPath is stored as all links are relative to whatever
// base URL the short link redirector is using.
func checkSync(l link) error {

	// get link doc from shortPath
	ld, err := getLinkDoc(l.shortPath)
	if err == mgo.ErrNotFound {
		fmt.Println("No Link doc found so need to create one")
	}
	if err != nil && err != mgo.ErrNotFound {
		return errors.Wrap(err, "mongo query failed")
	}

	// No we can check if anything has changed, and decide if we want to
	doSync := false

	// Make sure the link doc has dates set
	nilTime := time.Time{}
	if ld.CreatedAt == nilTime {
		fmt.Println("CreatedAt is nil - will set to now")
		ld.CreatedAt = time.Now()
		doSync = true
	}
	if ld.UpdatedAt == nilTime {
		fmt.Println("UpdatedAt is nil - will set to now")
		ld.UpdatedAt = time.Now()
		doSync = true
	}
	if ld.Title != l.title {
		ld.Title = l.title
		doSync = true
	}
	if ld.LongUrl != l.longURL {
		ld.LongUrl = l.longURL
		doSync = true
	}

	// sync if needed
	if doSync == true {
		fmt.Println("Do the sync!")
		err := sync(ld, l.shortPath)
		return errors.Wrap(err, "sync failed")
	}

	return nil
}

// getLinkDoc fetches a doc from the Links collection.If the doc is not found it
// returns a valid, but empty, value of type resources.Link.
func getLinkDoc(shortPath string) (resource.Link, error) {

	var l resource.Link

	// Always set this as it is used as the KEY for Links docs
	l.ShortUrl = shortPath

	// connect to collection
	c, err := DS.MongoDB.LinksCol()
	if err != nil {
		return l, errors.Wrap(err, "error connecting to collection")
	}

	// query
	s := bson.M{"shortUrl": shortPath}
	err = c.Find(s).One(&l)
	// Not found is ok, return the empty value
	if err == mgo.ErrNotFound {
		return l, nil /// no error, just not found
	}
	// Some other error
	if err != nil {
		return l, errors.Wrap(err, "mongo query failed")
	}

	// An existing record found
	return l, nil
}

// sync upserts a doc to the Links collection. The selector is the shortPath value, eg /r123.
// This value will always be the same and we don't know if we are creating a new record
// or modifying an existing one. So shortPath is the best identifier.
func sync(ld resource.Link, shortPath string) error {

	c, err := DS.MongoDB.LinksCol()
	if err != nil {
		return errors.Wrap(err, "error connecting to collection")
	}

	s := bson.M{"shortUrl": shortPath}
	u := bson.M{"$set": ld}
	_, err = c.Upsert(s, u)
	if err != nil {
		return errors.Wrap(err, "upsert failed")
	}

	return nil
}

// getResourceIDs fetches all of the resourceData ids with the specified active (soft-delete) status, from the primary db.
func getResourceIDs(active bool) ([]int, error) {

	var ids []int

	// active is stored as 0/1 in MySQL, true/false in MongoDB
	var a int
	if active == true {
		a = 1
	} else {
		a = 0
	}

	// Only resources with proper urls, even though there are very few with relative urls
	// Don't need back days - need to check everything.
	query := "SELECT id FROM ol_resource WHERE resource_url LIKE 'http%' AND active = ?"
	rows, err := DS.MySQL.Session.Query(query, a)
	if err != nil {
		return nil, errors.New("Error executing query - " + err.Error())
	}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, errors.New("Error scanning row - " + err.Error())
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// syncActiveFlag ensures that the active status is the same for resources stored in MySQL and in MongoDB.
// The primary DB is the authoritative source so always use the ol_resource.active field value.
func syncActiveFlag() error {

	// Fetch all inactive resources ids from primary db
	inactiveIDs, err := getResourceIDs(false)
	if err != nil {
		return err
	}
	fmt.Println("Found", len(inactiveIDs), "inactive resources")

	// Set the corresponding Resources docs active field to 'false'
	if err := setResourcesActiveField(inactiveIDs, false); err != nil {
		return err
	}
	// Set the corresponding Links docs active field to 'false'
	if err := setLinksActiveField(inactiveIDs, false); err != nil {
		return err
	}

	// Fetch all active resources ids from primary db
	activeIDs, err := getResourceIDs(true)
	if err != nil {
		return err
	}
	fmt.Println("Found", len(activeIDs), "active resources")

	// Set the corresponding Resources docs active field to 'true'
	if err := setResourcesActiveField(activeIDs, true); err != nil {
		return err
	}

	// Set the corresponding Links docs active field to 'false'
	return setLinksActiveField(activeIDs, true)
}

// setResourcesActiveField sets the active field for Resources docs identified by the ids list
func setResourcesActiveField(ids []int, active bool) error {

	if len(ids) == 0 {
		fmt.Println("No ids passed in to setResourcesActiveField()... nothing to do.")
		return nil
	}

	fmt.Printf("Setting %v Resources to active: %v", len(ids), active)

	c, err := DS.MongoDB.ResourcesCollection()
	if err != nil {
		return errors.Wrap(err, "error connecting to collection")
	}

	s := bson.M{"id": bson.M{"$in": ids}}
	u := bson.M{"$set": bson.M{"active": active}}
	ci, err := c.UpdateAll(s, u)
	if err != nil {
		return errors.Wrap(err, "Mongo query error")
	}
	fmt.Println("...", ci.Updated, "records were updated")

	return nil
}

// setLinksActiveField sets the active field for Links docs identified by the ids list. There is no id field
// in Links docs so use the shortUrl by pre-pending 'r', so id 1789 becomes 'r1789'
func setLinksActiveField(ids []int, active bool) error {

	if len(ids) == 0 {
		fmt.Println("No ids passed in to setLinksActiveField()... nothing to do.")
		return nil
	}

	// prepend ids with 'r' to create selector
	var rids []string
	for _, v := range ids {
		rid := os.Getenv("MAPPCPD_SHORT_LINK_PREFIX") + strconv.Itoa(v)
		rids = append(rids, rid)
	}

	fmt.Printf("Setting %v Links to active: %v", len(ids), active)

	c, err := DS.MongoDB.LinksCol()
	if err != nil {
		return errors.Wrap(err, "error connecting to collection")
	}

	s := bson.M{"shortUrl": bson.M{"$in": rids}}
	u := bson.M{"$set": bson.M{"active": active}}
	ci, err := c.UpdateAll(s, u)
	if err != nil {
		return errors.Wrap(err, "Mongo query error")
	}

	fmt.Println("...", ci.Updated, "records were updated")

	return nil
}

// updatePubmedData checks and updates the ol_resources record for resources that were sourced from Pubmed
func updatePubmedData() {

	// fetch pubmed resources
	xr, err := resourcesByAttribute("pubmed")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// set new ol_resource.Attributes for each resourceData
	for _, r := range xr {

		// last keyword in ol_resource.keywords is the pubmed article id
		pubmedId, err := strconv.Atoi(lastKeyword(r.Keywords))
		if err != nil {
			fmt.Println("Last keyword does not appear to be an id")
			os.Exit(1)
		}

		// set Attributes related to pubmed data
		r.pubmedData(strconv.Itoa(pubmedId))
	}
}

// resourceByAttribute fetches resourceData records that have a string LIKE 'pattern' in the ol_resource.Attributes field
func resourcesByAttribute(pattern string) ([]resourceData, error) {

	var xr []resourceData

	query := "SELECT id, name, keywords, Attributes FROM ol_resource " +
		"WHERE Attributes LIKE '%" + pattern + "%' " +
		"AND updated_at >= NOW() - INTERVAL " + strconv.Itoa(backdays) + " DAY "
	rows, err := DS.MySQL.Session.Query(query)
	if err != nil {
		return xr, err
	}

	for rows.Next() {
		var r resourceData
		var attributes []byte
		err := rows.Scan(&r.ID, &r.Name, &r.Keywords, &attributes)
		if err != nil {
			return xr, err
		}

		// Unmarshal attributes JSOn string, into Attributes value
		json.Unmarshal(attributes, &r.Attributes)

		xr = append(xr, r)
	}

	return xr, nil
}

// lastKeyword returns the last keyword in a comma separated list of keywords.
func lastKeyword(list string) string {
	xs := strings.Split(list, ",")
	return xs[len(xs)-1]
}

// updatePubmedData fetches data for an article by Pubmed id
//
// The JSON response from Pubmed is shaped as shown below, ie the field of interest is named after the id of the article.
// Hence, it is easiest to unmarshal the response into a map[string]interface{}
//
//  {
//	  "meta": [],
//	  "result": {
//	  	  "1234": {
//			  "field1": "value1",
//			  "field2": "value2",
//			  "field3": "value3",
//			  "fulljournalname": "Medicine"
//		  }
//	  }
//  }
//
// Example: https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esummary.fcgi?db=pubmed&id=25963440&retmode=json
func (r *resourceData) pubmedData(articleID string) {

	url := fmt.Sprintf("https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esummary.fcgi?db=pubmed&id=%s&retmode=json", articleID)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Could not GET", url)
	}
	defer res.Body.Close()

	xb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Could not read response body")
	}

	// unmarshall the entire response body
	rb := make(map[string]interface{})
	json.Unmarshal(xb, &rb)

	// assert "result" field
	rsf := rb["result"].(map[string]interface{})

	// assert the field with name equal to the id of the article
	idf := rsf[articleID].(map[string]interface{})

	// set required fields
	r.Attributes.SourceID = articleID

	v, ok := idf["fulljournalname"].(string)
	if ok {
		r.Attributes.SourceName = v
	}

	v, ok = idf["source"].(string)
	if ok {
		r.Attributes.SourceNameAbbrev = v
	}

	v, ok = idf["volume"].(string)
	if ok {
		r.Attributes.SourceVolume = v
	}

	v, ok = idf["issue"].(string)
	if ok {
		r.Attributes.SourceIssue = v
	}

	v, ok = idf["pages"].(string)
	if ok {
		r.Attributes.SourcePages = v
	}

	v, ok = idf["pubdate"].(string)
	if ok {
		r.Attributes.SourcePubDate = v
	}

	// Need to sort out the date... the ACTUAL date
	r.bestDate(idf)

	// update the main record
	err = updateResource(*r)
	if err != nil {
		fmt.Println("Error updating the record:", err)
	}
}

// bestDate attempts to find the best publish date fromt he available data
func (r *resourceData) bestDate(data map[string]interface{}) {

	fmt.Println("Looking for best date...")

	// Month in Pubmed data is *usually* a 3 character string, eg 'May', but sometimes it is a two-character
	// number, eg '05'. So this hack is to determine which prior to creating time value.
	months := map[string]string{
		"Jan": "1", "Feb": "2", "Mar": "3", "Apr": "4", "May": "5", "Jun": "6",
		"Jul": "7", "Aug": "8", "Sep": "9", "Oct": "10", "Nov": "11", "Dec": "12",
	}

	// first 'best' options are "pubdate" and "epubdate" - "epubdate" is usually a bit earlier than "pubdate"
	fmt.Print("Try pubdate: ")

	pubdate, ok := data["pubdate"].(string)
	if ok {
		fmt.Println(pubdate)
		xs := strings.Split(pubdate, " ")

		if len(xs) == 3 {
			r.PubYear, r.PubMonth, r.PubDay = xs[0], xs[1], xs[2]
		}

		if len(xs) == 2 {
			r.PubYear, r.PubMonth = xs[0], xs[1]
		}
	}

	// If the month value is a key in the months array, set it to a 'numerical' value ie '5' instead of 'May'
	m, ok := months[r.PubMonth]
	if ok {
		r.PubMonth = m
	}

	// No day, so set to 1
	if r.PubDay == "" {
		r.PubDay = "1"
	}

	// if we can parse this date then we are sweet
	ts := r.PubYear + "-" + r.PubMonth + "-" + r.PubDay
	_, err := time.Parse("2006-1-2", ts)
	if err != nil {
		fmt.Println("Error parsing date -", err)
		return
	}

	r.PubDate = ts
	fmt.Println("Best publish date:", r.PubDate)
}

// updateResource updates the ol_resource record
func updateResource(r resourceData) error {

	// attributes stores as a JSON string
	attributes, err := json.Marshal(r.Attributes)
	if err != nil {
		return err
	}

	query := `UPDATE ol_resource SET
              updated_at = NOW(),
              presented_on = '%s',
              presented_year = '%s',
              presented_month = '%s',
              presented_date = '%s',
			  attributes = '%s'
			  WHERE id = %v LIMIT 1`
	query = fmt.Sprintf(query, r.PubDate, r.PubYear, r.PubMonth, r.PubDay, attributes, r.ID)

	_, err = DS.MySQL.Session.Exec(query)

	return err
}
