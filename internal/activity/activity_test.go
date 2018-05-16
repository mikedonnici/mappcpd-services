package activity_test

import (
	"log"
	"testing"

	"github.com/mappcpd/web-services/internal/activity"
	"github.com/mappcpd/web-services/testdata"
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

func TestActivityCount(t *testing.T) {
	xa, err := activity.All(db.Store)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 5, len(xa))
}

func TestActivityTypesCount(t *testing.T) {

	cases := []struct {
		id    int
		count int
	}{
		{1, 0},
		{3, 0},
		{20, 9},
		{24, 5},
	}

	for _, c := range cases {
		xa, err := activity.Types(db.Store, c.id)
		if err != nil {
			t.Fatalf("Database error: %s", err)
		}
		helper.Result(t, c.count, len(xa))
	}
}

func TestActivityByID(t *testing.T) {
	cases := []struct {
		id   int
		name string
	}{
		{4, "Presentation"},
		{23, "Group Learning"},
	}

	for _, c := range cases {
		a, err := activity.ByID(db.Store, c.id)
		if err != nil {
			t.Fatalf("Database error: %s", err)
		}
		helper.Result(t, c.name, a.Name)
	}
}

func TestActivityByTypeID(t *testing.T) {
	cases := []struct {
		typeID     int
		activityID int
	}{
		{2, 20},
		{13, 21},
		{28, 23},
		{36, 24},
	}

	for _, c := range cases {
		a, err := activity.ByTypeID(db.Store, c.typeID)
		if err != nil {
			t.Fatalf("Database error: %s", err)
		}
		helper.Result(t, c.activityID, a.ID)
	}
}
