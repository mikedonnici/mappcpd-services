package organisation_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/mikedonnici/mappcpd-services/internal/organisation"
	"github.com/mikedonnici/mappcpd-services/testdata"
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

func TestPingDatabase(t *testing.T) {
	err := db.Store.MySQL.Session.Ping()
	if err != nil {
		t.Fatal("Could not ping database")
	}
}

func TestOrganisationByID(t *testing.T) {
	org, err := organisation.ByID(db.Store, 1)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, "ABC Organisation", org.Name)
}

func TestOrganisationDeepEqual(t *testing.T) {

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

	o, err := organisation.ByID(db.Store, 1)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	res := reflect.DeepEqual(exp, o)
	helper.Result(t, true, res)
}

// Test data has 2 parent organisations
func TestOrganisationCount(t *testing.T) {
	xo, err := organisation.All(db.Store)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 2, len(xo))
}

// Test data has 3 child organisations belonging to parent id 1
func TestChildOrganisationCount(t *testing.T) {
	o, err := organisation.ByID(db.Store, 1)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 3, len(o.Groups))
}
