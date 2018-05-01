package datastore

import (
	"os"

	"gopkg.in/mgo.v2"
)

// MongoDBConnection represents a connection to a MongoDB server
// and includes convenience methods for accessing each collection
type MongoDBConnection struct {
	url     string
	db      string
	Source  string
	session *mgo.Session
}

// MongoQuery is used to map query fields in a request to mgo functions
type MongoQuery struct {
	Find   map[string]interface{} `json:"find"`
	Select map[string]interface{} `json:"select"`
	Limit  int                    `json:"limit"`
	Sort   string                 `json:"sort"`
}

// Do executes the query on a collection (c), and scans the results (r)
func (mq MongoQuery) Do(c *mgo.Collection, r *[]interface{}) error {

	// Start to build the Query...
	q := c.Find(mq.Find).Select(mq.Select).Limit(mq.Limit)

	// ... only add sort if there is a value there
	if len(mq.Sort) > 0 {
		q.Sort(mq.Sort)
	}

	// .All runs the query, scans into r and returns an error, if present
	return q.All(r)
}

// Connect to to MongoDB, returns an error if it fails
func (m *MongoDBConnection) Connect() error {

	m.url = os.Getenv("MAPPCPD_MONGO_URL")
	m.db = os.Getenv("MAPPCPD_MONGO_DBNAME")
	m.Source = os.Getenv("MAPPCPD_MONGO_DESC")

	var err error
	m.session, err = mgo.Dial(m.url)
	if err != nil {
		return err
	}

	return nil
}

// MembersCollection returns a pointer to the Members collection
func (m *MongoDBConnection) MembersCollection() (*mgo.Collection, error) {

	return m.session.DB(m.db).C("Members"), nil
}

// ActivitiesCol returns a pointer to the Activities collection
func (m *MongoDBConnection) ActivitiesCol() (*mgo.Collection, error) {

	return m.session.DB(m.db).C("Activities"), nil
}

// ResourcesCollection returns a pointer to the Resources collection
func (m *MongoDBConnection) ResourcesCollection() (*mgo.Collection, error) {

	return m.session.DB(m.db).C("Resources"), nil
}

// ModulesCollection returns a pointer to the Modules collection
func (m *MongoDBConnection) ModulesCollection() (*mgo.Collection, error) {

	return m.session.DB(m.db).C("Modules"), nil
}

// LinksCol returns a pointer to the Links collection
func (m *MongoDBConnection) LinksCol() (*mgo.Collection, error) {

	return m.session.DB(m.db).C("Links"), nil
}

// RecurringCol returns a pointer to the Recurring collection
func (m *MongoDBConnection) RecurringCol() (*mgo.Collection, error) {

	return m.session.DB(m.db).C("Recurring"), nil
}

// Close terminates the session
func (m *MongoDBConnection) Close() {
	m.session.Close()
}
