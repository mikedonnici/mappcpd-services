package testdata

import (
	"fmt"

	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/nleof/goyesql"
	"github.com/pkg/errors"
	"github.com/hashicorp/go-uuid"
)

// Hard coded for local dev and Travis CI
const MySQLDSN = "root:password@tcp(localhost:3306)/"
const MongoDSN = "mongodb://localhost/mapp_demo"

var schemaQueries = goyesql.MustParseFile("testdata/schema.sql")
var tableQueries = goyesql.MustParseFile("testdata/tables.sql")
var dataQueries = goyesql.MustParseFile("testdata/data.sql")

type TestDB struct {
	Name  string
	MySQL datastore.MySQLConnection
}

// NewTestDB returns a pointer to a TestDB
func NewTestDB() *TestDB {
	s, _ := uuid.GenerateUUID()
	t := TestDB{
		Name: fmt.Sprintf("%v_test", s[0:7]),
		MySQL: datastore.MySQLConnection{
			DSN:  MySQLDSN,
			Desc: "test database",
		},
	}
	return &t
}

// Setup creates and populates the test database
func (t *TestDB) Setup() error {

	err := t.MySQL.Connect()
	if err != nil {
		return errors.Wrap(err, "Error establishing session with MySQL")
	}

	query := fmt.Sprintf(schemaQueries["create-test-schema"], t.Name)
	_, err = t.MySQL.Session.Exec(query)
	if err != nil {
		return errors.Wrap(err, "Error creating test schema")
	}

	// Update session to connect to new database
	t.MySQL.DSN = t.MySQL.DSN + t.Name
	err = t.MySQL.Connect()
	if err != nil {
		t.TearDown()
		return errors.Wrap(err, "Error connecting to the test database")
	}

	for _, q := range tableQueries {
		query = fmt.Sprintf(q, t.Name)
		_, err = t.MySQL.Session.Exec(query)
		if err != nil {
			t.TearDown()
			return errors.Wrap(err, "Error creating tables")
		}
	}

	for _, q := range dataQueries {
		query = fmt.Sprintf(q, t.Name)
		_, err = t.MySQL.Session.Exec(query)
		if err != nil {
			t.TearDown()
			return errors.Wrap(err, "Error inserting data - "+query)
		}
	}

	return nil
}

func (t *TestDB) TearDown() error {
	query := fmt.Sprintf(schemaQueries["drop-test-schema"], t.Name)
	_, err := t.MySQL.Session.Exec(query)
	if err != nil {
		return errors.Wrap(err, "Error deleting test schema")
	}
	return nil
}
