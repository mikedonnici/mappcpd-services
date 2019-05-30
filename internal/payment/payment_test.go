package payment_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/payment"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

func TestAll(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("payment", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testByID", testByID)
		t.Run("testByIDs", testByIDs)
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

// fetch payment by id, verify amount
func testByID(t *testing.T) {
	cases := []struct {
		arg  int     // payment id
		want float64 // payment amount
	}{
		{1, 108.95},
		{2, 10.10},
		{3, 20.20},
		{4, 20.31},
		{5, 30.32},
	}

	for _, c := range cases {
		p, err := payment.ByID(ds, c.arg)
		if err != nil {
			t.Errorf("payment.ByID(%d) err = %s", c.arg, err)
		}
		got := p.Amount
		if got != c.want {
			t.Errorf("Payment.Amount = %v, want %v", got, c.want)
		}
		xb, _ := json.MarshalIndent(p, "", "  ")
		t.Log(string(xb))
	}

}

// fetch multiple Payment values by IDs, check count and sum of all payments
func testByIDs(t *testing.T) {
	cases := []struct {
		arg       []int   // payment IDs
		wantCount int     // number of results
		wantSum   float64 // sum of payment totals
	}{
		{[]int{1}, 1, 108.95},
		{[]int{1, 101}, 1, 108.95},
		{[]int{2, 3, 5}, 3, 60.62},
	}

	for _, c := range cases {
		xp, err := payment.ByIDs(ds, c.arg)
		if err != nil {
			t.Errorf("payment.ByIDs(%v) err = %s", c.arg, err)
		}
		gotCount := len(xp)
		if gotCount != c.wantCount {
			t.Errorf("payment.ByIDs(%v) count = %d, want %d", c.arg, gotCount, c.wantCount)
		}
		var gotSum float64
		for _, p := range xp {
			gotSum += p.Amount
		}
		if gotSum != c.wantSum {
			t.Errorf("payment.ByIDs(%v) sum amounts = %v, want %v", c.arg, gotSum, c.wantSum)
		}
	}
}
