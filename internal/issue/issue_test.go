package issue_test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/cardiacsociety/web-services/internal/issue"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

func TestIssue(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("issue", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testIssueByID", testIssueByID)
		t.Run("testIssueTypeByID", testIssueTypeByID)
		t.Run("testInsertRowErrorIDNotNil", testInsertRowErrorIDNotNil)
		t.Run("testInsertRowErrorNoTypeID", testInsertRowErrorNoTypeID)
		t.Run("testInsertRowErrorNoDescription", testInsertRowErrorNoDescription)
		t.Run("testInsertRowErrorAssociationNoMemberID", testInsertRowErrorAssociationNoMemberID)
		t.Run("testInsertRowErrorAssociation", testInsertRowErrorAssociation)
		t.Run("testInsertRowErrorAssociationID", testInsertRowErrorAssociationID)
		t.Run("testInsertRowErrorAssociationEntity", testInsertRowErrorAssociationEntity)
		t.Run("testInsertRow", testInsertRow)
		t.Run("testInsertRowWithAssociation", testInsertRowWithAssociation)
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
		t.Fatalf("Ping() err = %s", err)
	}
}

func testIssueByID(t *testing.T) {
	arg := 1 // issue id in test data
	got, err := issue.ByID(ds, arg)
	if err != nil {
		t.Fatalf("issue.ByID(%d) err = %s", arg, err)
	}

	// This is what we expect in return
	want := issue.Issue{
		ID:       1,
		Resolved: true,
		Visible:  true,
		Type: issue.Type{
			ID:          1,
			Name:        "Invoice Raised",
			Description: "A new invoice has been raised and is pending payment.",
			Category: issue.Category{
				ID:          4,
				Name:        "Finance",
				Description: "Issues relating to Subscriptions, Invoicing and Payments.",
			},
		},
		Description:   "A new invoice has been raised and is pending payment. (INV0001)",
		Action:        "Members can pay online or by alternate methods specified on the invoice.",
		Notes:         nil,
		MemberID:      502,
		Association:   "invoice",
		AssociationID: 1,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("issue.ByID(%d) got !DeepEqual want", arg)
	}
}

// test selecting just the issue type
func testIssueTypeByID(t *testing.T) {
	arg := 1 // issue type id in test data
	got, err := issue.TypeByID(ds, arg)
	if err != nil {
		t.Fatalf("issueTypeByID(%d) err = %s", arg, err)
	}

	// This is what we expect in return
	want := issue.Type{
		ID:          1,
		Name:        "Invoice Raised",
		Description: "A new invoice has been raised and is pending payment.",
		Action:      "Members can pay online or by alternate methods specified on the invoice.",
		Category: issue.Category{
			ID:          4,
			Name:        "Finance",
			Description: "Issues relating to Subscriptions, Invoicing and Payments.",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("issue.TypeByID(%d) got !DeepEqual want", arg)
	}
}

// test an attempt to insert an issue row when the Issue.ID has a value
func testInsertRowErrorIDNotNil(t *testing.T) {
	i := issue.Issue{
		ID:          1,
		Type:        issue.Type{ID: 2},
		Description: "This is the description",
		Action:      "This is what must be done",
	}
	err := i.InsertRow(ds)
	want := issue.ErrorIDNotNil
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

// test an attempt to insert an issue row with Issue.Type.ID not set
func testInsertRowErrorNoTypeID(t *testing.T) {
	i := issue.Issue{
		Description: "This is the description",
		Action:      "This is what must be done",
	}
	err := i.InsertRow(ds)
	want := issue.ErrorNoTypeID
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

// test an attempt to insert an issue row with Issue.Description not set
func testInsertRowErrorNoDescription(t *testing.T) {
	i := issue.Issue{
		Type:   issue.Type{ID: 2},
		Action: "This is what must be done",
	}
	err := i.InsertRow(ds)
	want := issue.ErrorNoDescription
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

func testInsertRowErrorAssociationNoMemberID(t *testing.T) {
	i := issue.Issue{
		Type:          issue.Type{ID: 2},
		Description:   "This is the description",
		Action:        "This is what must be done",
		MemberID:      0, // err
		AssociationID: 345,
		Association:   "application",
	}
	err := i.InsertRow(ds)
	want := issue.ErrorAssociationNoMemberID
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

func testInsertRowErrorAssociation(t *testing.T) {
	i := issue.Issue{
		Type:          issue.Type{ID: 2},
		Description:   "This is the description",
		Action:        "This is what must be done",
		MemberID:      123,
		AssociationID: 345,
		Association:   "", // err
	}
	err := i.InsertRow(ds)
	want := issue.ErrorAssociation
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

func testInsertRowErrorAssociationID(t *testing.T) {
	i := issue.Issue{
		Type:          issue.Type{ID: 2},
		Description:   "This is the description",
		Action:        "This is what must be done",
		MemberID:      123,
		AssociationID: 0, // err
		Association:   "application",
	}
	err := i.InsertRow(ds)
	want := issue.ErrorAssociationID
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

func testInsertRowErrorAssociationEntity(t *testing.T) {
	i := issue.Issue{
		Type:          issue.Type{ID: 2},
		Description:   "This is the description",
		Action:        "This is what must be done",
		MemberID:      123,
		AssociationID: 345,
		Association:   "unknownentity", // err
	}
	err := i.InsertRow(ds)
	want := issue.ErrorAssociationEntity
	if err == nil {
		t.Errorf("Issue.InsertRow() err = nil, want %q", want)
	}
}

// test insert a row without an association
func testInsertRow(t *testing.T) {
	i := issue.Issue{
		Type:        issue.Type{ID: 2},
		Description: "This is the description",
		Action:      "This is what must be done",
	}
	err := i.InsertRow(ds)
	if err != nil {
		t.Errorf("Issue.InsertRow() err = %s", err)
	}
}

// test insert a row with an association
func testInsertRowWithAssociation(t *testing.T) {
	i := issue.Issue{
		Type:          issue.Type{ID: 2},
		Description:   "This is the description",
		Action:        "This is what must be done",
		MemberID:      123,
		AssociationID: 456,
		Association:   "application",
	}
	err := i.InsertRow(ds)
	if err != nil {
		t.Errorf("Issue.InsertRow() err = %s", err)
	}

	// Verify the association
	iss, err := issue.ByID(ds, i.ID)
	if err != nil {
		t.Fatalf("issue.ByID(%d) err = %s", i.ID, err)
	}
	gotAssociationID := iss.AssociationID
	wantAssociationID := i.AssociationID
	if gotAssociationID != wantAssociationID {
		t.Errorf("Issue.AssociationID = %d, want %d", gotAssociationID, wantAssociationID)
	}
	gotAssociation := iss.Association
	wantAssociation := i.Association
	if gotAssociation != wantAssociation {
		t.Errorf("Issue.Association = %q, want %q", gotAssociation, wantAssociation)
	}

	t.Log(toJSON(iss))
}

func toJSON(i interface{}) string {
	xb, _ := json.MarshalIndent(i, "", " ")
	return string(xb)
}
