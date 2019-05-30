package note_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/note"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

func TestNote(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("note", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testNoteContent", testNoteContent)
		t.Run("testNoteType", testNoteType)
		t.Run("testMemberNote", testMemberNote)
		t.Run("testNoteFirstAttachmentUrl", testNoteFirstAttachmentUrl)
		t.Run("testInsertRowErrorIDNotNil", testInsertRowErrorIDNotNil)
		t.Run("testInsertRowErrorNoMemberID", testInsertRowErrorNoMemberID)
		t.Run("testInsertRowErrorNoTypeID", testInsertRowErrorNoTypeID)
		t.Run("testInsertRowErrorNoContent", testInsertRowErrorNoContent)
		t.Run("testInsertRow", testInsertRow)
		t.Run("testInsertRowAssociation", testInsertRowAssociation)
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
		t.Fatal("Could not ping database")
	}
}

func testNoteContent(t *testing.T) {
	cases := []struct {
		arg  int
		want string
	}{
		{1, "Application note"},
		{2, "Issue raised"},
	}
	for _, c := range cases {
		n, err := note.ByID(ds, c.arg)
		if err != nil {
			t.Errorf("note.ByID(%d) err = %s", c.arg, err)
		}
		got := n.Content
		if got != c.want {
			t.Errorf("Note.Content = %q, want %q", got, c.want)
		}
	}
}

func testNoteType(t *testing.T) {
	cases := []struct {
		arg  int
		want string
	}{
		{1, "General"},
		{2, "System"},
	}
	for _, c := range cases {
		n, err := note.ByID(ds, c.arg)
		if err != nil {
			t.Errorf("note.ByID(%d) err = %s", c.arg, err)
		}
		got := n.Type
		if got != c.want {
			t.Errorf("Note.Type = %q, want %q", got, c.want)
		}
	}
}

func testMemberNote(t *testing.T) {
	arg := 1 // member id
	xn, err := note.ByMemberID(ds, 1)
	if err != nil {
		t.Fatalf("note.ByMemberID(%d) err = %s", arg, err)
	}
	got := len(xn)
	want := 3
	if got != want {
		t.Errorf("note.ByMemberID(%d) count = %d, want %d", arg, got, want)
	}
}

func testNoteFirstAttachmentUrl(t *testing.T) {
	cases := []struct {
		arg  int    // note id
		want string // file url
	}{
		{1, "https://cdn.test.com/note/1/1-filename.ext"},
	}

	for _, c := range cases {
		n, err := note.ByID(ds, c.arg)
		if err != nil {
			t.Errorf("note.ByID(%d) err = %s", c.arg, err)
		}
		var got string
		if len(n.Attachments) > 0 {
			got = n.Attachments[0].URL
		}
		if got != c.want {
			t.Errorf("Note.Attachments[0].URL = %s, want %s", got, c.want)
		}
	}
}

// test an attempt to insert an issue row when the Issue.ID has a value
func testInsertRowErrorIDNotNil(t *testing.T) {
	n := note.Note{
		ID:       1,
		TypeID:   10001,
		MemberID: 1,
		Content:  "This is the note content",
	}
	var err error
	err = n.InsertRow(ds)
	got := err.Error()
	want := note.ErrorIDNotNil
	if got != want {
		t.Errorf("Note.InsertRow() err = %q, want %q", got, want)
	}
}

// test an attempt to insert a note row with no MemberID
func testInsertRowErrorNoMemberID(t *testing.T) {
	n := note.Note{
		TypeID:  10001,
		Content: "This is the note content",
	}
	var err error
	err = n.InsertRow(ds)
	got := err.Error()
	want := note.ErrorNoMemberID
	if got != want {
		t.Errorf("Note.InsertRow() err = %q, want %q", got, want)
	}
}

// test an attempt to insert a note row with no TypeID
func testInsertRowErrorNoTypeID(t *testing.T) {
	n := note.Note{
		MemberID: 1,
		Content:  "This is the note content",
	}
	var err error
	err = n.InsertRow(ds)
	got := err.Error()
	want := note.ErrorNoTypeID
	if got != want {
		t.Errorf("Note.InsertRow() err = %q, want %q", got, want)
	}
}

// test an attempt to insert a note row with no content
func testInsertRowErrorNoContent(t *testing.T) {
	n := note.Note{
		MemberID: 1,
		TypeID: 10001,
	}
	var err error
	err = n.InsertRow(ds)
	got := err.Error()
	want := note.ErrorNoContent
	if got != want {
		t.Errorf("Note.InsertRow() err = %q, want %q", got, want)
	}
}

// test insert a row - will always have an association record for memberID
func testInsertRow(t *testing.T) {
	n := note.Note{
		TypeID: 10001,
		MemberID: 1,
		Content:  "This is the note content",
	}
	err := n.InsertRow(ds)
	if err != nil {
		t.Errorf("Note.InsertRow() err = %s", err)
	}
	// Re-fetch the note and verify the member id
	n2, err := note.ByID(ds, n.ID)
	if err != nil {
		t.Fatalf("note.ByID(%d) err = %s", n.ID, err)
	}
	got := n2.MemberID
	want := n.MemberID
	if got != want {
		t.Errorf("note.MemberID = %d, want %d", got, want)
	}
}

// test insert a row with an association with other data
func testInsertRowAssociation(t *testing.T) {
	n := note.Note{
		TypeID: 10001,
		MemberID: 1,
		Content:  "This is the note content",
		Association: "application",
		AssociationID: 1,
	}
	err := n.InsertRow(ds)
	if err != nil {
		t.Errorf("Note.InsertRow() err = %s", err)
	}
	// Re-fetch the note and verify the association
	n2, err := note.ByID(ds, n.ID)
	if err != nil {
		t.Fatalf("note.ByID(%d) err = %s", n.ID, err)
	}
	got := n2.Association
	want := n.Association
	if got != want {
		t.Errorf("note.Association = %q, want %q", got, want)
	}
}
