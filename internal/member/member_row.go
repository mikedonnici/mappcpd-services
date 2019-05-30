package member

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cardiacsociety/web-services/internal/issue"
	"github.com/cardiacsociety/web-services/internal/note"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// Foreign Key values for creating required record relationships
const (
	fileNoteTypeID            = 10006
	newApplicationIssueTypeID = 10
)

// Row represents a raw record from the member table in the SQL database. This
// type is primarily for inserting new records. Junction table data are
// represented with []int containing a list of foreign key ids for the relevant
// table. The comment next to JSON tags match the relational db columns names.
type Row struct {
	ID               int    `json:"id"`
	RoleID           int    `json:"roleId"`           // acl_member_role_id
	NamePrefixID     int    `json:"titleId"`          // a_name_prefix_id
	CountryID        int    `json:"countryId"`        // country_id
	ConsentDirectory bool   `json:"consentDirectory"` // consent_directory
	ConsentContact   bool   `json:"consentContact"`   // consent_contact
	UpdatedAt        string `json:"updatedAt"`        // updated_at
	DateOfBirth      string `json:"dateOfBirth"`      // date_of_birth
	Gender           string `json:"gender"`           // gender
	FirstName        string `json:"firstName"`        // first_name
	MiddleNames      string `json:"middleNames"`      // middle_names
	LastName         string `json:"lastName"`         // last_name
	PostNominal      string `json:"postNominal"`      // suffix
	Mobile           string `json:"mobile"`           // mobile_phone
	PrimaryEmail     string `json:"primaryEmail"`     // primary_email

	// The following fields are values represented in junction tables
	Qualifications []QualificationRow `json:"qualifications"`
	Specialities   []SpecialityRow    `json:"interests"`
	Positions      []PositionRow      `json:"positions"`
	Accreditations []AccreditationRow `json:"accreditations"`
	Tags           []TagRow           `json:"tags"`
	Contacts       []ContactRow       `json:"contacts"`

	// Application-related info
	Application ApplicationRow `json:"application"`
}

// QualificationRow represents a member qualification in a junction table. The
// ID of the junction record and the member id are not represented as they are
// not really required.
type QualificationRow struct {
	QualificationID int    `json:"qualificationId"`
	OrganisationID  int    `json:"organisationId"`
	YearObtained    int    `json:"year"`
	Abbreviation    string `json:"abbreviation"`
	Comment         string `json:"comment"`
}

// PositionRow represents a member position in a junction table. The ID of the
// junction record and the member id are not represented as they are not really
// required.
type PositionRow struct {
	PositionID     int    `json:"positionId"`
	OrganisationID int    `json:"organisationId"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
	Comment        string `json:"comment"`
}

// SpecialityRow represents a member speciality in a junction table. The ID of
// the junction record and the member id are not represented as they are not
// really required.
type SpecialityRow struct {
	SpecialityID int    `json:"specialityId"`
	Preference   int    `json:"preference"`
	Comment      string `json:"comment"`
}

// AccreditationRow represents a member accreditation in a junction table. The
// ID of the junction record and the member id are not represented as they are
// not really required.
type AccreditationRow struct {
	AccreditationID int    `json:"accreditationID"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	Comment         string `json:"comment"`
}

// TagRow represents a member tag in a junction table. The ID of the junction
// record and the member id are not represented as they are not really required.
type TagRow struct {
	TagID int `json:"tagID"`
}

// ApplicationRow represents a member's application.
type ApplicationRow struct {
	ID          int
	ForTitleID  int    `json:"forTitleId"`
	NominatorID int    `json:"nominatorId"`
	SeconderID  int    `json:"seconderId"`
	Comment     string `json:"note"`
	FileNote    string `json:"fileNote"`
	FileNoteID  int    `json:"fileNoteId"`
}

// ContactRow represents a contact location. The ID of the junction record and
// the member id are not represented as they are not really required.
type ContactRow struct {
	TypeID    int    `json:"contactTypeId"`
	Phone     string `json:"phone"`
	Fax       string `json:"fax"`
	Email     string `json:"email"`
	Web       string `json:"web"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	Address3  string `json:"address3"`
	Locality  string `json:"locality"`
	State     string `json:"state"`
	Postcode  string `json:"postcode"`
	CountryID int    `json:"countryId"`
}

// StatusRow represents a member's status record
type StatusRow struct {
	ID       int
	MemberID int
	StatusID int
	Current  bool
	Comment  string
}

// Insert inserts a member row into the database. If successful it will set the
// member id.
func (r *Row) Insert(ds datastore.Datastore) error {

	// convert bools to 0/1
	var consentDirectory, consentContact int
	if r.ConsentDirectory {
		consentDirectory = 1
	}
	if r.ConsentContact {
		consentContact = 1
	}

	// gender stored as 'M' or 'F', so capitalise first letter of gender string
	r.Gender = strings.ToUpper(string(strings.TrimSpace(r.Gender)[0]))

	res, err := ds.MySQL.Session.Exec(queries["insert-member-row"],
		r.RoleID,
		r.NamePrefixID,
		r.CountryID,
		consentDirectory,
		consentContact,
		r.DateOfBirth,
		r.Gender,
		r.FirstName,
		r.MiddleNames,
		r.LastName,
		r.PostNominal,
		r.Mobile,
		r.PrimaryEmail)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("LastInsertID() err = %s", err)
	}
	r.ID = int(id) // from int64

	err = r.insertQualifications(ds)
	if err != nil {
		return fmt.Errorf("insertQualifications() err = %s", err)
	}

	err = r.insertPositions(ds)
	if err != nil {
		return fmt.Errorf("insertPositions() err = %s", err)
	}

	err = r.insertSpecialities(ds)
	if err != nil {
		return fmt.Errorf("insertSpecialities() err = %s", err)
	}

	err = r.insertAccreditations(ds)
	if err != nil {
		return fmt.Errorf("insertAccreditations() err = %s", err)
	}

	err = r.insertTags(ds)
	if err != nil {
		return fmt.Errorf("insertTags() err = %s", err)
	}

	err = r.insertContacts(ds)
	if err != nil {
		return fmt.Errorf("insertContacts() err = %s", err)
	}

	err = r.insertApplication(ds)
	if err != nil {
		return fmt.Errorf("insertApplication() err = %s", err)
	}

	err = r.insertFileNote(ds)
	if err != nil {
		return fmt.Errorf("insertFileNote() err = %s", err)
	}

	err = r.insertIssue(ds)
	if err != nil {
		return fmt.Errorf("insertIssue() err = %s", err)
	}

	return nil
}

// insertQualifications inserts the member qualifications present in the Row value
func (r *Row) insertQualifications(ds datastore.Datastore) error {
	for _, q := range r.Qualifications {
		err := q.insert(ds, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertPositions inserts the member positions present in the Row value
func (r *Row) insertPositions(ds datastore.Datastore) error {
	for _, p := range r.Positions {
		err := p.insert(ds, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertSpecialities inserts the member specialities present in the Row value
func (r *Row) insertSpecialities(ds datastore.Datastore) error {
	for _, s := range r.Specialities {
		err := s.insert(ds, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertAccreditations inserts the member accreditations present in the Row value
func (r *Row) insertAccreditations(ds datastore.Datastore) error {
	for _, a := range r.Accreditations {
		err := a.insert(ds, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertTags inserts the member tags present in the Row value
func (r *Row) insertTags(ds datastore.Datastore) error {
	for _, t := range r.Tags {
		err := t.insert(ds, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertApplication creates an application record for the member and sets the
// Application ID on success.
func (r *Row) insertApplication(ds datastore.Datastore) error {
	id, err := r.Application.insert(ds, r.ID)
	if err != nil {
		return err
	}
	r.Application.ID = id
	return nil
}

// insertContacts inserts the member contact rows
func (r *Row) insertContacts(ds datastore.Datastore) error {
	for _, c := range r.Contacts {
		err := c.insert(ds, r.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertFileNote creates a note record associated with this member id and sets
// the Application.FileNoteID on success.
func (r *Row) insertFileNote(ds datastore.Datastore) error {

	// Empty note content will return an error, so ensure it has a value
	if r.Application.FileNote == "" {
		r.Application.FileNote = "Files attached"
	}

	n := note.Note{
		MemberID: r.ID,
		TypeID:   fileNoteTypeID,
		Content:  r.Application.FileNote,
	}

	// noteInsertRow will set note.ID
	err := n.InsertRow(ds)
	if err != nil {
		return err
	}
	r.Application.FileNoteID = n.ID
	return nil
}

// insertIssue raises an issue, of the appropriate type, relating to the new
// application
func (r *Row) insertIssue(ds datastore.Datastore) error {

	// get default deescription and action for the issue type
	issType, err := issue.TypeByID(ds, newApplicationIssueTypeID)
	if err != nil {
		return err
	}

	i := issue.Issue{
		Type:        issue.Type{ID: newApplicationIssueTypeID},
		MemberID:    r.ID,
		Description: issType.Description,
		Action:      issType.Action,
	}
	return i.InsertRow(ds)
}

// insert a member qualification row in the junction table
func (qr QualificationRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-qualification-row"],
		memberID,
		qr.QualificationID,
		qr.OrganisationID,
		qr.YearObtained,
		qr.Abbreviation,
		qr.Comment)
	return err
}

// insert a member position row in the junction table
func (pr PositionRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-position-row"],
		memberID,
		pr.PositionID,
		pr.OrganisationID,
	)
	return err
}

// insert a member speciality row in the junction table
func (sr SpecialityRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-speciality-row"],
		memberID,
		sr.SpecialityID,
		sr.Preference,
		sr.Comment)
	return err
}

// insert a member accreditation row in the junction table
func (ar AccreditationRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-accreditation-row"],
		memberID,
		ar.AccreditationID,
		ar.StartDate,
		ar.EndDate,
		ar.Comment)
	return err
}

// insert a member tag row in the junction table
func (tr TagRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-tag-row"],
		memberID,
		tr.TagID)
	return err
}

// insert methods creates a new application record, returns id on success
func (ar ApplicationRow) insert(ds datastore.Datastore, memberID int) (int, error) {
	res, err := ds.MySQL.Session.Exec(queries["insert-member-application-row"],
		memberID,
		ar.NominatorID,
		ar.SeconderID,
		ar.ForTitleID,
		ar.Comment)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (cr ContactRow) insert(ds datastore.Datastore, memberID int) error {
	_, err := ds.MySQL.Session.Exec(queries["insert-member-contact-row"],
		memberID,
		cr.TypeID,
		cr.CountryID,
		cr.Phone,
		cr.Fax,
		cr.Email,
		cr.Web,
		cr.Address1,
		cr.Address2,
		cr.Address3,
		cr.Locality,
		cr.State,
		cr.Postcode,
	)
	return err
}

// insert a member status row and, if it is set to current, ensure it is the
// only record with current = 1
func (sr StatusRow) insert(ds datastore.Datastore, memberID int) error {

	sr.MemberID = memberID

	// In the db current is stores as an int, so convert to int from bool
	var current int
	if sr.Current {
		current = 1
	}
	res, err := ds.MySQL.Session.Exec(queries["insert-member-status-row"],
		memberID,
		sr.StatusID,
		current,
		sr.Comment,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	sr.ID = int(id)

	// If true also need to set current = 0 for all other status
	// records for the member - can only have one status at a time.
	if sr.Current {
		_, err := ds.MySQL.Session.Exec(queries["update-member-current-status"], sr.ID, memberID)
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertRowFromJSON creates a new member (applicant) Row from a JSON object as well as various
// related rows required for the application process.
func InsertRowFromJSON(ds datastore.Datastore, s string) (Row, error) {
	r := Row{}

	err := json.Unmarshal([]byte(s), &r)
	if err != nil {
		return r, err
	}
	err = r.Insert(ds)
	if err != nil {
		return r, err
	}
	return r, nil
}
