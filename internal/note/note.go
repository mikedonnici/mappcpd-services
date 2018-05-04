// Packages notes provides access to Notes data
package note

import (
	"github.com/nleof/goyesql"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

var queries = goyesql.MustParseFile("queries.sql")

// Note represents a record of a comment, document or anything else. A Note can be linked to a member and
// other entities such as an application, or an issue
type Note struct {
	ID            int          `json:"id" bson:"id"`
	Type          string       `json:"type" bson:"type"`
	MemberID      int          `json:"memberId" bson:"memberId"`
	DateCreated   string       `json:"dateCreated" bson:"dateCreated"`
	DateUpdated   string       `json:"dateUpdated" bson:"dateUpdated"`
	DateEffective string       `json:"dateEffective" bson:"dateEffective"`
	Content       string       `json:"content" bson:"content"`
	Attachments   []Attachment `json:"attachments" bson:"attachments"`
}

// Attachment is a file linked to a note
type Attachment struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

// ByID fetches a Note from the default datastore
func ByID(id int) (Note, error) {
	return noteByID(id, datastore.MySQL)
}

// ByIDStore fetches a Note from the specified datastore - used for testing
func ByIDStore(id int, conn datastore.MySQLConnection) (Note, error) {
	return noteByID(id, conn)
}

// ByMemberID fetches all the notes linked to a Member
func ByMemberID(memberID int) ([]Note, error) {
	return notesByMemberID(memberID, datastore.MySQL)
}

// ByMemberIDStore fetches all the notes linked to a Member from the specified datastore - used for testing
func ByMemberIDStore(memberID int, conn datastore.MySQLConnection) ([]Note, error) {
	return notesByMemberID(memberID, conn)
}

// noteByID fetches a Note record from the specified data store
func noteByID(id int, conn datastore.MySQLConnection) (Note, error) {

	n := Note{ID: id}

	query := queries["select-note"] + " WHERE wn.id = ?"
	err := conn.Session.QueryRow(query, id).Scan(
		&n.ID,
		&n.Type,
		&n.MemberID,
		&n.DateCreated,
		&n.DateUpdated,
		&n.DateEffective,
		&n.Content,
	)
	if err != nil {
		return n, err
	}

	n.Attachments, err = attachments(n.ID, conn)

	return n, err
}

func notesByMemberID(memberID int, conn datastore.MySQLConnection) ([]Note, error) {

	var xn []Note

	query := queries["select-note"] + " WHERE m.id = ? ORDER BY wn.effective_on DESC"
	rows, err := conn.Session.Query(query, memberID)
	if err != nil {
		return xn, err
	}
	defer rows.Close()

	for rows.Next() {
		n := Note{}
		rows.Scan(
			&n.ID,
			&n.Type,
			&n.MemberID,
			&n.DateCreated,
			&n.DateUpdated,
			&n.DateEffective,
			&n.Content,
		)

		var err error
		n.Attachments, err = attachments(n.ID, conn)
		if err != nil {
			return xn, nil
		}

		xn = append(xn, n)
	}

	return xn, nil
}

func attachments(noteID int, conn datastore.MySQLConnection) ([]Attachment, error) {

	var xa []Attachment

	query := queries["select-attachments"] + " WHERE wa.wf_note_id = ?"
	rows, err := conn.Session.Query(query, noteID)
	if err != nil {
		return xa, err
	}
	defer rows.Close()

	for rows.Next() {
		a := Attachment{}
		err := rows.Scan(&a.ID, &a.Name, &a.URL)
		if err != nil {
			return xa, err
		}
		xa = append(xa, a)
	}

	return xa, nil
}
