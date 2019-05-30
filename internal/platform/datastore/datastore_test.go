package datastore_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

var db = testdata.NewDataStore()
var helper = testdata.NewHelper()

func TestMain(m *testing.M) {
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.TearDownMySQL()

	m.Run()
}

func TestNewDatastoreMySQL(t *testing.T) {
	ds := datastore.New()
	ds.MySQL.DSN = testdata.MySQLDSN
	ds.MySQL.Desc = "MySQL test database"
	err := ds.ConnectMySQL()
	helper.Result(t, nil, err)
}

func TestNewDatastoreMongoDB(t *testing.T) {
	ds := datastore.New()
	ds.MongoDB.DSN = testdata.MongoDSN
	ds.MongoDB.DBName = "test"
	ds.MongoDB.Desc = "MongoDB test database"
	err := ds.ConnectMongoDB()
	helper.Result(t, nil, err)
}
