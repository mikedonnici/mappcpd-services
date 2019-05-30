package application_test

import (
	"database/sql"
	"log"
	"reflect"
	"testing"

	"github.com/cardiacsociety/web-services/internal/application"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

func TestApplication(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("application", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testByID", testByID)
		t.Run("testByIDs", testByIDs)
		t.Run("testByID_notFound", testByID_notFound)
		t.Run("testByMemberID", testByMemberID)
		t.Run("testByMemberID_notFound", testByMemberID_notFound)
		t.Run("testByNonExistentMemberID", testByNonExistentMemberID)
		t.Run("testQuery", testQuery)
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

// fetch an application by id, verify member id
func testByID(t *testing.T) {

	cases := []struct {
		arg  int // application id
		want int // member id
	}{
		{1, 502},
		{2, 482},
		{3, 488},
	}

	for _, c := range cases {
		got, err := application.ByID(ds, c.arg)
		if err != nil {
			t.Errorf("application.ByID(%d) err = %s", c.arg, err)
		}
		if got.MemberID != c.want {
			t.Errorf("Application.MemberID = %d, want %d", got.MemberID, c.want)
		}
	}
}

// fetch applications by a list of IDs
func testByIDs(t *testing.T) {

	cases := []struct {
		arg  []int // application IDs
		want int   // count
	}{
		{[]int{101}, 0},
		{[]int{1}, 1},
		{[]int{1, 2, 3}, 3},
		{[]int{1, 2, 101}, 2},
	}

	for _, c := range cases {
		xa, err := application.ByIDs(ds, c.arg)
		if err != nil {
			t.Errorf("application.ByIDs(%d) err = %s", c.arg, err)
		}
		got := len(xa)
		if got != c.want {
			t.Errorf("Application.ByIDs count = %d, want %d", got, c.want)
		}
	}
}

// attempt fetch an application by id that does not exist
func testByID_notFound(t *testing.T) {
	arg := 101 // does not exist
	_, err := application.ByID(ds, arg)
	if err == nil {
		t.Errorf("application.ByID(%d) err = nil, want %s", arg, sql.ErrNoRows)
	}
}

// fetch application records by member id
func testByMemberID(t *testing.T) {
	cases := []struct {
		arg  int   // member id
		want []int // application ids
	}{
		{502, []int{1, 6}},
	}

	for _, c := range cases {
		xa, err := application.ByMemberID(ds, c.arg)
		if err != nil {
			t.Errorf("application.ByMemberID(%d) err = %s", c.arg, err)
		}
		var got []int
		for _, a := range xa {
			got = append(got, a.ID)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("application.ByMemberID(%d) = %v, want %v", c.arg, got, c.want)
		}

	}
}

// attempt fetch application records by member id where no applications exist
func testByMemberID_notFound(t *testing.T) {
	arg := 1 // member id 1 has no applications in test
	xa, err := application.ByMemberID(ds, arg)
	if err != nil {
		t.Errorf("application.ByMemberID(%d) err = %s", arg, err)
	}
	got := len(xa)
	want := 0
	if got != want {
		t.Errorf("application.ByMemberID(%d) len = %d, want %d", arg, got, want)
	}
}

// attempt fetch application records by member id that does not exist
func testByNonExistentMemberID(t *testing.T) {
	arg := 101
	xa, err := application.ByMemberID(ds, arg)
	if err != nil {
		t.Errorf("application.ByMemberID(%d) err = %s", arg, err)
	}
	got := len(xa)
	want := 0
	if got != want {
		t.Errorf("application.ByMemberID(%d) len = %d, want %d", arg, got, want)
	}
}

// test generic query function, specify clause and check expected result count
func testQuery(t *testing.T) {
	cases := []struct {
		arg  string
		want int
	}{
		{"", 6},
		{"AND member_id = 488", 1},
		{"AND member_id = 502", 2},
		{"AND member_id = 101", 0},
		{"AND applied_on > '2017-01-01'", 1},
		{"AND ma.id IN (1,2,3)", 3},
	}
	for _, c := range cases {
		xa, err := application.Query(ds, c.arg)
		if err != nil {
			t.Errorf("application.Query() err = %s", err)
		}
		got := len(xa)
		if got != c.want {
			t.Errorf("application.Query() count = %d, want %d", got, c.want)
		}
	}
}

// fetch some test data and ensure excel report is not returning an error
func testExcelReport(t *testing.T) {

	ids := []int{1, 2, 3} // application records
	want := 4             // expect 4 rows - heading and 3 records

	xp, err := application.ByIDs(ds, ids)
	if err != nil {
		t.Fatalf("application.ByIDs() err = %s", err)
	}
	f, err := application.ExcelReport(ds, xp)
	if err != nil {
		t.Fatalf("application.ExcelReport() err = %s", err)
	}

	rows := f.GetRows(f.GetSheetName(f.GetActiveSheetIndex())) // rows is [][]string
	got := len(rows)
	if got != want {
		t.Errorf("GetRows() row count = %d, want %d", got, want)
	}
}
