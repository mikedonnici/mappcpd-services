package organisation_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/cardiacsociety/web-services/internal/organisation"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

func TestOrganisation(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("organisation", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testOrganisationByID", testOrganisationByID)
		t.Run("testOrganisationDeepEqual", testOrganisationDeepEqual)
		t.Run("testOrganisationCount", testOrganisationCount)
		t.Run("testChildOrganisationCount", testChildOrganisationCount)
		t.Run("testOrganisationByTypeID", testOrganisationByTypeID)
	})
}

func setup() (datastore.Datastore, func()) {
	var db = testdata.NewDataStore()
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalf("SetupMySQL() err = %s", err)
	}
	return db.Store, func() {
		db.TearDownMySQL()
	}
}

func testPingDatabase(t *testing.T) {
	err := ds.MySQL.Session.Ping()
	if err != nil {
		t.Fatalf("Ping() err = %s", err)
	}
}

func testOrganisationByID(t *testing.T) {
	org, err := organisation.ByID(ds, 1)
	if err != nil {
		t.Fatalf("organisation.ByID() err = %s", err)
	}
	got := org.Name
	want := "ABC Organisation"
	if got != want {
		t.Errorf("organisation.ByID() Name = %q, want %q", got, want)
	}
}

func testOrganisationDeepEqual(t *testing.T) {

	exp := organisation.Organisation{
		ID:   1,
		Code: "ABC",
		Name: "ABC Organisation",
		Groups: []organisation.Organisation{
			{ID: 3, Code: "ABC-1", Name: "ABC Sub1"},
			{ID: 4, Code: "ABC-2", Name: "ABC Sub2"},
			{ID: 5, Code: "ABC-3", Name: "ABC Sub3"},
		},
	}

	o, err := organisation.ByID(ds, 1)
	if err != nil {
		t.Fatalf("organisation.ByID() err = %s", err)
	}

	got := reflect.DeepEqual(exp, o)
	want := true
	if got != want {
		t.Errorf("reflect.DeepEqual() = %v, want %v", got, want)
	}
}

// Test data has 2 parent organisations
func testOrganisationCount(t *testing.T) {
	xo, err := organisation.All(ds)
	if err != nil {
		t.Fatalf("organisation.All() err = %s", err)
	}
	got := len(xo)
	want := 2
	if got != want {
		t.Errorf("organisation.All() count = %d, want %d", got, want)
	}
}

// Test data has 3 child organisations belonging to parent id 1
func testChildOrganisationCount(t *testing.T) {
	o, err := organisation.ByID(ds, 1)
	if err != nil {
		t.Fatalf("organisation.ByID() err = %s", err)
	}
	got := len(o.Groups)
	want := 3
	if got != want {
		t.Errorf("organisation.ByID().Groups count = %d, want %d", got, want)
	}
}

// test organisations by type id
func testOrganisationByTypeID(t *testing.T) {
	cases := []struct {
		arg  int // type id
		want int // count
	}{
		{arg: 1, want: 2},
		{arg: 2, want: 1},
		{arg: 4, want: 1},
	}

	for _, c := range cases {
		xo, err := organisation.ByTypeID(ds, c.arg)
		if err != nil {
			t.Fatalf("organisation.ByTypeID(%d) err = %s", c.arg, err)
		}
		got := len(xo)
		if got != c.want {
			t.Errorf("organisation.ByTypeID(%d) count = %d, want %d", c.arg, got, c.want)
		}
	}
}
