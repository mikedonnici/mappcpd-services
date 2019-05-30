package activity_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/activity"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

func TestActivity(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("activity", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testActivityCount", testActivityCount)
		t.Run("testActivityTypesCount", testActivityTypesCount)
		t.Run("testActivityByID", testActivityByID)
		t.Run("testActivityByTypeID", testActivityByTypeID)
	})
}

func setup() (datastore.Datastore, func()) {
	var db = testdata.NewDataStore()
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalf("db.SetupMySQL() err = %s", err)
	}
	return db.Store, func() {
		err := db.TearDownMySQL()
		if err != nil {
			log.Fatalf("db.TearDownMySQL() err = %s", err)
		}
	}
}

func testPingDatabase(t *testing.T) {
	err := ds.MySQL.Session.Ping()
	if err != nil {
		t.Fatalf("Ping() err = %s", err)
	}
}

func testActivityCount(t *testing.T) {
	xa, err := activity.All(ds)
	if err != nil {
		t.Fatalf("activity.All() err = %s", err)
	}
	got := len(xa)
	want := 5
	if got != want {
		t.Errorf("activity.All() = %d, want %d", got, want)
	}
}

func testActivityTypesCount(t *testing.T) {
	cases := []struct {
		arg  int
		want int
	}{
		{1, 0},
		{3, 0},
		{20, 9},
		{24, 5},
	}
	for _, c := range cases {
		xa, err := activity.Types(ds, c.arg)
		if err != nil {
			t.Fatalf("activity.Types() err = %s", err)
		}
		got := len(xa)
		if got != c.want {
			t.Errorf("activity.Types() count = %d, want %d", got, c.want)
		}
	}
}

func testActivityByID(t *testing.T) {
	cases := []struct {
		arg  int
		want string
	}{
		{4, "Presentation"},
		{23, "Group Learning"},
	}
	for _, c := range cases {
		a, err := activity.ByID(ds, c.arg)
		if err != nil {
			t.Fatalf("activity.ByID() err = %s", err)
		}
		got := a.Name
		if got != c.want {
			t.Errorf("activity.ByID() Activity.Name = %q, want %q", got, c.want)
		}
	}
}

// fetch activity by type id, echeck the correct activity id was returned
func testActivityByTypeID(t *testing.T) {
	cases := []struct {
		arg  int
		want int
	}{
		{2, 20},
		{13, 21},
		{28, 23},
		{36, 24},
	}
	for _, c := range cases {
		a, err := activity.ByTypeID(ds, c.arg)
		if err != nil {
			t.Fatalf("activity.ByTypeID() err = %s", err)
		}
		got := a.ID
		if got != c.want {
			t.Errorf("activity.ByTypeID() Activity.ID = %d, want %d", got, c.want)
		}
	}
}
