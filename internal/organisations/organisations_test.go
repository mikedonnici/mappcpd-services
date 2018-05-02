package organisations_test

import (
	"testing"

	"github.com/mappcpd/web-services/testdata"
)

const success = "\u2713"
const failure = "\u2717"

var db = testdata.NewTestDB()

func TestPingDatabase(t *testing.T) {
	err := db.MySQL.Session.Ping()
	if err != nil {
		t.Fatal("Could not ping database")
	}
}

func TestTableNames(t *testing.T) {

}

func fail(t *testing.T) {
	db.TearDown()
	t.Fatal("failed")
}

