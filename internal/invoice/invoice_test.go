package invoice_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/invoice"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
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

// test fetch an invoice by id, verify amount
func testByID(t *testing.T) {
	cases := []struct {
		arg  int     // invoice id
		want float64 // amount
	}{
		{1, 110.11},
		{2, 220.22},
	}

	for _, c := range cases {
		i, err := invoice.ByID(ds, c.arg)
		if err != nil {
			t.Errorf("invoice.ByID(%d) err = %s", c.arg, err)
		}
		got := i.Amount
		if got != c.want {
			t.Errorf("invoice.ByID(%d) Amount = %v, want %v", c.arg, got, c.want)
		}
	}
}

// fetch multiple invoices by IDs, check count and sum of amounts
func testByIDs(t *testing.T) {
	cases := []struct {
		arg       []int   // payment IDs
		wantCount int     // number of results
		wantSum   float64 // sum of invoice amounts
	}{
		{[]int{1}, 1, 110.11},
		{[]int{1, 2}, 2, 330.33},
		{[]int{1, 2, 3, 4}, 2, 330.33}, // 3 and 4 don't exist
	}

	for _, c := range cases {
		xi, err := invoice.ByIDs(ds, c.arg)
		if err != nil {
			t.Errorf("invoice.ByIDs(%v) err = %s", c.arg, err)
		}
		gotCount := len(xi)
		if gotCount != c.wantCount {
			t.Errorf("invoice.ByIDs(%v) count = %d, want %d", c.arg, gotCount, c.wantCount)
		}
		var gotSum float64
		for _, i := range xi {
			gotSum += i.Amount
		}
		if gotSum != c.wantSum {
			t.Errorf("invoice.ByIDs(%v) sum amounts = %v, want %v", c.arg, gotSum, c.wantSum)
		}
	}
}

// fetch some test data and ensure excel report is not returning an error
func testExcelReport(t *testing.T) {

	ids := []int{1, 2} // invoice records
	want := 4          // expect 4 rows - heading, 2 records and a total row

	xp, err := invoice.ByIDs(ds, ids)
	if err != nil {
		t.Fatalf("invoice.ByIDs() err = %s", err)
	}
	f, err := invoice.ExcelReport(ds, xp)
	if err != nil {
		t.Fatalf("invoice.ExcelReport() err = %s", err)
	}

	rows := f.GetRows(f.GetSheetName(f.GetActiveSheetIndex())) // rows is [][]string
	got := len(rows)
	if got != want {
		t.Errorf("GetRows() row count = %d, want %d", got, want)
	}
}
