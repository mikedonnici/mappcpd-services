package cpd_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/cpd"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

var db = testdata.NewDataStore()
var helper = testdata.NewHelper()

func TestCPD(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("CPD", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testCPDByID", testCPDByID)
		t.Run("testCPDByMemberID", testCPDByMemberID)
		t.Run("testCPDQuery", testCPDQuery)
		t.Run("testAddCPD", testAddCPD)
		t.Run("testUpdateCPD", testUpdateCPD)
		t.Run("testDuplicateOf", testDuplicateOf)
		t.Run("testDelete", testDelete)
	})
}

func setup() (datastore.Datastore, func()) {
	var db = testdata.NewDataStore()
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

func testCPDByID(t *testing.T) {

	cases := []struct {
		arg  int    // id
		want string // description
	}{
		{1, "BJJ like Bruno Malfacine"},
		{2, "Ate sausages and eggs"},
		{3, "Baked bread"},
	}

	for _, c := range cases {
		cpd, err := cpd.ByID(ds, c.arg)
		if err != nil {
			t.Fatalf("cpd.ByID(%d) err = %s", c.arg, err)
		}
		got := cpd.Description
		if got != c.want {
			t.Errorf("cpd.ByID(%d).Description = %q, want %q", c.arg, got, c.want)
		}
	}
}

func testCPDByMemberID(t *testing.T) {
	arg := 1 // member id
	xc, err := cpd.ByMemberID(ds, arg)
	if err != nil {
		t.Fatalf("cpd.ByMemberID(%d) err = %s", arg, err)
	}
	got := len(xc)
	want := 3
	if got != want {
		t.Fatalf("cpd.ByMemberID(%d) count = %d, want = %d", arg, got, want)
	}
}

func testCPDQuery(t *testing.T) {
	xc, err := cpd.Query(ds, "WHERE cma.description LIKE '%Bruno%'")
	if err != nil {
		t.Fatalf("cpd.Query() err = %s", err)
	}
	got := len(xc)
	want := 1
	if got != want {
		t.Fatalf("cpd.Query() count = %d, want = %d", got, want)
	}
}

func testAddCPD(t *testing.T) {
	c := cpd.Input{
		MemberID:    1,
		ActivityID:  24,
		TypeID:      25,
		Date:        "2018-05-07",
		Quantity:    2.25,
		Description: "I added this record",
		Evidence:    false,
	}
	id, err := cpd.Add(ds, c)
	if err != nil {
		t.Fatalf("cpd.Add() err = %s", err)
	}

	// fetch the newly added record, and verify the description
	r, err := cpd.ByID(ds, id)
	if err != nil {
		t.Fatalf("cpd.ByID(%d) err = %s", c.ID, err)
	}
	got := c.Description
	want := r.Description
	if got != want {
		t.Fatalf("cpd.ByID(%d).Description = %q, want %q", c.ID, got, want)
	}
}

func testUpdateCPD(t *testing.T) {
	c := cpd.Input{
		ID:          2,
		MemberID:    1,
		ActivityID:  24,
		TypeID:      25,
		Date:        "2018-05-07",
		Quantity:    2.25,
		Description: "The description was updated",
		Evidence:    false,
	}
	err := cpd.Update(ds, c)
	if err != nil {
		t.Fatalf("cpd.Update() err = %s", err)
	}

	// fetch updated record, and verify the description
	r, err := cpd.ByID(ds, c.ID)
	if err != nil {
		t.Fatalf("cpd.ByID(%d) err = %s", c.ID, err)
	}
	got := c.Description
	want := r.Description
	if got != want {
		t.Fatalf("cpd.ByID(%d).Description = %q, want %q", c.ID, got, want)
	}
}

func testDuplicateOf(t *testing.T) {

	// fetch cpd record
	arg := 1 // cpd id
	a, err := cpd.ByID(ds, arg)
	if err != nil {
		t.Fatalf("cpd.ByID(%d) err = %s", arg, err)
	}

	// create a duplicate
	i := cpd.Input{
		MemberID:    a.MemberID,
		ActivityID:  a.Activity.ID,
		TypeID:      a.Type.ID,
		Date:        a.Date,
		Description: a.Description,
		Evidence:    a.Evidence,
		UnitCredit:  a.CreditData.UnitCredit,
		Quantity:    a.Credit,
	}

	// got should be the duplicate id, that is, same as arg
	got, err := cpd.DuplicateOf(ds, i)
	if err != nil {
		t.Fatalf("cpd.DuplicateOf() err = %s", err)
	}
	want := arg
	if got != want {
		t.Errorf("cpd.DuplicateOf() = %d, want %d", got, want)
	}
}

func testDelete(t *testing.T) {

	// get a count before deleting
	xc, err := cpd.Query(ds, "")
	if err != nil {
		t.Fatalf("cpd.Query() err = %s", err)
	}
	countBefore := len(xc)

	// delete one cpd record
	memberID := 1
	cpdID := 3
	err = cpd.Delete(ds, memberID, cpdID)
	if err != nil {
		t.Fatalf("cpd.Delete() err = %s", err)
	}

	// get the count after deleting
	xc, err = cpd.Query(ds, "")
	if err != nil {
		t.Fatalf("cpd.Query() err = %s", err)
	}
	countAfter := len(xc)

	want := countBefore - 1
	got := countAfter
	if got != want {
		t.Errorf("cpd.Query() count = %d, want %d", got, want)
	}
}
