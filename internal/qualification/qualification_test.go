package qualification_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/qualification"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

func TestQualification(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("Qualifications", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testAll", testAll)
	})
}

func setup() (datastore.Datastore, func()) {
	db := testdata.NewDataStore()
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalf("SetupMySQL() err = %s", err)
	}
	return db.Store, func() {
		err := db.TearDownMySQL()
		if err != nil {
			log.Fatalf("TearDownMySQL() err = %s", err)
		}
	}
}

func testPingDatabase(t *testing.T) {
	err := ds.MySQL.Session.Ping()
	if err != nil {
		t.Fatalf("MySQL.Session.Ping() err = %s", err)
	}
}

// fetch the list of qualifications
func testAll(t *testing.T) {
	xq, err := qualification.All(ds)
	if err != nil {
		t.Fatalf("qualification.All() err = %s", err)
	}
	got := len(xq)
	want := 29 // from testdata
	if got != want {
		t.Errorf("qualification.All() count = %d, want %d", got, want)
	}
}
