package organisations_test

import (
	"reflect"
	"testing"
	"log"

	"github.com/mappcpd/web-services/internal/organisations"
	"github.com/mappcpd/web-services/testdata"
)

var db = testdata.NewTestDB()
var helper = testdata.NewHelper()

// todo does not exit on first failed test?
func TestMain(m *testing.M) {
	err := db.Setup()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.TearDown()

	m.Run()
}

func TestPingDatabase(t *testing.T) {
	err := db.MySQL.Session.Ping()
	if err != nil {
		t.Fatal("Could not ping database")
	}
}

func TestOrganisationByID(t *testing.T) {
	org, err := organisations.OrganisationByIDStore(1, db.MySQL)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, "ABC Organisation", org.Name)
}

func TestOrganisationDeepEqual(t *testing.T) {

	exp := organisations.Organisation{
		ID:   1,
		Name: "ABC Organisation",
		Code: "ABC",
	}

	org, err := organisations.OrganisationByIDStore(1, db.MySQL)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	res := reflect.DeepEqual(exp, org)
	helper.Result(t, true, res)
}

// Test data has 2 parent organisations
func TestOrganisationListCount(t *testing.T) {
	l, err := organisations.OrganisationsListStore(db.MySQL)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 2, len(l))
}

// Test data has 3 child organisations belonging to parent id 1
func TestChildOrganisationsListCount(t *testing.T) {
	l, err := organisations.ChildOrganisationsStore(1, db.MySQL)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 3, len(l))
}