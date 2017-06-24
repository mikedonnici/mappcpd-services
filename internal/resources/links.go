package resources

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/utility"
)

type Link struct {
	_id            bson.ObjectId `json:"_id" bson:"_id"`
	CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt" bson:"updatedAt"`
	ShortUrl       string        `json:"shortUrl" bson:"shortUrl"`
	LongUrl        string        `json:"longUrl" bson:"longUrl"`
	Title          string        `json:"title" bson:"title"`
	Clicks         int           `json:"clicks" bson:"clicks"`
	LastStatusCode int           `json:"lastStatusCode" bson:"lastStatusCode"`
}

func (l *Link) DocSave() error {

	// Get pointer to the Links collection
	lc, err := datastore.MongoDB.LinksCol()
	if err != nil {
		log.Printf("Error getting pointer to Links collection: %s\n", err.Error())
		return err
	}

	// Selector for Upsert - no MySQL id here so use the long url, could use _id
	s := map[string]string{"longUrl": l.LongUrl}
	// Upsert
	_, err = lc.Upsert(s, &l)
	if err != nil {
		log.Printf("Error updating Links doc: %s\n", err.Error())
		return err
	}

	return nil
}

// DocLinksOne returns a single link doc, unmarshaled into the proper struct.
// Note this DOES return an error when nothing is found
func DocLinksOne(q map[string]interface{}) (Link, error) {

	l := Link{}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt"})

	lc, err := datastore.MongoDB.LinksCol()
	if err != nil {
		return l, err
	}
	err = lc.Find(q).One(&l)
	if err != nil {
		return l, err
	}

	return l, nil
}
