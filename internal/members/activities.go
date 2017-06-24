package members

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gopkg.in/go-playground/validator.v9"

	"github.com/mappcpd/web-services/internal/activities"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// MemberActivityDoc is the document format for an activity that is
// recorded by a member - that is, a CPD diary entry
type MemberActivityDoc struct {
	ID          int                         `json:"id" bson:"id"`
	MemberID    int                         `json:"memberId" bson:"memberId"`
	CreatedAt   time.Time                   `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time                   `json:"updatedAt" bson:"updatedAt"`
	Date        string                      `json:"date" bson:"date"`
	DateISO     time.Time                   `json:"dateISO" bson:"dateISO"`
	Credit      float32                     `json:"credit" bson:"credit"`
	Description string                      `json:"description" bson:"description"`
	Category    activities.ActivityCategory `json:"category" bson:"category"`
	Activity    activities.Activity         `json:"activity" bson:"activity"`
	CreditData  activities.ActivityCredit   `json:"creditData" bson:"creditData"`
}

// MemberActivityRow represents the minimum data to add or update a Member Activity.
// 'Row' implies a representation of the relevant SQL table row.
type MemberActivityRow struct {
	ID          int     `json:"ID"`
	MemberID    int     `json:"memberID"`
	ActivityID  int     `json:"activityID" validate:"required,min=1"`
	Evidence    int     `json:"evidence"`
	Date        string  `json:"date" validate:"required"`
	Quantity    float32 `json:"quantity" validate:"required"`
	UnitCredit  float32 `json:"unitCredit"`
	Description string  `json:"description" validate:"required"`
}

// MemberActivities is a collection of MemberActivityDoc values
type MemberActivities []MemberActivityDoc

// MemberActivityByID fetches a member activity record by id
func MemberActivityByID(id int) (*MemberActivityDoc, error) {

	// Create Activity value
	a := MemberActivityDoc{ID: id}

	// Coalesce any NULL-able fields
	query := `SELECT
		cma.member_id,
		cma.activity_on,
		COALESCE(cma.description, ''),
		(cma.quantity * cma.points_per_unit),
		cma.quantity,
		'no unit code in model',
		COALESCE(cau.name, ''),
		COALESCE(cau.description, ''),
		cma.points_per_unit,
		cac.id,
		"cac.Code",
		COALESCE(cac.name, ''),
		COALESCE(cac.description, ''),
		ca.id,
		COALESCE(ca.code, ''),
		COALESCE(ca.name, ''),
		COALESCE(ca.description, '')
		FROM ce_m_activity cma
		LEFT JOIN ce_activity ca ON cma.ce_activity_id = ca.id
		LEFT JOIN ce_activity_unit cau ON ca.ce_activity_unit_id = cau.id
		LEFT JOIN ce_activity_category cac ON ca.ce_activity_category_id = cac.id
		WHERE cma.id = ?`

	err := datastore.MySQL.Session.QueryRow(query, id).Scan(
		&a.MemberID,
		&a.Date,
		&a.Description,
		&a.Credit,
		&a.CreditData.Quantity,
		&a.CreditData.UnitCode,
		&a.CreditData.UnitName,
		&a.CreditData.UnitDescription,
		&a.CreditData.UnitCredit,
		&a.Category.ID,
		&a.Category.Code,
		&a.Category.Name,
		&a.Category.Description,
		&a.Activity.ID,
		&a.Activity.Code,
		&a.Activity.Name,
		&a.Activity.Description,
	)
	if err != nil {
		return &a, err
	}

	// Add ISODate for MongoDB from date string
	// Note the first arg is a REFERENCE time - it must always be the same
	// date and time and is defined in the time package: Mon Jan 2 15:04:05 MST 2006
	// So here we just use the bits we need to match MySQL date string: YYYY-MM-DD
	a.DateISO, err = time.Parse("2006-01-02", a.Date)
	if err != nil {
		log.Printf("Error creating ISODate: %s", err.Error())
	}

	return &a, nil
}

// MemberActivityRowByID fetches an ActivityRow value by ID
func MemberActivityRowByID(id int) (*MemberActivityRow, error) {

	// Create ActivityRow value
	a := MemberActivityRow{ID: id}

	// Coalesce any NULL-able fields
	query := `SELECT
		cma.member_id,
		ca.id,
		cma.evidence,
		cma.activity_on,
		cma.quantity,
		COALESCE(cma.description, '')
		FROM ce_m_activity cma
		LEFT JOIN ce_activity ca ON cma.ce_activity_id = ca.id
		WHERE cma.id = ?`

	err := datastore.MySQL.Session.QueryRow(query, id).Scan(
		&a.MemberID,
		&a.ActivityID,
		&a.Evidence,
		&a.Date,
		&a.Quantity,
		&a.Description,
	)
	if err != nil {
		return &a, err
	}

	return &a, nil
}

// MemberActivitiesByMemberID fetches activities for a particular member
func MemberActivitiesByMemberID(id int) ([]MemberActivityDoc, error) {

	activities := MemberActivities{}

	sql := fmt.Sprintf("SELECT id from ce_m_activity WHERE member_id = %v", id)
	rows, err := datastore.MySQL.Session.Query(sql)
	if err != nil {
		return activities, err
	}
	defer rows.Close()

	var activityIDs []int

	for rows.Next() {

		var id int
		rows.Scan(&id)
		activityIDs = append(activityIDs, id)

		a, err := MemberActivityByID(id)
		if err != nil {
			return activities, err
		}
		activities = append(activities, *a)

	}

	return activities, nil
}

// UpdateMemberActivityDoc updates the JSON-formatted activity record in the Doc DB (MongoDB)
func UpdateMemberActivityDoc(a *MemberActivityDoc, w *sync.WaitGroup) {

	// Make the selector for Upsert
	id := map[string]int{"id": a.ID}

	// Get pointer to the collection
	c, err := datastore.MongoDB.ActivitiesCol()
	if err != nil {
		log.Printf("Error getting pointer to Activities collection: %s\n", err.Error())
	}

	// Upsert
	_, err = c.Upsert(id, &a)
	if err != nil {
		log.Printf("Error updating Activity id %s in Activities collection: %s\n", a.ID, err.Error())
	}

	// Tell wait group we're done, if it was passed in
	w.Done()
	log.Printf("Updated Activity id %s Activities collection\n", a.ID)
}

// AddMemberActivity inserts a new member activity in the MySQL db and returns the new id on success.
func AddMemberActivity(a MemberActivityRow) (int64, error) {

	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return 0, err
	}

	// Look up the credit-per-unit for this type of activity...
	uc, err := activities.ActivityUnitCredit(a.ActivityID)
	if err != nil {
		return 0, err
	}
	a.UnitCredit = uc

	query := `INSERT INTO ce_m_activity
	(member_id, ce_activity_id, evidence, created_at, updated_at,
	activity_on, quantity, points_per_unit, description)
	VALUES("%v", "%v", "%v", NOW(), NOW(), "%v", "%v", "%v", "%v")`
	query = fmt.Sprintf(query, a.MemberID, a.ActivityID, a.Evidence, a.Date, a.Quantity, a.UnitCredit, a.Description)

	// Get result of the the query execution...
	r, err := datastore.MySQL.Session.Exec(query)
	if err != nil {
		return 0, err
	}

	// Get the new id...
	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateMemberActivity updates an existing member activity record in the MySQL db
func UpdateMemberActivity(a MemberActivityRow) error {

	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return err
	}

	// Look up the value of this type of activity
	uc, err := activities.ActivityUnitCredit(a.ActivityID)
	if err != nil {
		return err
	}
	a.UnitCredit = uc

	query := `UPDATE ce_m_activity SET
	ce_activity_id= "%v",
	evidence= "%v",
	updated_at = NOW(),
	activity_on = "%v",
	quantity= "%v",
	points_per_unit= "%v",
	description = "%v"
	WHERE id = %v
	LIMIT 1`
	query = fmt.Sprintf(query, a.ActivityID, a.Evidence, a.Date, a.Quantity, a.UnitCredit, a.Description, a.ID)
	_, err = datastore.MySQL.Session.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// Save a recurring activity
//func (a *members.RecurringActivity) Save() error {
//
//	return nil
//}
