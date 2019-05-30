package member_test

import (
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/member"
	"github.com/cardiacsociety/web-services/internal/note"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

// Note names: ds2 and setup2() as package-level identifiers must be unique -
// ds and setup() exist in member_test.go
var ds2 datastore.Datastore

func TestMemberRow(t *testing.T) {

	var teardown func()
	ds2, teardown = setup2()
	defer teardown()
	//ds2, _ = setup2() // keep db

	t.Run("member_row", func(t *testing.T) {
		t.Run("testInsertRow", testInsertRow)
		t.Run("testInsertRowJSON", testInsertRowJSON)
	})
}

func setup2() (datastore.Datastore, func()) {
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

// testInsertRow tests the creation of a new member record
func testInsertRow(t *testing.T) {
	m := member.Row{}
	m.RoleID = 2
	m.NamePrefixID = 1
	m.CountryID = 17
	m.ConsentDirectory = true
	m.ConsentContact = true
	m.UpdatedAt = "2019-01-01"
	m.DateOfBirth = "1970-11-03"
	m.Gender = "M"
	m.FirstName = "Mike"
	m.MiddleNames = "Peter"
	m.LastName = "Donnici"
	m.PostNominal = "B.Sc.Agr"
	m.Mobile = "0402 400 191"
	m.PrimaryEmail = "michael@8o8.io"

	m.Qualifications = []member.QualificationRow{
		member.QualificationRow{
			QualificationID: 11,
			OrganisationID:  222,
			YearObtained:    1992,
			Abbreviation:    "B.Sc.Agr.",
			Comment:         "Major in Crop Science",
		},
		member.QualificationRow{
			QualificationID: 22,
			OrganisationID:  223,
			YearObtained:    1996,
			Abbreviation:    "Grad. Cert. Computing",
			Comment:         "Distance education",
		},
	}

	m.Positions = []member.PositionRow{
		member.PositionRow{
			PositionID:     11,
			OrganisationID: 222,
			StartDate:      "2010-01-01",
			EndDate:        "2012-12-31",
			Comment:        "This is a comment",
		},
		member.PositionRow{
			PositionID:     22,
			OrganisationID: 223,
			StartDate:      "2010-01-01",
			EndDate:        "2012-12-31",
			Comment:        "This is a comment",
		},
	}

	m.Specialities = []member.SpecialityRow{
		member.SpecialityRow{
			SpecialityID: 11,
			Preference:   1,
			Comment:      "This is a comment",
		},
	}

	m.Accreditations = []member.AccreditationRow{
		member.AccreditationRow{
			AccreditationID: 11,
			StartDate:       "2010-01-01",
			EndDate:         "2012-12-31",
			Comment:         "This is a comment",
		},
	}

	m.Tags = []member.TagRow{
		member.TagRow{
			TagID: 1,
		},
		member.TagRow{
			TagID: 2,
		},
		member.TagRow{
			TagID: 3,
		},
	}

	m.Contacts = []member.ContactRow{
		member.ContactRow{
			TypeID:    2, // Directory
			CountryID: 14,
			Phone:     "02 444 66 789",
			Fax:       "02 444 66 890",
			Email:     "any@oldemail.com",
			Web:       "https://thesite.com",
			Address1:  "Leve 12",
			Address2:  "123 Some Street",
			Address3:  "Some large building",
			Locality:  "CityTown",
			State:     "NewShire",
			Postcode:  "1234",
		},
		member.ContactRow{
			TypeID:    1, // Mail
			CountryID: 14,
			Address1:  "Level 12",
			Address2:  "123 Some Street",
			Address3:  "Some large building",
			Locality:  "CityTown",
			State:     "NewShire",
			Postcode:  "1234",
		},
	}

	err := m.Insert(ds2)
	if err != nil {
		t.Fatalf("member.Row.Insert() err = %s", err)
	}
	if m.ID == 0 {
		t.Errorf("member.Row.ID = 0, want > 0")
	}

	// verify a few things about the member record
	mem, err := member.ByID(ds2, m.ID)
	if err != nil {
		t.Fatalf("member.ByID(%d) err = %s", m.ID, err)
	}

	// check number of qualifications
	want := 2
	got := len(mem.Qualifications)
	if got != want {
		t.Errorf("Member.Qualifcations count = %d, want %d", got, want)
	}

	// check number of positions
	want = 2
	got = len(mem.Positions)
	if got != want {
		t.Errorf("Member.Positions count = %d, want %d", got, want)
	}

	// check number of specialities
	want = 1
	got = len(mem.Specialities)
	if got != want {
		t.Errorf("Member.Specialities count = %d, want %d", got, want)
	}

	// check number of accreditations
	want = 1
	got = len(mem.Accreditations)
	if got != want {
		t.Errorf("Member.Accreditations count = %d, want %d", got, want)
	}

	// check number of tags
	want = 3
	got = len(mem.Tags)
	if got != want {
		t.Errorf("Member.Tags count = %d, want %d", got, want)
	}

	// check number of contacts
	want = 2
	got = len(mem.Contact.Locations)
	if got != want {
		t.Errorf("Member.Contact.Locations count = %d, want %d", got, want)
	}
}

// testInsertRowJSON tests the creation of a new member record from a JSON doc
func testInsertRowJSON(t *testing.T) {

	// When this test is passing, below is the format for JSON posted to create a new application
	j := `{
		"roleId" : 2,
		"countryId": 14, 
		"gender": "Male",
		"titleId": 5,
		"firstName": "Mike",
		"middleNames": "Peter",
		"lastName": "Donnici",
		"dateOfBirth": "1970-11-03",
		"primaryEmail": "michael@somewhere.com",
		"mobile": "+61402400191",
		"consentDirectory": true,
		"consentContact": true,

		"qualifications": [
			{
				"qualificationId": 2,
				"name": "Bachelor of Medicine, Bachelor of Surgery",
				"abbreviation": "MBBS",
				"year": 2000,
				"organisationId": 237,
				"organisationName": "University of Sydney"
			},
			{
				"qualificationId": 3,
				"name": "Bachelor of Science",
				"abbreviation": "BSc",
				"year": 1998,
				"organisationId": 237,
				"organisationName": "University of Sydney"
			}
		],

		"interests": [
			{
				"specialityId": 36,
				"name": "Physiotherapist"
			},
			{
				"specialityId": 37,
				"name": "Radiographer"
			},
			{
				"specialityId": 38,
				"name": "Rehab Exercise and Prevention"
			}
		],

		"contacts": [
			{
				"contactTypeId": 1,
				"address1": "123 Some Street",
				"address2": "The second liner",
				"address3": "Third floor",
				"locality": "C-Bay",
				"state": "NSW",
				"postcode": "2999",
				"countryId": 14,
				"phone": "02 6122 3456",
				"fax": "02 6134 5555",
				"email": "bas@das.io",
				"web": "https://baz.io"
			},
			{
				"contactTypeId": 2,
				"address1": "123 Some Street",
				"address2": "The second liner",
				"address3": "Third floor",
				"locality": "C-Bay",
				"state": "NSW",
				"postcode": "2999",
				"countryId": 14,
				"phone": "02 6122 3456",
				"fax": "02 6134 5555",
				"email": "bas@das.io",
				"web": "https://baz.io"
			}
		],

		"positions": [
			{
				"positionId": 1,
				"organisationId": 9
			},
			{
				"positionId": 2,
				"organisationId": 7
			},
			{
				"positionId": 3,
				"organisationId": 10
			}
		],

		"tags": [
			{
				"tagId": 4  
			}
		],

		"application": {
			"forTitleId": 2,
			"nominatorId": 399,
			"note": "qualification note: some additional info about my qualifications\r\nnominators note: some additional info about my nominators\r\nishr: true\r\nagreePrivacy: true\r\nagreeConstitution: true\r\nconsentRequestInfo: true",
			"fileNote": "Uploaded by applicant"
		}
	}`

	row, err := member.InsertRowFromJSON(ds2, j)
	if err != nil {
		t.Fatalf("member.RowFromJSON() err = %s", err)
	}

	// verify a few things about the member record
	mem, err := member.ByID(ds2, row.ID)
	if err != nil {
		t.Fatalf("member.ByID(%d) err = %s", row.ID, err)
	}

	// check number of qualifications
	want := 2
	got := len(mem.Qualifications)
	if got != want {
		t.Errorf("Member.Qualifications count = %d, want %d", got, want)
	}

	// check number of positions
	want = 3
	got = len(mem.Positions)
	if got != want {
		t.Errorf("Member.Positions count = %d, want %d", got, want)
	}

	// check number of specialities
	want = 3
	got = len(mem.Specialities)
	if got != want {
		t.Errorf("Member.Specialities count = %d, want %d", got, want)
	}

	// check number of tags
	want = 1
	got = len(mem.Tags)
	if got != want {
		t.Errorf("Member.Tags count = %d, want %d", got, want)
	}

	// check number of contacts
	want = 2
	got = len(mem.Contact.Locations)
	if got != want {
		t.Errorf("Member.Contact.Locations count = %d, want %d", got, want)
	}

	// Check file note was created
	xn, err := note.ByMemberID(ds2, row.ID)
	if err != nil {
		t.Errorf("note.ByMemberID(%d) err = %s", row.ID, err)
	}
	want = 1
	got = len(xn)
	if got != want {
		t.Errorf("note.ByMemberID() count = %d, want %d", got, want)
	}
}
