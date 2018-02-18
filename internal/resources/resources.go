package resources

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"database/sql"
	"encoding/json"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/mappcpd/web-services/internal/constants"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/utility"
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

// ResourceByID fetches a resource by id, from the MySQL db
func ResourceByID(id int) (*Resource, error) {

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

	err := datastore.MySQL.Session.QueryRow(query, id).Scan(
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
		msg := fmt.Sprintf("ResourceByID() could not find record with id %v", id)
		log.Println(msg, err)
		return &r, errors.Wrap(err, msg)
	case err != nil:
		msg := "ResourceByID() sql error"
		log.Println(msg, err)
		return &r, errors.Wrap(err, msg)
	}

	// Convert MySQL date time strings to time.Time
	r.CreatedAt, err = time.Parse(constants.MySQLTimestampFormat, createdAt)
	if err != nil {
		msg := fmt.Sprintf("ResourceByID() record %v - could not Parse created_at", id)
		fmt.Println(msg, err)
		//os.Exit(1)
	}
	r.UpdatedAt, _ = time.Parse(constants.MySQLTimestampFormat, updatedAt)
	if err != nil {
		msg := fmt.Sprintf("ResourceByID() record %v - could not Parse updated_at", id)
		fmt.Println(msg, err)
		//os.Exit(1)
	}
	r.PubDate.Date, err = time.Parse(constants.MySQLDateFormat, presentedOn)
	if err != nil {
		msg := fmt.Sprintf("ResourceByID() record %v - could not Parse presented_on", id)
		fmt.Println(msg, err)
		//os.Exit(1)
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
			fmt.Println("ResourceByID() could not unmarshal attributes json string for resource id", r.ID, " - it might be malformed")
			// return &r, err
		}
	}

	return &r, nil
}

// DocResourcesAll searches the Resource collection. Receives query(q) and projection(p)
// It returns []interface{} so that only the projected fields are present. The down side of
// this is that the fields are returned in alphabetical order so it is not as readable
// as the Member struct. Option might be to use the Member struct when no projection
// is specified. TODO - see if we can use a the proper struct when there is no projection
func DocResourcesAll(q map[string]interface{}, p map[string]interface{}) ([]interface{}, error) {

	resources, err := datastore.MongoDB.ResourcesCol()
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
func DocResourcesLimit(q map[string]interface{}, p map[string]interface{}, l int) ([]interface{}, error) {

	r := []interface{}{}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt"})

	resources, err := datastore.MongoDB.ResourcesCol()
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
func DocResourcesOne(q map[string]interface{}) (Resource, error) {

	r := Resource{}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt"})

	resources, err := datastore.MongoDB.ResourcesCol()
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
func QueryResourcesCollection(mq datastore.MongoQuery) ([]interface{}, error) {

	// results
	r := []interface{}{}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(mq.Find, []string{"updatedAt", "createdAt"})

	// get a pointer to the resources collection
	c, err := datastore.MongoDB.ResourcesCol()
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

// SyncResource synchronises the Resource record from MySQL -> MongoDB
func SyncResource(r *Resource) {

	// Fetch the current Doc (if there) and compare updatedAt
	r2, err := DocResourcesOne(bson.M{"id": r.ID})
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
	go UpdateResourceDoc(r, &w)
	w.Wait()
}

// UpdateMemberDoc updates the JSON-formatted member record in MongoDB
func UpdateResourceDoc(r *Resource, w *sync.WaitGroup) {

	// Make the selector for Upsert
	id := map[string]int{"id": r.ID}

	// Get pointer to the collection
	mc, err := datastore.MongoDB.ResourcesCol()
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
func (r *Resource) Save() (int, error) {

	// At this point we will not save a resource without a target url... defeats the purpose
	if r.ResourceURL == "" {
		m := "Cannot r.Save() when there is no resource url"
		fmt.Println(m)
		return 0, errors.New(m)

	}

	// To prevent duplicate resources we need to check for a match based on suitable field(s).
	// For now we can use the resource_url field as there is no point(?) having multiple Resources that point
	// to the same location. However, it might be better to match on multiple fields?
	// If we do find a duplicate we then need to check if we should update the resource as something may have changed.
	// This is an issue when doing batch resource inserts because we may load the same batch, find many duplicates,
	// update them all and hence change the updated_at stamp. This then causes other services to update (mongr->algr)
	// when they don't need to.
	// Note that we don't yet have the field value for r.ID, so we can set it here if we do find a duplicate resource_url,
	// otherwise set it further along, when a new record is created in MySQL.
	// todo flaky check for duplicate, if func returns zero we create a duplicate! only thing saving it is return for blank resource url above
	r.ID = DuplicateResourceURL(r.ResourceURL)
	if r.ID > 0 {
		//fmt.Println("")
		//fmt.Println("#################################################################################")
		//fmt.Printf("Trying to Save() a resource that has the same resource_url (%v) as existing resource id %v\n", r.ResourceURL, r.ID)
		//fmt.Println("Checking if this resource should be updated instead...")

		// See if there is any real need to update, rather than just updating willy nilly :)
		// First, fetch the current Resource with this ID...
		r2, err := ResourceByID(r.ID)
		if err != nil {
			fmt.Printf("Could not fetch the (possible) duplicate resource id %v, err: %v - skipping this record", r.ID, err)
			return int(r.ID), err
		}

		//fmt.Println("")
		//fmt.Println("New Resource:")
		//fmt.Println(r)
		//fmt.Println("")
		//fmt.Println("Existing Resource:")
		//fmt.Println(r2)

		// Resource (r) is our update candidate, r2 is the possible duplicate from the database. Check
		// key fields to see if there is anything that actually needs to be updated...
		if ResourceDeepEqual(r, r2) {
			fmt.Println("Resource", r.ID, "has no significant changes, skipping update")
			return int(r.ID), err
		}

		// Not deeply equal, so update
		fmt.Print("Resource ", r.ID, " has changed and needs to be updated... ")
		err = r.Update(r.ID)
		if err != nil {
			fmt.Printf("Update failed for Resource %v with err: %v - skipping this record", r.ID, err)
			return int(r.ID), err
		}

		// Set ShortUrl
		err = r.SetShortURL()
		if err != nil {
			fmt.Println("setShortURL error:", err)
			return int(r.ID), err
		}

		fmt.Println("done!")
		return int(r.ID), nil
	}

	// No duplicate so add a new resource...
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
		fmt.Println("Could not marshal attributes prior to insert - ignoring attributes completely")
	} else {
		attributes = strings.Replace(string(xb), "\"", "\\\"", -1)
	}

	query = fmt.Sprintf(query, r.TypeID, 1, r.Primary,
		r.CreatedAt.Format(constants.MySQLTimestampFormat),
		r.UpdatedAt.Format(constants.MySQLTimestampFormat),
		r.PubDate.Date.Format(constants.MySQLDateFormat),
		r.PubDate.Year,
		r.PubDate.Month,
		r.PubDate.Day,
		r.Name, r.Description, keywords,
		r.ResourceURL, r.ShortURL, r.ThumbnailURL, attributes)

	res, err := datastore.MySQL.Session.Exec(query)
	if err != nil {
		fmt.Println("Error with query: \n", query, "\n", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	// Set the ID now we have it, for any subsequent convenience
	r.ID = int(id)

	// A little useful output for the logs...
	fmt.Println("Added a new resource with ID", r.ID)

	// Make sure we set the shortURL - this sets the value in the MySQL field: ol_resource.short_url
	// *AFTER* it creates the record in MongoDB... otherwise, the short link will be activated but will
	// not be redirected.
	err = r.SetShortURL()
	if err != nil {
		fmt.Println("setShortURL error:", err)
	}

	return r.ID, nil
}

// Update performs an update operation on a Resource record in MySQL
func (r *Resource) Update(id int) error {

	// Only difference with Save() query is that created_at is not included
	query := "UPDATE ol_resource SET ol_resource_type_id = %v, active = %v, `primary` = %v," +
		`updated_at = "%v", presented_on = "%v", presented_year = "%v", presented_month = "%v", presented_date = "%v",
		name = "%v", description= "%v", keywords= "%v",
		resource_url = "%v", short_url = "%v", thumbnail_url = "%v"
		WHERE id = %v`

	// Create comma separated keyword list for MySQL field from r.Keywords []string
	keywords := strings.Join(r.Keywords, ",")

	query = fmt.Sprintf(query, r.TypeID, 1, r.Primary,
		r.UpdatedAt.Format(constants.MySQLTimestampFormat),
		r.PubDate.Date.Format(constants.MySQLDateFormat),
		r.PubDate.Year,
		r.PubDate.Month,
		r.PubDate.Day,
		r.Name, r.Description, keywords,
		r.ResourceURL, r.ShortURL, r.ThumbnailURL,
		id)

	_, err := datastore.MySQL.Session.Exec(query)
	if err != nil {
		fmt.Println("Query error:")
		return err
	}

	return nil
}

// SetShortURL sets the short_url fields in MySQL ol_resource.short_url. This is here for convenience so that when
// a Resource is added via the API, and we don't yet have an ID, we can set the short url without making another API call.
// This is a bit "hackish", however the current MappCPD application will send a user to the short url if it exists in the db
// -e, if there is a value in ol_resource.short_url. So BEFORE we set this value we need to make sure a record exists in the Links
// collection in MongoDB as this is what the short url service (linkr) refers to when doing short link redirecting
func (r *Resource) SetShortURL() error {

	// Don't set a short url if there is no long url as a target!
	if r.ResourceURL == "" {
		return errors.New("Cannot set a short url when there is no long url!")
	}

	// Check the Links collection for this Resource (URL)
	q := bson.M{"longUrl": r.ResourceURL}
	l, err := DocLinksOne(q)
	if err != nil {
		switch err {
		case mgo.ErrNotFound:
			fmt.Println("SetShortURL no Links doc so will create one...")
			l.CreatedAt = time.Now()
			l.UpdatedAt = time.Now()
			l.ShortUrl = "r" + strconv.Itoa(r.ID)
			l.LongUrl = r.ResourceURL
			l.Title = r.Name
			err := l.DocSave()
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
	_, err = datastore.MySQL.Session.Exec(query)
	if err != nil {
		fmt.Println("SQL error with query: ", query, " -", err)
		return err
	}

	return nil
}

// DuplicateResourceURL checks the MySQL db for the existence of a resource_url, and returns the FIRST id or 0
// This is here to help prevent the addition of duplicate resources. In the .Save() func if a duplicate is found,
// that is, a duplicate resource_url, .Update() will be run instead. This means the id is required so we know which record
// to update. It is possible there could be more than one duplicate so, for now, return the FIRST duplicate with QueryRow()
func DuplicateResourceURL(url string) int {

	var c int

	// Don't allow this to match an empty url string
	if url == "" {
		return c
	}

	query := "SELECT id FROM ol_resource WHERE active = 1 AND resource_url = ?"
	datastore.MySQL.Session.QueryRow(query, url).Scan(&c)

	return c
}

func ResourceDeepEqual(r, r2 *Resource) bool {

	if r.Name != r2.Name {
		//fmt.Println("")
		//fmt.Println("Name not equal")
		//fmt.Println(r.Name, "&", r2.Name)
		return false
	}
	if r.Description != r2.Description {
		//fmt.Println("")
		//fmt.Println("Descriptions not equal")
		//fmt.Println(r.Description, "&", r2.Description)
		return false
	}
	if !reflect.DeepEqual(r.PubDate, r2.PubDate) {
		//fmt.Println("")
		//fmt.Println("PubDate not equal")
		//fmt.Println(r.PubDate, "&", r2.PubDate)
		return false
	}
	if !reflect.DeepEqual(r.Keywords, r2.Keywords) {
		//fmt.Println("")
		//fmt.Println("Keywords not equal...")
		//for i := range r.Keywords {
		//	if r.Keywords[i] != r2.Keywords[i] {
		//		fmt.Println(r.Keywords[i], " NOT EQUAL TO", r2.Keywords[i])
		//	}
		//}
		//fmt.Println(r.Keywords, "&", r2.Keywords)
		return false
	}

	return true
}
