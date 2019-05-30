package resource

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/utility"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Resource record
type Resource struct {
	OID          bson.ObjectId          `json:"_id,omitempty" bson:"_id,omitempty"`
	ID           int                    `json:"id" bson:"id"`
	Active       bool                   `json:"active" bson:"active"`
	CreatedAt    time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt" bson:"updatedAt"`
	PubDate      PubDate                `json:"pubDate" bson:"pubDate"`
	TypeID       int                    `json:"typeId" bson:"typeId"`
	Type         string                 `json:"type" bson:"type"`
	Primary      bool                   `json:"primary" bson:"primary"`
	Name         string                 `json:"name" bson:"name"`
	Description  string                 `json:"description" bson:"description"`
	Keywords     []string               `json:"keywords" bson:"keywords"`
	ResourceURL  string                 `json:"resourceUrl" bson:"resourceUrl"`
	ShortURL     string                 `json:"shortUrl" bson:"shortUrl"`
	ThumbnailURL string                 `json:"thumbnailUrl" bson:"thumbnailUrl"`
	Attributes   map[string]interface{} `json:"attributes" bson:"attributes"`
}

type PubDate struct {
	Date  time.Time `json:"date" bson:"date"`
	Year  int       `json:"year" bson:"year"`
	Month int       `json:"month" bson:"month"`
	Day   int       `json:"day" bson:"day"`
}

type Resources []Resource

// ByID fetches a resource by id, from the MySQL data
func ByID(ds datastore.Datastore, id int) (*Resource, error) {

	// Set up a new empty Member
	r := Resource{ID: id}

	// Coalesce any NULL-able fields
	query := `
	SELECT
	olr.active,
	olr.created_at,
	olr.updated_at,
	COALESCE(olr.presented_on, ''),
	COALESCE(olr.presented_year, ''),
	COALESCE(olr.presented_month, ''),
	COALESCE(olr.presented_date, ''),
	olr.ol_resource_type_id,
	COALESCE(olrt.name, ''),
	olr.primary,
	COALESCE(olr.name, ''),
	COALESCE(olr.description, ''),
	COALESCE(olr.keywords, ''),
	COALESCE(olr.resource_url, ''),
	COALESCE(olr.short_url, ''),
	COALESCE(olr.thumbnail_url, ''),
	COALESCE(olr.attributes, '')
	FROM ol_resource olr
	LEFT JOIN ol_resource_type olrt ON olr.ol_resource_type_id = olrt.id
        WHERE olr.id = ?`

	// Hold these until we fix them up
	var createdAt string
	var updatedAt string
	var presentedOn string
	var presentedYear string
	var presentedMonth string
	var presentedDay string
	var keywords string
	var attributes string

	err := ds.MySQL.Session.QueryRow(query, id).Scan(
		&r.Active,
		&createdAt,
		&updatedAt,
		&presentedOn,
		&presentedYear,
		&presentedMonth,
		&presentedDay,
		&r.TypeID,
		&r.Type,
		&r.Primary,
		&r.Name,
		&r.Description,
		&keywords,
		&r.ResourceURL,
		&r.ShortURL,
		&r.ThumbnailURL,
		&attributes,
	)
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("ByID() could not find record with id %v", id)
		log.Println(msg, err)
		return &r, errors.Wrap(err, msg)
	case err != nil:
		msg := "ByID() sql error"
		log.Println(msg, err)
		return &r, errors.Wrap(err, msg)
	}

	// Convert MySQL date time strings to time.Time
	r.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		msg := fmt.Sprintf("ByID() record %v - could not Parse created_at", id)
		fmt.Println(msg, err)
		//os.Exit(1)
	}
	r.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
	if err != nil {
		msg := fmt.Sprintf("ByID() record %v - could not Parse updated_at", id)
		fmt.Println(msg, err)
		//os.Exit(1)
	}

	if presentedOn != "" && presentedOn != "0000-00-00" {
		r.PubDate.Date, err = time.Parse("2006-01-02", presentedOn)
		if err != nil {
			msg := fmt.Sprintf("ByID() record %v - could not Parse presented_on", id)
			fmt.Println(msg, err)
			//os.Exit(1)
		}
	}

	// Convert year, month and day strings to int
	r.PubDate.Year, err = strconv.Atoi(presentedYear)
	if err != nil {
		r.PubDate.Year = 0
		//return &r, err
	}
	r.PubDate.Month, err = strconv.Atoi(presentedMonth)
	if err != nil {
		r.PubDate.Month = 0
		//return &r, err
	}
	r.PubDate.Day, err = strconv.Atoi(presentedDay)
	if err != nil {
		r.PubDate.Day = 0
		//return &r, err
	}

	// Convert keyword comma-separated string to a []string
	r.Keywords = strings.Split(keywords, ",")
	for i, v := range r.Keywords {
		r.Keywords[i] = strings.Trim(v, " ")
	}

	// Unmarshal the JSON attributes (string) from MySQL
	xb := []byte(attributes)
	if len(xb) > 0 {
		err = json.Unmarshal(xb, &r.Attributes)
		if err != nil {
			// for now let's not treat this as an error as it is overkill
			// to bomb out if the JSON is not formed correctly
			fmt.Println("ByID() could not unmarshal attributes json string for resource id", r.ID, " - it might be malformed")
			// return &r, err
		}
	}

	return &r, nil
}

// DocResourcesAll searches the Resource collection. Receives query(query) and projection(p)
// It returns []interface{} so that only the projected fields are present. The down side of
// this is that the fields are returned in alphabetical order so it is not as readable
// as the Member struct. Option might be to use the Member struct when no projection
// is specified. TODO - see if we can use a the proper struct when there is no projection
func DocResourcesAll(ds datastore.Datastore, q map[string]interface{}, p map[string]interface{}) ([]interface{}, error) {

	resources, err := ds.MongoDB.ResourcesCollection()
	if err != nil {
		return nil, err
	}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt"})

	// Run query and return results
	var r []interface{}
	err = resources.Find(q).Select(p).All(&r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// DocResourcesLimit returns n resources
func DocResourcesLimit(ds datastore.Datastore, q map[string]interface{}, p map[string]interface{}, l int) ([]interface{}, error) {

	r := []interface{}{}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt"})

	resources, err := ds.MongoDB.ResourcesCollection()
	if err != nil {
		return r, err
	}
	err = resources.Find(q).Select(p).Limit(l).All(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// DocResourcesOne returns one resource, unmarshaled into the proper struct so no projection allowed here
func DocResourcesOne(ds datastore.Datastore, q map[string]interface{}) (Resource, error) {

	r := Resource{}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt"})

	resources, err := ds.MongoDB.ResourcesCollection()
	if err != nil {
		return r, err
	}
	err = resources.Find(q).One(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// QueryResourcesCollection ... queries the resources collection :)
func QueryResourcesCollection(ds datastore.Datastore, mq datastore.MongoQuery) ([]interface{}, error) {

	// results
	r := []interface{}{}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(mq.Find, []string{"updatedAt", "createdAt"})

	// get a pointer to the resources collection
	c, err := ds.MongoDB.ResourcesCollection()
	if err != nil {
		return r, err
	}

	// execute query, scan results into r
	err = mq.Do(c, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// FetchResources returns values of type Resource from the Resources collection in MongoDB, based on the query and
// limited by the value of limit. If limit is 0 all results are returned.
func FetchResources(ds datastore.Datastore, query map[string]interface{}, limit int) ([]Resource, error) {

	var data []Resource

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(query, []string{"updatedAt", "createdAt"})

	c, err := ds.MongoDB.ResourcesCollection()
	if err != nil {
		return nil, err
	}
	err = c.Find(query).Limit(limit).All(&data)

	return data, err
}

// SyncResource synchronises the Resource record from MySQL -> MongoDB
func SyncResource(ds datastore.Datastore, r *Resource) {

	// Fetch the current Doc (if there) and compare updatedAt
	r2, err := DocResourcesOne(ds, bson.M{"id": r.ID})
	if err != nil {
		log.Println("Target document error: ", err, "- so do an upsert")
	}

	msg := fmt.Sprintf("Resource id %v - MySQL updated %v, MongoDB updated %v", r.ID, r.UpdatedAt, r2.UpdatedAt)
	if r.UpdatedAt.Equal(r2.UpdatedAt) {
		msg += " - NO need to sync"
		log.Println(msg)
		return
	}
	msg += " - syncing..."
	log.Println(msg)

	// Update the document in the Members collection
	var w sync.WaitGroup
	w.Add(1)
	go updateResourceDoc(ds, r, &w)
	w.Wait()
}

// UpdateMemberDoc updates the JSON-formatted member record in MongoDB
func updateResourceDoc(ds datastore.Datastore, r *Resource, w *sync.WaitGroup) {

	// Make the selector for Upsert
	id := map[string]int{"id": r.ID}

	// Get pointer to the collection
	mc, err := ds.MongoDB.ResourcesCollection()
	if err != nil {
		log.Printf("Error getting pointer to Resources collection: %s\n", err.Error())
		return
	}

	// Upsert
	_, err = mc.Upsert(id, &r)
	if err != nil {
		log.Printf("Error updating document in Resources collection: %s\n", err.Error())
	}

	// Tell wait group we're done, if it was passed in
	if w != nil {
		w.Done()
	}

	log.Println("Updated Resources document")
}

// Save a Resource to MySQL. Returns the id of the new record, and an error. If the record appears to be a duplicate
// then this will hand off to .Update() to see if the record should be updated instead.
func (r *Resource) Save(ds datastore.Datastore) (int, error) {

	// Don't save a resource without a target url
	if r.ResourceURL == "" {
		return 0, errors.New("Cannot save a resource without a url")
	}

	// set r.ID if there is a matching resource url in the database and update only if it is not an exact match
	r.ID = ResourceByURL(ds, r.ResourceURL)
	if r.ID > 0 {

		nothingToUpdate, err := r.ExactDatabaseMatch(ds)
		if err != nil {
			msg := fmt.Sprintf("Error checking equality of Resource id %v value with its counterpart in the database: %s - cannot update", r.ID, err)
			return int(r.ID), errors.New(msg)
		}
		if nothingToUpdate {
			msg := fmt.Sprintf("Resource id %v update appears to be identical with its counterpart in the database - nothing to update", r.ID)
			return int(r.ID), errors.New(msg)
		}

		err = r.Update(ds, r.ID)
		if err != nil {
			msg := fmt.Sprintf("Update failed for Resource %v with err: %v", r.ID, err)
			return int(r.ID), errors.New(msg)
		}

		return int(r.ID), nil
	}

	return r.Add(ds)
}

// Add inserts a Resource into the mysql database and returns the new record id
func (r *Resource) Add(ds datastore.Datastore) (int, error) {

	query := "INSERT INTO ol_resource (ol_resource_type_id, active, `primary`," +
		`created_at, updated_at, presented_on, presented_year, presented_month, presented_date,
		name, description, keywords,
		resource_url, short_url, thumbnail_url, attributes)
		VALUES (%v, %v, %v,
		"%v", "%v", "%v", "%v", "%v", "%v",
		"%v", "%v", "%v",
		"%v", "%v", "%v", "%v")`

	// Create comma separated keyword list for MySQL field from r.Keywords []string
	keywords := strings.Join(r.Keywords, ",")

	// Marshall the Attributes field to a string for MySQL
	// The ugly strings.Replace() below is to be able to insert string literal JSON
	// into ol_resource.attributes in MySQL.
	// To store the JSON as a string, eg {"free": true, "public": true, "source": "Pubmed"}
	// it needs to be escaped for the INSERT, thus:
	// {\"free\": true, \"public\": true, \"source\": \"Pubmed\"}
	attributes := ""
	xb, err := json.Marshal(r.Attributes)
	if err != nil {
		log.Println("Could not marshal attributes prior to insert - ignoring attributes completely")
	} else {
		attributes = strings.Replace(string(xb), "\"", "\\\"", -1)
	}

	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()

	query = fmt.Sprintf(query, r.TypeID, 1, r.Primary,
		r.CreatedAt.Format("2006-01-02 15:04:05"),
		r.UpdatedAt.Format("2006-01-02 15:04:05"),
		r.PubDate.Date.Format("2006-01-02"),
		r.PubDate.Year,
		r.PubDate.Month,
		r.PubDate.Day,
		r.Name, r.Description, keywords,
		r.ResourceURL, r.ShortURL, r.ThumbnailURL, attributes)

	res, err := ds.MySQL.Session.Exec(query)
	if err != nil {
		msg := fmt.Sprintf("Error with query: %s\nError: %s", query, err)
		return 0, errors.New(msg)
	}

	id, err := res.LastInsertId()
	if err != nil {
		msg := fmt.Sprintf("Error fetching last insert id: %s", err)
		return 0, errors.New(msg)
	}
	r.ID = int(id)
	log.Println("Added a new resource with ID", r.ID)

	err = r.SetShortURL(ds)
	if err != nil {
		log.Println("Error setting short url for new resource:", err)
	}

	return r.ID, nil
}

// Update performs an update operation on a Resource record in MySQL
func (r *Resource) Update(ds datastore.Datastore, id int) error {

	// Only difference with Save() query is that created_at is not included
	query := "UPDATE ol_resource SET ol_resource_type_id = %v, active = %v, `primary` = %v," +
		`updated_at = "%v", presented_on = "%v", presented_year = "%v", presented_month = "%v", presented_date = "%v",
		name = "%v", description= "%v", keywords= "%v",
		resource_url = "%v", short_url = "%v", thumbnail_url = "%v"
		WHERE id = %v`

	// Create comma separated keyword list for MySQL field from r.Keywords []string
	keywords := strings.Join(r.Keywords, ",")

	r.UpdatedAt = time.Now()

	query = fmt.Sprintf(query, r.TypeID, 1, r.Primary,
		r.UpdatedAt.Format("2006-01-02 15:04:05"),
		r.PubDate.Date.Format("2006-01-02"),
		r.PubDate.Year,
		r.PubDate.Month,
		r.PubDate.Day,
		r.Name, r.Description, keywords,
		r.ResourceURL, r.ShortURL, r.ThumbnailURL,
		id)

	_, err := ds.MySQL.Session.Exec(query)
	if err != nil {
		fmt.Println("Query error:")
		return err
	}

	err = r.SetShortURL(ds)
	if err != nil {
		msg := fmt.Sprintf("setShortURL error: %v", err)
		return errors.New(msg)
	}

	return nil
}

// ExactDatabaseMatch returns true if a exactly matching record exists in the database. The bool response
// is only useful if there is no error, otherwise a false negative might lead to duplicate records being created
func (r *Resource) ExactDatabaseMatch(ds datastore.Datastore) (bool, error) {

	dbResource, err := ByID(ds, r.ID)
	if err != nil {
		msg := fmt.Sprintf("Could not fetch existing resource id %v, err: %v - skipping this record", r.ID, err)
		return false, errors.New(msg)
	}

	if ResourceDeepEqual(r, dbResource) {
		return true, nil
	}

	return false, nil
}

// SetShortURL sets the short_url fields in MySQL ol_resource.short_url. This is here for convenience so that when
// a Resource is added via the API, and we don't yet have an ID, we can set the short url without making another API call.
// This is a bit "hackish", however the current MappCPD application will send a user to the short url if it exists in the data
// -e, if there is a value in ol_resource.short_url. So BEFORE we set this value we need to make sure a record exists in the Links
// collection in MongoDB as this is what the short url service (linkr) refers to when doing short link redirecting
func (r *Resource) SetShortURL(ds datastore.Datastore) error {

	// Don't set a short url if there is no long url as a target!
	if r.ResourceURL == "" {
		return errors.New("Cannot set a short url when there is no long url!")
	}

	// Check the Links collection for this Resource (URL)
	q := bson.M{"longUrl": r.ResourceURL}
	l, err := DocLinksOne(ds, q)
	if err != nil {
		switch err {
		case mgo.ErrNotFound:
			fmt.Println("SetShortURL no Links doc so will create one...")
			l.CreatedAt = time.Now()
			l.UpdatedAt = time.Now()
			l.ShortUrl = "r" + strconv.Itoa(r.ID)
			l.LongUrl = r.ResourceURL
			l.Title = r.Name
			err := l.DocSave(ds)
			if err != nil {
				fmt.Println(err)
				return err
			}
		default:
			fmt.Println("Error:", err, "- cannot continue setting short link")
			return err
		}
	}

	// Finally, update the ol_resource.short_url value
	shortUrl := os.Getenv("MAPPCPD_SHORT_LINK_URL") + "/" + os.Getenv("MAPPCPD_SHORT_LINK_PREFIX") + strconv.Itoa(r.ID)
	query := fmt.Sprintf("UPDATE ol_resource SET short_url = \"%v\" WHERE id = %v", shortUrl, r.ID)
	fmt.Println("SetShortLinkURL():", query)
	_, err = ds.MySQL.Session.Exec(query)
	if err != nil {
		fmt.Println("SQL error with query: ", query, " -", err)
		return err
	}

	return nil
}

// ResourceByURL checks the MySQL data for the existence of a resource_url, and returns the FIRST id or 0
// This is here to help prevent the addition of duplicate resources. In the .Save() func if a duplicate is found,
// that is, a duplicate resource_url, .Update() will be run instead. This means the id is required so we know which record
// to update. It is possible there could be more than one duplicate so, for now, return the FIRST duplicate with QueryRow()
func ResourceByURL(ds datastore.Datastore, url string) int {

	var c int

	// Don't allow this to match an empty url string
	if url == "" {
		return c
	}

	query := "SELECT id FROM ol_resource WHERE active = 1 AND resource_url = ?"
	ds.MySQL.Session.QueryRow(query, url).Scan(&c)

	return c
}

// ResourceDeepEqual compares key fields in a Resource and returns true if all are equal
func ResourceDeepEqual(r, r2 *Resource) bool {

	if r.Name != r2.Name {
		return false
	}
	if r.Description != r2.Description {
		return false
	}
	if !reflect.DeepEqual(r.PubDate, r2.PubDate) {
		return false
	}
	if !reflect.DeepEqual(r.Keywords, r2.Keywords) {
		return false
	}

	return true
}

// Sync saves the Resource to the document database.
func (r *Resource) Sync(ds datastore.Datastore) error {
	return r.SaveDoc(ds)
}

// SaveDoc upserts Resource doc to MongoDB
func (r *Resource) SaveDoc(ds datastore.Datastore) error {

	rc, err := ds.MongoDB.ResourcesCollection()
	if err != nil {
		return fmt.Errorf("resource.SaveDoc() err = %s", err)
	}

	selector := map[string]int{"id": r.ID}
	_, err = rc.Upsert(selector, &r)
	if err != nil {
		return fmt.Errorf("resource.SaveDoc() err = %s", err)
	}

	return nil
}
