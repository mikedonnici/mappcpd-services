/*
	Package notes_test provides integration test for notes
*/
package notes_test

import (
	"log"
	"testing"

	"github.com/mappcpd/web-services/internal/notes"
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

func TestNoteByID(t *testing.T) {
	res, err := notes.NoteByIDStore(2, db.MySQL)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, "Issue raised.", res.Content)
}
