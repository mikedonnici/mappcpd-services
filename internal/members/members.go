package members

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"database/sql"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"

	"github.com/mappcpd/web-services/internal/notes"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/utility"
)

// Note trying to scan NULL db values into strings throws an error. This is discussed here:
// https://github.com/go-sql-driver/mysql/issues/34
// Using []byte is a workaround but then need to convert back to strings. So I've used
// COALESCE() in any SQL where a NULL value is possible... it is a problem with the db
// so might as well make the db deal with it :)

// this file contains the Member "model" -  a struct that maps to the JSON representation
// of the member record represented as a document, and can be unpacked to be mapped to the
// relational model ofr a member

// Member defines struct for member record
type Member struct {
	_id       string    `json:"_id" bson:"_id"`
	ID        int       `json:"id" bson:"id"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`

	// Active refers to the members status in relation to the organisation, ie ms_m_status.ms_status_id = 1 (MySQL)
	// In this model this really belongs in the memberships, however is here from simplicity.
	Active         bool            `json:"active" bson:"active""`
	Title          string          `json:"title" bson:"title"`
	FirstName      string          `json:"firstName" bson:"firstName"`
	MiddleNames    string          `json:"middleNames" bson:"middleNames"`
	LastName       string          `json:"lastName" bson:"lastName"`
	PostNominal    string          `json:"postNominal" bson:"postNominal"`
	Gender         string          `json:"gender" bson:"gender"`
	DateOfBirth    string          `json:"dateOfBirth" bson:"dateOfBirth"`
	Memberships    []Membership    `json:"memberships" bson:"memberships"`
	Contact        MemberContact   `json:"contact" bson:"contact"`
	Qualifications []Qualification `json:"qualifications" bson:"qualifications"`
	Positions      []Position      `json:"positions" bson:"positions"`
	// omitempty to exclude this from sync
	RecurringActivities []RecurringActivity `json:"recurringActivities,omitempty" bson:"recurringActivities,omitempty"`
}

type Members []Member

// Contact struct holds all contact information for a member
type MemberContact struct {
	EmailPrimary   string           `json:"emailPrimary" bson:"emailPrimary"`
	EmailSecondary string           `json:"emailSecondary" bson:"emailSecondary"`
	Mobile         string           `json:"mobile" bson:"mobile"`
	Locations      []MemberLocation `json:"locations" bson:"locations"`

	// Flags that indicate members consent to appear in the directory, and to have contact details shared in directory
	Directory bool `json:"directory" bson:"directory"`
	Consent   bool `json:"consent" bson:"consent"`
}

// Location defines a contact place or contact 'card'
type MemberLocation struct {
	Preference  int    `json:"order" bson:"order"`
	Description string `json:"type" bson:"type"`
	Address     string `json:"address" bson:"address"`
	City        string `json:"city" bson:"city"`
	State       string `json:"state" bson:"state"`
	Postcode    string `json:"postcode" json:"postcode"`
	Country     string `json:"country" bson:"country"`
	Phone       string `json:"phone" bson:"phone"`
	Fax         string `json:"fax" bson:"fax"`
	Email       string `json:"email" bson:"email"`
	URL         string `json:"url" bson:"url"`
}

// Membership holds all of the details for membership to an organisation
type Membership struct {
	OrgID        string            `json:"orgId" bson:"orgId"`
	OrgCode      string            `json:"orgCode" bson:"orgCode"`
	OrgName      string            `json:"orgName" bson:"orgName"`
	Title        string            `json:"title" bson:"title"`
	TitleHistory []MembershipTitle `json:"titleHistory" bson:"titleHistory"`
}

// MembershipTitle refers to the standing, rank or type of membership within an organisation
type MembershipTitle struct {
	Date        string `json:"date" bson:"date"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Comment     string `json:"comment" bson:"comment"`
}

// Qualification is a formal qualification such as a degree, Masters, PHD etc
type Qualification struct {
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Year        string `json:"year" bson:"year"`
}

// Position is an appointment to a board, council or similar
type Position struct {
	OrgCode     string `json:"orgCode" bson:"orgCode"`
	OrgName     string `json:"orgName" bson:"orgName"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Start       string `json:"start" bson:"start"`
	End         string `json:"end" bson:"end"`
}

// SetActive sets the Active boolean value
func (m *Member) SetActive() error {

	// Assume inactive unless otherwise
	m.Active = false

	// store result from db query
	var active int

	query := `SELECT ms_status_id from ms_m_status WHERE
	active = 1 AND current = 1 AND member_id = ?`
	err := datastore.MySQL.Session.QueryRow(query, m.ID).Scan(&active)
	switch {
	case err == sql.ErrNoRows:
		// for the no rows case spit out a message but we don't need to bomb out with a 500
		// otherwise no record will be returned to the caller
		msg := fmt.Sprintf(".SetActive() could not find a record with 'current'/'active' = 1 for member id %v - will assume members 'active' status is false - ", m.ID)
		log.Println(msg, err)
		return nil

	case err != nil:
		msg := ".SetActive() failed"
		log.Println(msg, err)
		return errors.Wrap(err, msg)
	}

	// Found a record, if it has a value of 1 then member is active, otherwise it remains as false
	if active == 1 {
		m.Active = true
	}

	return nil
}

// SetTitle sets the title (Mr, Prof, Dr) and Post nominal, if any
func (m *Member) SetTitle() error {

	query := `SELECT
	COALESCE(a.name, '') FROM a_name_prefix a
	RIGHT JOIN member m ON m.a_name_prefix_id = a.id
	WHERE m.id = ?`
	err := datastore.MySQL.Session.QueryRow(query, m.ID).Scan(&m.Title)
	switch {
	case err == sql.ErrNoRows:
		// Do nothing... there is just no title
		return nil
	case err != nil:
		msg := ".SetTitle() failed"
		log.Println(msg, err)
		return errors.Wrap(err, msg)
	}

	return nil
}

// SetContactLocations populates the Contact.Locations []Location
func (m *Member) SetContactLocations() error {

	query := `SELECT
	     COALESCE(mpct.name, ''),
             CONCAT(
             	COALESCE(mpmc.address1, ''), '\n',
             	COALESCE(mpmc.address2, ''), '\n',
             	COALESCE(mpmc.address3, '')
             	),
             COALESCE(mpmc.locality, ''),
             COALESCE(mpmc.state, ''),
             COALESCE(mpmc.postcode, ''),
             COALESCE(country.name, ''),
             COALESCE( mpmc.phone, ''),
             COALESCE(mpmc.fax, ''),
             COALESCE(mpmc.email, ''),
             COALESCE(mpmc.web, ''),
             COALESCE(mpct.order, '')
             FROM mp_m_contact mpmc
             LEFT JOIN mp_contact_type mpct ON mpmc.mp_contact_type_id = mpct.id
             LEFT JOIN country ON mpmc.country_id = country.id
             WHERE mpmc.member_id = ?
             GROUP BY mpmc.id
             ORDER BY mpct.order ASC`

	//log.Println(sql)

	rows, err := datastore.MySQL.Session.Query(query, m.ID)
	switch {
	case err == sql.ErrNoRows:
		// No rows
		return nil
	case err != nil:
		msg := ".SetContactLocations() sql error"
		log.Println(msg, err)
		return errors.Wrap(err, msg)
	}
	defer rows.Close()

	for rows.Next() {

		l := MemberLocation{}

		err := rows.Scan(
			&l.Description,
			&l.Address,
			&l.City,
			&l.State,
			&l.Postcode,
			&l.Country,
			&l.Phone,
			&l.Fax,
			&l.Email,
			&l.URL,
			&l.Preference,
		)
		if err != nil {
			msg := ".SetContactLocations() failed to scan row"
			log.Println(msg, err)
			return errors.Wrap(err, msg)
		}

		// Trim additional address newlines
		l.Address = strings.Trim(l.Address, "\n")

		m.Contact.Locations = append(m.Contact.Locations, l)
	}

	return nil
}

// GetMemberships populates the Memberships field with one or more
// Membership values
func (m *Member) SetMemberships() error {

	// TODO: SQL to fetch memberships - requires db changes

	// Force selection of more that one membership now for testing
	// Hard coded to CSANZ for now

	// TODO: Add a field called CustomData for any JSON specific to the Membership
	csanz := Membership{
		OrgID:   "csanz",
		OrgCode: "CSANZ",
		OrgName: "Cardiac Society of Australia and New Zealand",
	}

	m.Memberships = append(m.Memberships, csanz)

	return nil
}

// GetTitle populates the MembershipTitle field for a particular Membership.
// It receives the Membership index (mi) which points to the relevant item in []Membership
func (m *Member) SetMembershipTitle(mi int) error {

	// For now we will just set the Member.MembershipTitle field to a string
	// with the name of the title. TitleHistory contains all the details
	// Including the current title so storing them at Member.MembershipTitle is
	// somewhat redundant, and leaving the current title out of the
	// History seems silly as well, as it is part of the history.
	//t := MembershipTitle{}
	t := ""

	// TODO: Make this work for different organisations
	//sql := `SELECT
	//	COALESCE(mmt.granted_on, ''),
	//	"no code",
	//	COALESCE(mt.name, ''),
	//	COALESCE(mt.description, '')
	//	FROM ms_title mt
	//	INNER JOIN ms_m_title mmt ON mt.id = mmt.ms_title_id
	//	WHERE mmt.member_id = ?
	//	AND current = 1
	//	ORDER BY mmt.id DESC
	//	LIMIT 1`

	query := `SELECT
		COALESCE(mt.name, '')
		FROM ms_title mt
		INNER JOIN ms_m_title mmt ON mt.id = mmt.ms_title_id
		WHERE mmt.member_id = ?
		AND current = 1
		ORDER BY mmt.id DESC
		LIMIT 1`

	err := datastore.MySQL.Session.QueryRow(query, m.ID).Scan(
		//&t.Date,
		//&t.Code,
		&t,
		//&t.Description,
	)

	// IMPORTANT - if no rows are found this is NOT an error here
	// it just means there are no memberships. There is a big HACK
	// here whereby the membership is deleted if there is no membership title
	switch {
	case err == sql.ErrNoRows:
		// remove the default membership as there is no title
		m.Memberships = []Membership{}
		return nil
	case err != nil:
		msg := ".SetMembershipTitle() sql error"
		log.Println(msg, err)
		return errors.Wrap(err, msg)
	default:
		// Set the MembershipTitle value for the Membership at index 'mi'
		m.Memberships[mi].Title = t
		return nil
	}
}

// GetTitleHistory populates the Member.TitleHistory field for the Membership
// at index 'mi. Very similar to GetTitle except there may be more than one
// MembershipTitle so it uses []MembershipTitle
func (m *Member) SetMembershipTitleHistory(mi int) error {

	query := `SELECT
		COALESCE(mmt.granted_on, ''),
		"no code",
		COALESCE(mt.name, ''),
		COALESCE(mt.description, ''),
		COALESCE(mmt.comment, '')
		FROM ms_title mt
		INNER JOIN ms_m_title mmt ON mt.id = mmt.ms_title_id
		WHERE mmt.member_id = ?
		ORDER BY mmt.id DESC`

	rows, err := datastore.MySQL.Session.Query(query, m.ID)
	switch {
	case err == sql.ErrNoRows:
		// no rows
		return nil
	case err != nil:
		msg := ".SetMembershipTitleHistory() sql error"
		log.Println(msg, err)
		return errors.Wrap(err, msg)
	}
	defer rows.Close()

	for rows.Next() {

		t := MembershipTitle{}
		err := rows.Scan(
			&t.Date,
			&t.Code,
			&t.Name,
			&t.Description,
			&t.Comment,
		)
		if err != nil {
			msg := ".SetMembershipTitleHistory() failed to scan row"
			fmt.Println(msg, err)
			return errors.Wrap(err, msg)
		}

		//notes := m.GetNotes("", "124")

		//t.Notes = notes

		// Append the historical title to the TitleHistory []MembershipTitle
		m.Memberships[mi].TitleHistory = append(m.Memberships[mi].TitleHistory, t)

	}

	return nil
}

func (m *Member) SetQualifications() error {

	query := `SELECT
	COALESCE(mq.short_name, ''),
	COALESCE(mq.name, ''),
	COALESCE(mq.description, ''),
	COALESCE(mmq.year, '')
	FROM mp_m_qualification mmq
	LEFT JOIN mp_qualification mq on mmq.mp_qualification_id = mq.id
	WHERE mmq.member_id = ?
	ORDER BY year DESC`

	rows, err := datastore.MySQL.Session.Query(query, m.ID)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		msg := ".SetQualifications() sql error"
		fmt.Println(msg, err)
		return errors.Wrap(err, msg)
	}
	defer rows.Close()

	for rows.Next() {

		q := Qualification{}

		err := rows.Scan(
			&q.Code,
			&q.Name,
			&q.Description,
			&q.Year,
		)
		if err != nil {
			msg := ".SetQualifications() failed to scan row"
			fmt.Println(msg, err)
			return errors.Wrap(err, msg)
		}

		// Trim additional address newlines
		//l.Address = strings.Trim(l.Address, "\n")

		m.Qualifications = append(m.Qualifications, q)
	}

	return nil
}

// Get Positions fetches the Positions held by a member
func (m *Member) SetPositions() error {

	query := `SELECT
	COALESCE(organisation.short_name, ''),
	COALESCE(organisation.name, ''),
	COALESCE(mp.short_name, ''),
	COALESCE(mp.name, ''),
	COALESCE(mp.description, ''),
	COALESCE(mmp.start_on, ''),
	COALESCE(mmp.end_on, '')
	FROM mp_m_position mmp
	LEFT JOIN mp_position mp ON mmp.mp_position_id = mp.id
	LEFT JOIN organisation ON mmp.organisation_id = organisation.id
	WHERE mmp.member_id = ?`

	rows, err := datastore.MySQL.Session.Query(query, m.ID)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		msg := ".SetPositions() sql error"
		fmt.Println(msg, err)
		return errors.Wrap(err, msg)
	}
	defer rows.Close()

	for rows.Next() {

		p := Position{}

		err := rows.Scan(
			&p.OrgCode,
			&p.OrgName,
			&p.Code,
			&p.Name,
			&p.Description,
			&p.Start,
			&p.End,
		)
		if err != nil {
			msg := ".SetPositions() failed to scan row"
			log.Println(msg, err)
			return errors.Wrap(err, msg)
		}

		// Trim additional address newlines
		//l.Address = strings.Trim(l.Address, "\n")

		m.Positions = append(m.Positions, p)
	}

	return nil
}

// GetNotes fetches notes relating, optionally those that relate to
// a particular entity 'e'. An 'entity' is a value in the db that
// describes the table (entity) to which the note is linked. For example,
// a note relating to a membership title would have the value mp_title
func (m *Member) GetNotes(entityName string, entityID string) []notes.Note {

	query := `SELECT
		wn.effective_on,
		wn.note,
		wna.association,
		wna.association_entity_id
		FROM wf_note wn
		LEFT JOIN wf_note_association wna ON wn.id = wna.wf_note_id
		WHERE wna.member_id = ?
		%s %s
		ORDER BY wn.effective_on DESC`

	// filter by entity name
	s1 := ""
	if len(entityName) > 0 {
		s1 = " AND " + entityName + " clause here"
	}

	// Further filter by a specific entity id
	s2 := ""
	if len(entityID) > 0 {
		s2 = " AND " + entityID + " clause here"
	}

	query = fmt.Sprintf(query, s1, s2)
	fmt.Println(query)

	// Get the notes relating to this title
	n1 := notes.Note{
		ID:            123,
		DateCreated:   "2016-01-01",
		DateUpdated:   "2016-02-02",
		DateEffective: "2016-03-03",
		Content:       "This is the actual note...",
	}

	n2 := notes.Note{
		ID:            123,
		DateCreated:   "2016-04-01",
		DateUpdated:   "2016-05-02",
		DateEffective: "2016-06-03",
		Content:       "This is the second note...",
	}

	return []notes.Note{n2, n1}
}

// MemberByID fetches a member record by id, populates a Member value
// with some of the data and returns a pointer to Member
func MemberByID(id int) (*Member, error) {

	// Set up a new empty Member
	m := Member{ID: id}

	// Coalesce any NULL-able fields
	query := `SELECT
	created_at,
	updated_at,
	COALESCE(first_name, ''),
	COALESCE(middle_names, ''),
	COALESCE(last_name, ''),
	CONCAT(COALESCE(suffix, ''), ' ', COALESCE(qualifications_other, '')),
	COALESCE(gender, ''),
	COALESCE(date_of_birth, ''),
	COALESCE(primary_email, ''),
    COALESCE(secondary_email, ''),
    COALESCE(mobile_phone, ''),
    consent_directory,
    consent_contact
    FROM member WHERE id = ?`

	// TODO - post nominal fields need to be more cleanly handled in the actual MappCPD application

	// Hold these until we fix them up
	var createdAt string
	var updatedAt string

	err := datastore.MySQL.Session.QueryRow(query, id).Scan(
		&createdAt,
		&updatedAt,
		&m.FirstName,
		&m.MiddleNames,
		&m.LastName,
		&m.PostNominal,
		&m.Gender,
		&m.DateOfBirth,
		&m.Contact.EmailPrimary,
		&m.Contact.EmailSecondary,
		&m.Contact.Mobile,
		&m.Contact.Directory,
		&m.Contact.Consent,
	)
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("MemberByID() could not find record with id %v", id)
		log.Println(msg, err)
		return &m, errors.Wrap(err, msg)
	case err != nil:
		msg := "MemberByID() sql error"
		log.Println(msg, err)
		return &m, errors.Wrap(err, msg)
	}

	// Convert MySQL date time strings to time.Time
	m.CreatedAt, _ = utility.DateTime(createdAt)
	m.UpdatedAt, _ = utility.DateTime(updatedAt)

	err = m.SetActive()
	if err != nil {
		msg := "MemberByID() failed to set Active status"
		log.Println(msg, err)
		return &m, errors.Wrap(err, msg)
	}

	err = m.SetTitle()
	if err != nil {
		msg := "MemberByID() failed to set title"
		log.Println(msg, err)
		return &m, errors.Wrap(err, msg)
	}

	err = m.SetContactLocations()
	if err != nil {
		msg := "MemberByID() failed to set contact locations"
		log.Println(msg, err)
		return &m, errors.Wrap(err, msg)
	}

	// Set Memberships
	err = m.SetMemberships()
	if err != nil {
		msg := "MemberByID() failed to set memberships"
		log.Println(msg, err)
		return &m, errors.Wrap(err, msg)
	}

	// TODO: Membership is a botch and causes all non-member to fail to sync
	// MembershipTitle is a child of each Membership, so we
	// we repeat this for each membership
	for i := range m.Memberships {

		// Set Membership.MembershipTitle for Membership value at [i]
		err = m.SetMembershipTitle(i)
		if err != nil {
			msg := "MemberByID() failed to set membership title"
			log.Println(msg, err)
			return &m, errors.Wrap(err, msg)
		}

		// Set Membership.TitleHistory for Membership value at [i]
		err = m.SetMembershipTitleHistory(i)
		if err != nil {
			msg := "MemberByID() failed to set membership title history"
			log.Println(msg, err)
			return &m, errors.Wrap(err, msg)
		}
	}

	// Set Qualifications
	err = m.SetQualifications()
	if err != nil {
		msg := "MemberByID() failed to set qualifications"
		log.Println(msg, err)
		return &m, errors.Wrap(err, msg)
	}

	// Set Positions
	err = m.SetPositions()
	if err != nil {
		msg := "MemberByID() failed to set positions"
		log.Println(msg, err)
		return &m, errors.Wrap(err, msg)
	}

	return &m, nil
}

// UpdateMember will update the MySQL member record
// 'm' is a map of the JSON body POSTed in which contains any parts of the
// member record we want to update.
func UpdateMember(m map[string]interface{}) error {

	do := false // a flag to decide if we can proceed

	// This maps the json field names with the MySQl column names
	dbMap := map[string]string{
		"gender":      "gender",
		"firstName":   "first_name",
		"middleNames": "middle_names",
		"lastName":    "last_name",
		"dateOfBirth": "date_of_birth",
	}

	query := "UPDATE member SET "
	for i, v := range dbMap {

		if _, ok := m[i]; ok {

			if do == true {
				query = query + ", "
			}
			do = true
			query = query + fmt.Sprintf(`%s="%s"`, v, m[i])
		}
	}
	query = query + fmt.Sprintf(` WHERE id = %v LIMIT 1`, m["id"])

	if do == true {
		fmt.Printf("UpdateMember(): %s\n", query)
		_, err := datastore.MySQL.Session.Exec(query)
		if err != nil {
			msg := "UpdateMember() sql error"
			log.Println(msg, err)
			return errors.Wrap(err, msg)
		}
	} else {
		msg := "UpdateMember() failed"
		err := errors.New("No valid fields posted")
		log.Println(msg, err)
		return errors.Wrap(err, msg)
	}

	return nil
}

// UpdateMemberDoc updates the JSON-formatted member record in MongoDB
func UpdateMemberDoc(m *Member, w *sync.WaitGroup) {

	// Make the selector for Upsert
	mid := map[string]int{"id": m.ID}

	// Get pointer to the Members collection
	mc, err := datastore.MongoDB.MembersCol()
	if err != nil {
		msg := "UpdateMemberDoc() could not get pointer to collection"
		log.Println(msg, err)
		return
	}

	// Upsert
	_, err = mc.Upsert(mid, &m)
	if err != nil {
		log.Printf("Error updating document in Members collection: %s\n", err.Error())
	}

	// Tell wait group we're done, if it was passed in
	if w != nil {
		w.Done()
	}

	log.Println("Updated document in Members collection")
}

// SyncMember synchronises the Member record from MySQL -> MongoDB
// Todo - this should not DECIDE on sync based on updated date... should just do one job
func SyncMember(m *Member) {

	// Fetch the current Doc (if there) and compare updatedAt
	m2, err := DocMembersOne(bson.M{"id": m.ID}, bson.M{})
	if err != nil {
		log.Println("Target document error: ", err, "- so do an upsert")
	}

	msg := fmt.Sprintf("Member id %v - MySQL updated at %s, MongoDB updated at %s", m.ID, m.UpdatedAt, m2.UpdatedAt)
	if m.UpdatedAt.Equal(m2.UpdatedAt) {
		msg += " - NO need to sync"
		log.Println(msg)
		return
	}
	msg += " - syncing..."
	log.Println(msg)

	// Update the document in the Members collection
	var w sync.WaitGroup
	w.Add(1)
	go UpdateMemberDoc(m, &w)
	w.Wait()
}

// DocMembersAll searches the Member collection. Receives query(q) and projection(p)
// It returns []interface{} so that only the projected fields are present. The down side of
// this is that the fields are returned in alphabetical order so it is not as readable
// as the Member struct. Option might be to use the Member struct when no projection
// is specified.
func DocMembersAll(q map[string]interface{}, p map[string]interface{}) ([]interface{}, error) {

	members, err := datastore.MongoDB.MembersCol()
	if err != nil {
		return nil, err
	}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt"})

	// Run query and return results
	var r []interface{}
	err = members.Find(q).Select(p).All(&r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// todo ... deprecate this func in favour of one that returns []Member?
func DocMembersLimit(q map[string]interface{}, p map[string]interface{}, l int) ([]interface{}, error) {

	members, err := datastore.MongoDB.MembersCol()
	if err != nil {
		return nil, err
	}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt"})

	// Run query and return results
	var r []interface{}
	err = members.Find(q).Select(p).Limit(l).All(&r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// DocMembersOne returns one member, unmarshaled into the proper struct
func DocMembersOne(q map[string]interface{}, p map[string]interface{}) (Member, error) {

	m := Member{}

	members, err := datastore.MongoDB.MembersCol()
	if err != nil {
		return m, err
	}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt"})

	err = members.Find(q).Select(p).One(&m)
	if err != nil {
		return m, err
	}

	return m, nil
}

func SearchMembersCollection(q map[string]interface{}, p map[string]interface{}, l int) ([]Member, error) {

	members, err := datastore.MongoDB.MembersCol()
	if err != nil {
		return nil, err
	}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt"})

	// Run query and return results
	var xm []Member
	err = members.Find(q).Select(p).Limit(l).All(&xm)
	if err != nil {
		return nil, err
	}

	return xm, nil
}

// SaveDoc method upserts Member doc to MongoDB
func (m *Member) SaveDoc() error {

	// Make selector for Upsert
	mid := map[string]int{"id": m.ID}

	// Get pointer to the Members collection
	mc, err := datastore.MongoDB.MembersCol()
	if err != nil {
		log.Printf("Error getting pointer to Members collection: %s\n", err.Error())
		return err
	}

	// Upsert
	_, err = mc.Upsert(mid, &m)
	if err != nil {
		log.Printf("Error updating document in Members collection: %s\n", err.Error())
	}

	fmt.Println(".SaveDoc() succeeded")
	return nil
}

// UpdateDoc method updates Member doc to MongoDB
func (m *Member) UpdateDoc() error {

	// Selector
	mid := map[string]int{"id": m.ID}

	// Get pointer to the Members collection
	mc, err := datastore.MongoDB.MembersCol()
	if err != nil {
		log.Printf("Error getting pointer to Members collection: %s\n", err.Error())
		return err
	}

	// Update
	err = mc.Update(mid, &m)
	if err != nil {
		log.Printf("Error updating document in Members collection: %s\n", err.Error())
	}

	fmt.Println(".UpdateDoc() succeeded")
	return nil
}
