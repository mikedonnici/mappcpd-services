package note_test

import (
	"log"
	"testing"

	"github.com/mappcpd/web-services/internal/note"
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

func TestNoteContent(t *testing.T) {

	cases := []struct {
		ID     int
		Expect string
	}{
		{1, "Application note"},
		{2, "Issue raised"},
	}

	for _, c := range cases {
		r, err := note.ByID(db.Store, c.ID)
		if err != nil {
			t.Fatalf("Database error: %s", err)
		}
		helper.Result(t, c.Expect, r.Content)
	}
}

func TestNoteType(t *testing.T) {

	cases := []struct {
		ID     int
		Expect string
	}{
		{1, "General"},
		{2, "System"},
	}

	for _, c := range cases {
		r, err := note.ByID(db.Store, c.ID)
		if err != nil {
			t.Fatalf("Database error: %s", err)
		}
		helper.Result(t, c.Expect, r.Type)
	}
}

func TestMemberNote(t *testing.T) {

	xn, err := note.ByMemberID(db.Store, 1)
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 3, len(xn))
}

func TestNoteFirstAttachmentUrl(t *testing.T) {

	cases := []struct {
		ID     int
		Expect string
	}{
		{1, "https://cdn.test.com/note/1/1-filename.ext"},
	}

	for _, c := range cases {
		n, err := note.ByID(db.Store, c.ID)
		if err != nil {
			t.Fatalf("Database error: %s", err)
		}
		helper.Result(t, c.Expect, n.Attachments[0].URL)
	}
}
