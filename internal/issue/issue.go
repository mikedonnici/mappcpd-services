// Package issue provides management of issue data
package issue

import (
	"errors"
	"fmt"

	"github.com/cardiacsociety/web-services/internal/note"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// Error messages
const (
	ErrorIDNotNil              = "cannot insert an issue row because ID already has a value"
	ErrorNoTypeID              = "cannot insert an issue row because Type.ID field is not set"
	ErrorNoDescription         = "cannot insert an issue row because Description is empty"
	ErrorAssociationNoMemberID = "cannot associate issue with another entity unless a member id is specified"
	ErrorAssociation           = "association entity not specified"
	ErrorAssociationID         = "association entity ID not specified"
	ErrorAssociationEntity     = "association entity invalid"
)

// Issue represents a workflow issue
type Issue struct {
	ID          int
	Resolved    bool
	Visible     bool
	Description string
	Action      string

	// The following fields represent data associated with an issue. In the relational database this
	// is how Issues are linked to members and invoices. The association is optionsal and open
	// so as to allow Issue to be raised without any association (gloabl issues) or specifically
	//related to a member or invoice record. As such, any connections must be determined programatically.
	// At this stage issues can only be associated with "application" or "invoice" records
	MemberID      int    // if set, this issue will be associated with this member
	Association   string // either "application" or "invoice"
	AssociationID int    // the id of the associated application or invoice record

	Type  Type
	Notes []note.Note
}

// Type represents the sub-category of the issue, ie Category -> Type
type Type struct {
	ID          int
	Name        string
	Description string
	Action      string
	Category    Category
}

// Category represents the top-level categorisation of issues
type Category struct {
	ID          int
	Name        string
	Description string
}

// InsertRow creates a new issue row with fields from Issue
func (i *Issue) InsertRow(ds datastore.Datastore) error {
	switch {
	case i.ID > 0:
		return errors.New(ErrorIDNotNil)
	case i.Type.ID == 0:
		return errors.New(ErrorNoTypeID)
	case i.Description == "":
		return errors.New(ErrorNoDescription)
	}
	q := fmt.Sprintf(queries["insert-issue"], i.Type.ID, i.Description, i.Action)
	res, err := ds.MySQL.Session.Exec(q)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	i.ID = int(id) // from int64

	// Associate other data if fields are set
	if i.MemberID > 0 || i.AssociationID > 0 || i.Association != "" {
		err := i.checkAssociatioData()
		if err != nil {
			return err
		}
		q := fmt.Sprintf(queries["insert-issue-association"], i.ID, i.MemberID, i.AssociationID, i.Association)
		_, err = ds.MySQL.Session.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

// checkAssociatioData verifies fields required to associate an issue with other data
func (i *Issue) checkAssociatioData() error {
	// no association
	if i.MemberID == 0 && i.Association == "" && i.AssociationID == 0 {
		return nil
	}
	// associate with member only
	if i.MemberID > 0 && i.Association == "" && i.AssociationID == 0 {
		return nil
	}
	// associate with entity
	if i.Association != "" || i.AssociationID > 0 {
		switch {
		case i.MemberID == 0:
			return errors.New(ErrorAssociationNoMemberID)
		case i.Association == "":
			return errors.New(ErrorAssociation)
		case i.AssociationID == 0:
			return errors.New(ErrorAssociationID)
		case i.Association != "application" && i.Association != "invoice":
			return errors.New(ErrorAssociationEntity)
		}
	}
	return nil
}

// ByID fetches an issue by id
func ByID(ds datastore.Datastore, id int) (Issue, error) {
	i := Issue{}
	var resolved, visible int // for converting 0/1 to bool
	q := queries["select-issue-by-id"]
	err := ds.MySQL.Session.QueryRow(q, id).Scan(
		&i.ID,
		&resolved,
		&visible,
		&i.Description,
		&i.Action,
		&i.MemberID,
		&i.Association,
		&i.AssociationID,
		&i.Type.ID,
		&i.Type.Name,
		&i.Type.Description,
		&i.Type.Category.ID,
		&i.Type.Category.Name,
		&i.Type.Category.Description,
	)
	if resolved == 1 {
		i.Resolved = true
	}
	if visible == 1 {
		i.Visible = true
	}
	return i, err
}

// TypeByID fetches an issue type by id
func TypeByID(ds datastore.Datastore, id int) (Type, error) {
	t := Type{}
	q := queries["select-issue-type-by-id"]
	err := ds.MySQL.Session.QueryRow(q, id).Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.Action,
		&t.Category.ID,
		&t.Category.Name,
		&t.Category.Description,
	)
	return t, err
}
