package position_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/position"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

func TestAll(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("invoice", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testByID", testByID)
		t.Run("testByIDs", testByIDs)
		t.Run("testExcelReport", testExcelReport)
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

// fetch a member position by id, verify the member id and position name
func testByID(t *testing.T) {
	cases := []struct {
		arg          int // pos id
		wantMemberID int
		wantPosName  string
	}{
		{
			arg:          1,
			wantMemberID: 1,
			wantPosName:  "Affiliate",
		},
		{
			arg:          2,
			wantMemberID: 1,
			wantPosName:  "Member",
		},
	}

	for _, c := range cases {
		p, err := position.ByID(ds, c.arg)
		if err != nil {
			t.Errorf("position.ByID(%d) err = %s", c.arg, err)
		}
		gotMemberID := p.MemberID
		gotPosName := p.Name
		if gotMemberID != c.wantMemberID {
			t.Errorf("position.ByID(%d) MemberID = %d, want %d", c.arg, gotMemberID, c.wantMemberID)
		}
		if gotPosName != c.wantPosName {
			t.Errorf("position.ByID(%d) Name = %q, want %q", c.arg, gotPosName, c.wantPosName)
		}
	}
}

// fetch multiple positions, verify position names
func testByIDs(t *testing.T) {
	cases := []struct {
		arg  []int    // pos ids
		want []string // pos names
	}{
		{[]int{1, 2}, []string{"Affiliate", "Member"}},
	}

	for _, c := range cases {
		xp, err := position.ByIDs(ds, c.arg)
		if err != nil {
			t.Errorf("position.ByIDs() err = %s", err)
		}
		got := []string{}
		for _, p := range xp {
			got = append(got, p.Name)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("position.ByIDs() Names = %v, want %v", got, c.want)
		}
	}
}

// fetch some test data and ensure excel report is not returning an error
func testExcelReport(t *testing.T) {

	ids := []int{1, 2, 3} // position records
	want := 4             // expect 4 rows - heading and 3 records

	xp, err := position.ByIDs(ds, ids)
	if err != nil {
		t.Fatalf("position.ByIDs() err = %s", err)
	}
	f, err := position.ExcelReport(ds, xp)
	if err != nil {
		t.Fatalf("position.ExcelReport() err = %s", err)
	}
	
	rows := f.GetRows(f.GetSheetName(f.GetActiveSheetIndex())) // rows is [][]string
	got := len(rows)
	if got != want {
		t.Errorf("GetRows() row count = %d, want %d", got, want)
	}
}
