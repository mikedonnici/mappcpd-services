package testdata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/hashicorp/go-uuid"
	"github.com/nleof/goyesql"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

// Hard coded for local dev and Travis CI
const MySQLDSN = "root:password@tcp(localhost:3306)/"
const MongoDSN = "mongodb://localhost/"

var path = os.Getenv("GOPATH") + "/src/github.com/cardiacsociety/web-services/testdata/"
var schemaQueries = goyesql.MustParseFile(path + "schema.sql")
var tableQueries = goyesql.MustParseFile(path + "tables.sql")
var dataQueries = goyesql.MustParseFile(path + "data.sql")
var memberDocs = path + "members.json"
var resourcesDocs = path + "resources.json"

type TestStore struct {
	Name  string
	Store datastore.Datastore
}

// NewDataStore returns a pointer to a TestStore
func NewDataStore() *TestStore {

	s, _ := uuid.GenerateUUID()
	n := fmt.Sprintf("%v_test", s[0:7])

	t := TestStore{
		Name: n,
		Store: datastore.Datastore{
			MySQL: datastore.MySQLConnection{
				DSN:  MySQLDSN,
				Desc: "test MySQL database",
			},
			MongoDB: datastore.MongoDBConnection{
				DBName: n,
				DSN:    MongoDSN,
				Desc:   "test Mongo database",
			},
		},
	}
	return &t
}

// SetupMySQL creates and populates the test MySQL database
func (t *TestStore) SetupMySQL() error {

	err := t.Store.MySQL.Connect()
	if err != nil {
		return errors.Wrap(err, "Error establishing session with MySQL")
	}

	query := fmt.Sprintf(schemaQueries["create-test-schema"], t.Name)
	_, err = t.Store.MySQL.Session.Exec(query)
	if err != nil {
		return errors.Wrap(err, "Error creating test schema")
	}

	// Update session to connect to new database
	t.Store.MySQL.DSN = t.Store.MySQL.DSN + t.Name
	err = t.Store.MySQL.Connect()
	if err != nil {
		t.TearDownMySQL()
		return errors.Wrap(err, "Error connecting to the test database")
	}

	for _, q := range tableQueries {
		query = fmt.Sprintf(q, t.Name)
		_, err = t.Store.MySQL.Session.Exec(query)
		if err != nil {
			t.TearDownMySQL()
			return errors.Wrap(err, "Error creating tables")
		}
	}

	for _, q := range dataQueries {
		query = fmt.Sprintf(q, t.Name)
		_, err = t.Store.MySQL.Session.Exec(query)
		if err != nil {
			t.TearDownMySQL()
			return errors.Wrap(err, "Error inserting data - "+query)
		}
	}

	return nil
}

// SetupMongoDB creates and populates the test Mongo database
func (t *TestStore) SetupMongoDB() error {

	err := t.Store.MongoDB.Connect()
	if err != nil {
		return errors.Wrap(err, "Error establishing session with MongoDB")
	}

	err = t.Store.MongoDB.Session.Ping()
	if err != nil {
		return errors.Wrap(err, "Error pinging MongoDB")
	}

	// Import member data
	m := bson.M{}
	f, err := ioutil.ReadFile(memberDocs)
	if err != nil {
		return errors.Wrap(err, "Error reading members json file")
	}
	err = json.Unmarshal(f, &m)
	if err != nil {
		return errors.Wrap(err, "Unmarshal error - members")
	}
	err = t.Store.MongoDB.Session.DB(t.Store.MongoDB.DBName).C("Members").Insert(m)
	if err != nil {
		return errors.Wrap(err, "Error inserting member document")
	}

	// Import resources data
	var xr []bson.M
	f, err = ioutil.ReadFile(resourcesDocs)
	if err != nil {
		return errors.Wrap(err, "Error reading resources json file")
	}
	err = json.Unmarshal(f, &xr)
	if err != nil {
		return errors.Wrap(err, "Unmarshal error - resources")
	}

	for _, r := range xr {
		err = t.Store.MongoDB.Session.DB(t.Store.MongoDB.DBName).C("Resources").Insert(r)
		if err != nil {
			return errors.Wrap(err, "Error inserting resource document")
		}
	}

	return nil
}

func (t *TestStore) TearDownMySQL() error {
	query := fmt.Sprintf(schemaQueries["drop-test-schema"], t.Name)
	_, err := t.Store.MySQL.Session.Exec(query)
	if err != nil {
		return errors.Wrap(err, "Error deleting MySQL test database")
	}
	return nil
}

func (t *TestStore) TearDownMongoDB() error {
	err := t.Store.MongoDB.Session.DB(t.Store.MongoDB.DBName).DropDatabase()
	if err != nil {
		return errors.Wrap(err, "Error deleting MongoDB test database")
	}
	return nil
}
