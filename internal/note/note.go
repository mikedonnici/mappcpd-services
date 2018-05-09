// Packages notes provides access to Notes data
package note

import (
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

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

// ByID fetches a Note from the specified datastore - used for testing
func ByID(ds datastore.Datastore, id int) (Note, error) {
	return noteByID(ds, id)
}

// ByMemberID fetches all the notes linked to a Member from the specified datastore - used for testing
func ByMemberID(ds datastore.Datastore, memberID int) ([]Note, error) {
	return notesByMemberID(ds, memberID)
}

// noteByID fetches a Note record from the specified data store
func noteByID(ds datastore.Datastore, id int) (Note, error) {

	n := Note{ID: id}

	query := Queries["select-note"] + " WHERE wn.id = ?"
	err := ds.MySQL.Session.QueryRow(query, id).Scan(
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

	n.Attachments, err = attachments(ds, n.ID)

	return n, err
}

func notesByMemberID(ds datastore.Datastore, memberID int) ([]Note, error) {

	var xn []Note

	query := Queries["select-note"] + " WHERE m.id = ? ORDER BY wn.effective_on DESC"
	rows, err := ds.MySQL.Session.Query(query, memberID)
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
		n.Attachments, err = attachments(ds, n.ID)
		if err != nil {
			return xn, nil
		}

		xn = append(xn, n)
	}

	return xn, nil
}

func attachments(ds datastore.Datastore, noteID int) ([]Attachment, error) {

	var xa []Attachment

	query := Queries["select-attachments"] + " WHERE wa.wf_note_id = ?"
	rows, err := ds.MySQL.Session.Query(query, noteID)
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
