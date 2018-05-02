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
	helper.PrintResult(t, "Allied Health Council", org.Name)
}

func TestOrganisationDeepEqual(t *testing.T) {

	exp := organisations.Organisation{
		ID:   1,
		Name: "Allied Health Council",
		Code: "CL_AH",
	}

	org, err := organisations.OrganisationByIDStore(1, db.MySQL)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	res := reflect.DeepEqual(exp, org)
	helper.PrintResult(t, true, res)
}
