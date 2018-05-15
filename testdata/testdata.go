package testdata

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-uuid"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/nleof/goyesql"
	"github.com/pkg/errors"
)

// Hard coded for local dev and Travis CI
const MySQLDSN = "root:password@tcp(localhost:3306)/"
const MongoDSN = "mongodb://localhost/mapp_demo"

var path = os.Getenv("GOPATH") + "/src/github.com/mappcpd/web-services/testdata/"
var schemaQueries = goyesql.MustParseFile(path + "schema.sql")
var tableQueries = goyesql.MustParseFile(path + "tables.sql")
var dataQueries = goyesql.MustParseFile(path + "data.sql")

type TestDB struct {
	Name  string
	Store datastore.Datastore
}

// NewDataStore returns a pointer to a TestDB
func NewDataStore() *TestDB {
	s, _ := uuid.GenerateUUID()
	t := TestDB{
		Name: fmt.Sprintf("%v_test", s[0:7]),
		Store: datastore.Datastore{MySQL: datastore.MySQLConnection{
				DSN:  MySQLDSN,
				Desc: "test database",
			},
		},
	}
	return &t
}

// Setup creates and populates the test database
func (t *TestDB) Setup() error {

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
		t.TearDown()
		return errors.Wrap(err, "Error connecting to the test database")
	}

	for _, q := range tableQueries {
		query = fmt.Sprintf(q, t.Name)
		_, err = t.Store.MySQL.Session.Exec(query)
		if err != nil {
			t.TearDown()
			return errors.Wrap(err, "Error creating tables")
		}
	}

	for _, q := range dataQueries {
		query = fmt.Sprintf(q, t.Name)
		_, err = t.Store.MySQL.Session.Exec(query)
		if err != nil {
			t.TearDown()
			return errors.Wrap(err, "Error inserting data - "+query)
		}
	}

	return nil
}

func (t *TestDB) TearDown() error {
	query := fmt.Sprintf(schemaQueries["drop-test-schema"], t.Name)
	_, err := t.Store.MySQL.Session.Exec(query)
	if err != nil {
		return errors.Wrap(err, "Error deleting test schema")
	}
	return nil
}
