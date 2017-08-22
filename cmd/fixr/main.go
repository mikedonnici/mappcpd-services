package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"database/sql"

	"github.com/mappcpd/web-services/internal/resources"
	"github.com/34South/envr"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/pkg/errors"
)

var db *sql.DB
var ddb *mgo.Session

type link struct {
	id       int
	title    string
	shortPath string
	shortURL string
	longURL  string
}

func init() {
	envr.New("fixrEnv", []string{
		"MAPPCPD_MYSQL_URL",
		"MAPPCPD_MONGO_URL",
		"MAPPCPD_MONGO_DBNAME",
		"MAPPCPD_SHORT_LINK_URL",
		"MAPPCPD_SHORT_LINK_PREFIX",
	}).Auto()
}

func main() {

	fmt.Println("Fixing...")

	var err error // no shadowing!
	db, err = sql.Open("mysql", os.Getenv("MAPPCPD_MYSQL_URL"))
	if err != nil {
		log.Fatalln("Could not connect to MySQL server:", os.Getenv("MAPPCPD_MYSQL_URL"))
	}

	ddb, err = mgo.Dial(os.Getenv("MAPPCPD_MONGO_URL"))
	if err != nil {
		log.Fatal("Failed to establish a session with Mongo server - " + err.Error())
	}

	// Select resources that start with 'http%' so don't break relative URLs
	query := "SELECT id, name, short_url, resource_url FROM ol_resource " +
		"WHERE `active` = 1 AND `primary` = 1 AND resource_url LIKE 'http%' "

	rows, err := db.Query(query)
	for rows.Next() {

		l := link{}

		err := rows.Scan(&l.id, &l.title, &l.shortURL, &l.longURL)
		if err != nil {
			msg := fmt.Sprintf("Error scanning row with id %v", l.id, " - skipping this record")
			fmt.Println(msg)
			continue
		}

		// Work out what we expect, or need to set up for a short link
		// Custom short URL based on id of resource with no padding, eg r12, r3435
		l.shortPath = fmt.Sprintf("%v%v", os.Getenv("MAPPCPD_SHORT_LINK_PREFIX"), l.id)
		// The short_url should be
		expectedShortURL := fmt.Sprintf("%v/%v", os.Getenv("MAPPCPD_SHORT_LINK_URL"), l.shortPath)
		fmt.Println("Short path is:", l.shortPath)
		fmt.Println("Short URL should be:", expectedShortURL)


		// There are two scenarios:
		// 1. short_url already has the expected value in the primary store
		// Look for any differences between the primary record and the Links doc, if found then sync changes.
		//
		// 2. Missing, invalid or does not match expected based on config
		// Need to set the full short URL in the primary record, and then sync as above

		// Set in primary record, if not the expected value...
		if l.shortURL != expectedShortURL {
			fmt.Println("Short url is missing / invalid in primary db")
			l.shortURL = expectedShortURL
			err := setShortURL(l.id, l.shortURL)
			if err != nil {
				fmt.Println(errors.Cause(err))
			}
		}

		// Check for changes and sync if required...
		err = checkSync(l)
		if err != nil {
			fmt.Println(errors.Cause(err))
		}
	}
}

// setShortURL sets the value of the short_url field in the ol_resource (primary) record
func setShortURL(id int, shortURL string ) error {

		query := `UPDATE ol_resource SET short_url = "%v" WHERE id = %v LIMIT 1`
		query = fmt.Sprintf(query, shortURL, id)
		fmt.Println("Update SQL:", query)
		_, err := db.Exec(query)

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
func getLinkDoc(shortPath string) (resources.Link, error) {

	var l resources.Link

	// Always set this as it is used as the KEY for Links docs
	l.ShortUrl = shortPath

	// Query the database
	c := ddb.DB(os.Getenv("MAPPCPD_MONGO_DBNAME")).C("Links")
	s := bson.M{"shortUrl": shortPath}
	err := c.Find(s).One(&l)

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
func sync(ld resources.Link, shortPath string) error {

	c := ddb.DB(os.Getenv("MAPPCPD_MONGO_DBNAME")).C("Links")
	s := bson.M{"shortUrl": shortPath}
	u := bson.M{"$set": ld}
	_, err := c.Upsert(s, u)
	if err != nil {
		return errors.Wrap(err, "upsert failed")
	}

	return nil
}
