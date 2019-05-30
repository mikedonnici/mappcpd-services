// Package note provides management of notes data
package note

import (
	"database/sql"
	"errors"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// Error messages
const (
	ErrorIDNotNil          = "cannot insert a note row because Note.ID already has a value"
	ErrorNoMemberID        = "cannot insert a note row because Note.MemberID field is not set"
	ErrorNoTypeID          = "cannot insert a note row because Note.TypeID field is not set"
	ErrorNoContent         = "cannot insert a note row because Note.Content is empty"
	ErrorAssociation       = "association entity not specified"
	ErrorAssociationID     = "association entity ID not specified"
	ErrorAssociationEntity = "association entity invalid"
)

// Note represents a record of a comment, document or anything else. A Note is always linked to a member
// and can also be associated with an application or an issue
type Note struct {
	ID            int    `json:"id" bson:"id"`
	MemberID      int    `json:"memberId" bson:"memberId"`
	Type          string `json:"type" bson:"type"`
	TypeID        int
	Association   string
	AssociationID int
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

// InsertRow creates a new note row with fields from Note
func (n *Note) InsertRow(ds datastore.Datastore) error {
	switch {
	case n.ID > 0:
		return errors.New(ErrorIDNotNil)
	case n.MemberID == 0:
		return errors.New(ErrorNoMemberID)
	case n.TypeID == 0:
		return errors.New(ErrorNoTypeID)
	case n.Content == "":
		return errors.New(ErrorNoContent)
	}
	res, err := ds.MySQL.Session.Exec(queries["insert-note"], n.TypeID, n.Content)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	n.ID = int(id) // from int64

	// Notes differ from issues in that they always require an associated
	// record in wf_note_association. That is, they must always be associated with
	// at least a member id. The member id has already been checked (above).
	err = n.checkAssociationData()
	if err != nil {
		return err
	}
	_, err = ds.MySQL.Session.Exec(queries["insert-note-association"],
		n.ID,
		n.MemberID,
		NullInt(n.AssociationID),
		NullString(n.Association),
	)
	if err != nil {
		return err
	}

	return nil
}

// checkAssociatioData verifies fields required to associate an issue with other data. 
// To associate a Note record with another entity requires the entity name as a string, 
// ie 'application' or 'issue', as well as the id of the record from that entity table.
func (n *Note) checkAssociationData() error {
	// no association
	if n.Association == "" && n.AssociationID == 0 {
		return nil
	}
	if n.Association != "" || n.AssociationID > 0 {
		switch {
		case n.Association == "":
			return errors.New(ErrorAssociation)
		case n.AssociationID == 0:
			return errors.New(ErrorAssociationID)
		case n.Association != "application" && n.Association != "issue":
			return errors.New(ErrorAssociationEntity)
		}
	}
	return nil
}

// NullString allows an empty string value (nil) to be set to NULL in the database
func NullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// NullInt allows an empty int value (nil) to be set to NULL in the database
func NullInt(i int) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: int64(i),
		Valid: true,
	}
}

// ByID fetches a Note from the specified datastore
func ByID(ds datastore.Datastore, id int) (Note, error) {
	n := Note{ID: id}

	// handle potential null values
	association := sql.NullString{}
	associationID := sql.NullInt64{}

	err := ds.MySQL.Session.QueryRow(queries["select-note-by-id"], id).Scan(
		&n.ID,
		&n.Type,
		&n.TypeID,
		&association,
		&associationID,
		&n.MemberID,
		&n.DateCreated,
		&n.DateUpdated,
		&n.DateEffective,
		&n.Content,
	)
	if err != nil {
		return n, err
	}
	// set potentially NULL field values
	if association.String != "" && association.Valid {
		n.Association = association.String
	}
	if associationID.Int64 != 0 && associationID.Valid {
		n.AssociationID = int(associationID.Int64)
	}

	n.Attachments, err = attachments(ds, n.ID)

	return n, err
}

// ByMemberID fetches all the notes linked to a Member from the specified datastore
func ByMemberID(ds datastore.Datastore, memberID int) ([]Note, error) {
	var xn []Note
	q := queries["select-notes-by-member-id"] + " ORDER BY wn.effective_on DESC"
	rows, err := ds.MySQL.Session.Query(q, memberID)
	if err != nil {
		return xn, err
	}
	defer rows.Close()

	for rows.Next() {

		n := Note{}

		// handle potential null values
		association := sql.NullString{}
		associationID := sql.NullInt64{}

		err := rows.Scan(
			&n.ID,
			&n.Type,
			&n.TypeID,
			&association,
			&associationID,
			&n.MemberID,
			&n.DateCreated,
			&n.DateUpdated,
			&n.DateEffective,
			&n.Content,
		)
		if err != nil {
			return xn, err
		}

		// set potentially NULL field values
		if association.String != "" && association.Valid {
			n.Association = association.String
		}
		if associationID.Int64 != 0 && associationID.Valid {
			n.AssociationID = int(associationID.Int64)
		}

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

	query := queries["select-attachments"] + " WHERE wa.wf_note_id = ?"
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
