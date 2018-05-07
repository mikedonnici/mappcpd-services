package cpd_test

import (
	"log"
	"testing"

	"github.com/mappcpd/web-services/internal/cpd"
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

func TestCPDByID(t *testing.T) {

	cases := []struct {
		id   int
		desc string
	}{
		{1, "BJJ like Bruno Malfacine"},
		{2, "Ate sausages and eggs"},
		{3, "Baked bread"},
	}

	for _, c := range cases {
		cpd, err := cpd.ByIDStore(c.id, db.MySQL)
		if err != nil {
			t.Fatalf("Database error: %s", err)
		}
		helper.Result(t, c.desc, cpd.Description)
	}
}

func TestCPDByMemberID(t *testing.T) {
	xcpd, err := cpd.ByMemberIDStore(1, db.MySQL)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 3, len(xcpd))
}

func TestCPDQuery(t *testing.T) {
	xcpd, err := cpd.QueryStore("WHERE cma.description LIKE '%Bruno%'", db.MySQL)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 1, len(xcpd))
}

func TestAddCPD(t *testing.T) {
	c := cpd.Input{
		MemberID: 1,
		ActivityID: 24,
		TypeID: 25,
		Date: "2018-05-07",
		Quantity: 2.25,
		Description: "I added this record",
		Evidence: false,
	}
	id, err := cpd.AddStore(c, db.MySQL)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	// fetch the newly added record, and verify the description
	r, err := cpd.ByIDStore(id, db.MySQL)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	helper.Result(t, c.Description, r.Description)
}

func TestUpdateCPD(t *testing.T) {
	c := cpd.Input{
		ID: 2,
		MemberID: 1,
		ActivityID: 24,
		TypeID: 25,
		Date: "2018-05-07",
		Quantity: 2.25,
		Description: "The description was updated",
		Evidence: false,
	}
	err := cpd.UpdateStore(c, db.MySQL)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	r, err := cpd.ByIDStore(c.ID, db.MySQL)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}

	helper.Result(t, c.Description, r.Description)
}
