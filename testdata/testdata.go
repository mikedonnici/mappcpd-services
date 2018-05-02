package testdata

import (
	"fmt"
	"time"

	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/nleof/goyesql"
	"github.com/pkg/errors"
)

const DSN = "dev:password@tcp(localhost:3306)/"

var schemaQueries = goyesql.MustParseFile("../../testdata/schema.sql")
var tableQueries = goyesql.MustParseFile("../../testdata/tables.sql")
var dataQueries = goyesql.MustParseFile("../../testdata/data.sql")

type TestDB struct {
	Name  string
	MySQL datastore.MySQLConnection
}

// NewTestDB returns a pointer to a TestDB
func NewTestDB() *TestDB {
	t := TestDB{Name: fmt.Sprintf("%v_test", time.Now().Unix())}
	return &t
}

// Setup creates and populates the test database
func (t *TestDB) Setup() error {

	err := t.MySQL.ConnectSource(DSN)
	if err != nil {
		return errors.Wrap(err, "Error establishing session with MySQL")
	}

	query := fmt.Sprintf(schemaQueries["create-test-schema"], t.Name)
	_, err = t.MySQL.Session.Exec(query)
	if err != nil {
		return errors.Wrap(err, "Error creating test schema")
	}

	err = t.MySQL.ConnectSource(DSN + t.Name)
	if err != nil {
		return errors.Wrap(err, "Error connecting to the test database")
	}

	for _, q := range tableQueries {
		query = fmt.Sprintf(q, t.Name)
		_, err = t.MySQL.Session.Exec(query)
		if err != nil {
			return errors.Wrap(err, "Error creating tables")
		}
	}

	for _, q := range dataQueries {
		query = fmt.Sprintf(q, t.Name)
		_, err = t.MySQL.Session.Exec(query)
		if err != nil {
			return errors.Wrap(err, "Error inserting data")
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
