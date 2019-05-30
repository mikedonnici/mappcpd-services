package member_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/cardiacsociety/web-services/internal/member"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
	"github.com/matryer/is"
	"gopkg.in/mgo.v2/bson"
)

var ds datastore.Datastore

const doTeardown = true

func TestMember(t *testing.T) {

	if doTeardown {
		var teardown func()
		ds, teardown = setup()
		defer teardown()
	} else {
		ds, _ = setup()
	}

	t.Run("member", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testByID", testByID)
		t.Run("testSearchDocDB", testSearchDocDB)
		t.Run("testSaveDocDB", testSaveDocDB)
		t.Run("testSyncUpdated", testSyncUpdated)
		t.Run("testExcelReport", testExcelReport)
		t.Run("testExcelReportJournal", testExcelReportJournal)
		t.Run("testLapse", testLapse)
	})
}

func setup() (datastore.Datastore, func()) {
	var db = testdata.NewDataStore()
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalf("db.SetupMySQL() err = %s", err)
	}
	err = db.SetupMongoDB()
	if err != nil {
		log.Fatalln(err)
	}
	return db.Store, func() {
		err := db.TearDownMySQL()
		if err != nil {
			log.Fatalf("db.TearDownMySQL() err = %s", err)
		}
	}
}

func testPingDatabase(t *testing.T) {
	is := is.New(t)
	err := ds.MySQL.Session.Ping()
	is.NoErr(err) // Could not ping test database
}

func testByID(t *testing.T) {
	is := is.New(t)
	m, err := member.ByID(ds, 1)
	is.NoErr(err)                                              // Error fetching member by id
	is.True(m.Active)                                          // Active should be true
	is.Equal(m.LastName, "Donnici")                            // Last name incorrect
	is.True(len(m.Memberships) > 0)                            // No memberships
	is.Equal(m.Memberships[0].Title, "Associate")              // Incorrect membership title
	is.Equal(m.Contact.EmailPrimary, "michael@mesa.net.au")    // Email incorrect
	is.Equal(m.Contact.Mobile, "0402123123")                   // Mobile incorrect
	is.Equal(m.Contact.Locations[0].City, "Jervis Bay")        // Location city incorrect
	is.Equal(m.Qualifications[0].Name, "PhD")                  // Qualification incorrect
	is.Equal(m.Specialities[1].Name, "Cardiac Cath Lab Nurse") // Speciality incorrect
	//printJSON(*m)
}

func testSearchDocDB(t *testing.T) {
	is := is.New(t)
	q := bson.M{"id": 7821}
	m, err := member.SearchDocDB(ds, q)
	is.NoErr(err)                     // Error querying MongoDB
	is.Equal(m[0].LastName, "Rousos") // Last name incorrect
}

func testSaveDocDB(t *testing.T) {
	is := is.New(t)
	mem := member.Member{
		ID:          1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Active:      true,
		Title:       "Mr",
		FirstName:   "Michael",
		MiddleNames: []string{"Peter"},
		LastName:    "Donnici",
		Gender:      "M",
		DateOfBirth: "1970-11-03",
	}
	err := mem.SaveDocDB(ds)
	is.NoErr(err) // Error saving to MongoDB

	q := bson.M{"lastName": "Donnici"}
	xm, err := member.SearchDocDB(ds, q)
	m := xm[0]
	is.NoErr(err)     // Error querying MongoDB
	is.Equal(m.ID, 1) // ID should be 1
}

func testSyncUpdated(t *testing.T) {
	is := is.New(t)
	mem := member.Member{
		ID:          2,
		CreatedAt:   time.Now().Add(-10 * time.Duration(time.Minute)), // 10 mins ago
		UpdatedAt:   time.Now().Add(-10 * time.Duration(time.Minute)), // 10 mins ago
		Active:      true,
		Title:       "Mr",
		FirstName:   "Barry",
		LastName:    "White",
		Gender:      "M",
		DateOfBirth: "1945-03-15",
	}
	err := mem.SaveDocDB(ds)
	is.NoErr(err) // Error saving to MongoDB

	memUpdate := member.Member{
		ID:          2,
		CreatedAt:   time.Now().Add(-10 * time.Duration(time.Minute)), // 10 mins ago
		UpdatedAt:   time.Now(),                                       // should trigger update
		Active:      false,
		Title:       "Mr",
		FirstName:   "Barry",
		LastName:    "White",
		Gender:      "M",
		DateOfBirth: "1948-03-15",
	}
	err = memUpdate.SyncUpdated(ds)
	is.NoErr(err) // Error syncing to MongoDB

	q := bson.M{"lastName": "White"}
	xm, err := member.SearchDocDB(ds, q)
	m := xm[0]
	is.NoErr(err)                         // Error querying MongoDB
	is.Equal(m.ID, 2)                     // ID should be 2
	is.Equal(m.Active, false)             // Active should be false
	is.Equal(m.DateOfBirth, "1948-03-15") // DateOfBirth incorrect
}

// fetch some test data and ensure excel report is not returning an error
func testExcelReport(t *testing.T) {

	id := 1   // member record
	want := 2 // expect 2 rows - heading and 2 record

	m, err := member.ByID(ds, id)
	if err != nil {
		t.Fatalf("member.ByID() err = %s", err)
	}
	xm := []member.Member{*m}
	f, err := member.ExcelReport(xm)
	if err != nil {
		t.Fatalf("member.ExcelReport() err = %s", err)
	}

	rows := f.GetRows(f.GetSheetName(f.GetActiveSheetIndex())) // rows is [][]string
	got := len(rows)
	if got != want {
		t.Errorf("GetRows() row count = %d, want %d", got, want)
	}
}

// fetch some test data and ensure excel report (journal) is not returning an error
func testExcelReportJournal(t *testing.T) {

	id := 1   // member record
	want := 2 // expect 2 rows - heading and 2 record

	m, err := member.ByID(ds, id)
	if err != nil {
		t.Fatalf("member.ByID() err = %s", err)
	}
	xm := []member.Member{*m}
	f, err := member.ExcelReportJournal(xm)
	if err != nil {
		t.Fatalf("member.ExcelReportJournal() err = %s", err)
	}

	rows := f.GetRows(f.GetSheetName(f.GetActiveSheetIndex())) // rows is [][]string
	got := len(rows)
	if got != want {
		t.Errorf("GetRows() row count = %d, want %d", got, want)
	}
}

// test lapsing a member - @todo: check the actual result
func testLapse(t *testing.T) {
	m, err := member.ByID(ds, 1)
	if err != nil {
		t.Fatalf("ByID() err = %s", err)
	}
	err = m.Lapse(ds)
	if err != nil {
		t.Fatalf("member.Lapse() err = %s", err)
	}
}

func printJSON(m member.Member) {
	xb, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println("-------------------------------------------------------------------")
	fmt.Print(string(xb))
	fmt.Println("-------------------------------------------------------------------")
}
