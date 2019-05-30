package tag_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/tag"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

func TestTag(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("tag", func(t *testing.T) {
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
		t.Fatalf("Ping() err = %s", err)
	}
}

// fetch the list of tags
func testAll(t *testing.T) {
	xs, err := tag.All(ds)
	if err != nil {
		t.Fatalf("Tag.All() err = %s", err)
	}
	got := len(xs)
	want := 5 // only 5 in test data
	if got != want {
		t.Errorf("Tag.All() count = %d, want %d", got, want)
	}
}
