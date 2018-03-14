package members

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gopkg.in/go-playground/validator.v9"

	"github.com/mappcpd/web-services/internal/activities"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/pkg/errors"
)

// MemberActivity represents an instance of an activity recorded by a member - ie a CPD diary entry
type MemberActivity struct {
	ID            int                         `json:"id" bson:"id"`
	MemberID      int                         `json:"memberId" bson:"memberId"`
	CreatedAt     time.Time                   `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time                   `json:"updatedAt" bson:"updatedAt"`
	Date          string                      `json:"date" bson:"date"`
	DateISO       time.Time                   `json:"dateISO" bson:"dateISO"`
	//Quantity      float64                     `json:"quantity" bson:"quantity"`
	//CreditPerUnit float32                     `json:"creditPerUnit" bson:"creditPerUnit"`
	Credit        float64                     `json:"credit" bson:"credit"`
	Description   string                      `json:"description" bson:"description"`
	Category      activities.ActivityCategory `json:"category" bson:"category"`
	Activity      activities.Activity         `json:"activity" bson:"activity"`
	Type          activities.ActivityType     `json:"type" bson:"type"`
	CreditData    activities.ActivityCredit   `json:"creditData" bson:"creditData"`
}

// MemberActivityInput contains fields required to add / update a Member Activity.
type MemberActivityInput struct {
	ID          int     `json:"ID"`
	MemberID    int     `json:"memberId"`
	ActivityID  int     `json:"activityId" validate:"required,min=1"`
	TypeID      int     `json:"typeId" validate:"required,min=1"`
	Evidence    int     `json:"evidence"`
	Date        string  `json:"date" validate:"required"`
	Quantity    float64 `json:"quantity" validate:"required"`
	UnitCredit  float64 `json:"unitCredit"`
	Description string  `json:"description" validate:"required"`
}

// MemberActivityAttachment contains information about a file attached to a member activity
//type MemberActivityInput struct {
//	ID          int     `json:"ID"`
//	MemberID    int     `json:"memberId"`
//	ActivityID  int     `json:"activityId" validate:"required,min=1"`
//	TypeID      int     `json:"typeId" validate:"required,min=1"`
//	Evidence    int     `json:"evidence"`
//	Date        string  `json:"date" validate:"required"`
//	Quantity    float64 `json:"quantity" validate:"required"`
//	UnitCredit  float64 `json:"unitCredit"`
//	Description string  `json:"description" validate:"required"`
//}

// MemberActivities is a collection of MemberActivity values
type MemberActivities []MemberActivity

// MemberActivityByID fetches a member activity record by id
func MemberActivityByID(id int) (*MemberActivity, error) {

	a := MemberActivity{}

	query := selectMemberActivityQuery + ` WHERE cma.id = ?`

	err := datastore.MySQL.Session.QueryRow(query, id).Scan(
		&a.ID,
		&a.MemberID,
		&a.Date,
		&a.Description,
		&a.Credit,
		&a.CreditData.Quantity,
		&a.CreditData.UnitName,
		&a.CreditData.UnitCredit,
		&a.Category.ID,
		&a.Category.Name,
		&a.Category.Description,
		&a.Activity.ID,
		&a.Activity.Code,
		&a.Activity.Name,
		&a.Activity.Description,
		&a.Type.ID,
		&a.Type.Name,
	)
	if err != nil {
		fmt.Println(errors.Wrap(err, "scan error"))
		return &a, errors.Wrap(err, "scan error")
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

// MemberActivitiesByMemberID fetches activities for a particular member
func MemberActivitiesByMemberID(memberID int) ([]MemberActivity, error) {

	activities := MemberActivities{}

	query := selectMemberActivityQuery + ` WHERE member_id = ? ORDER BY activity_on DESC`

	rows, err := datastore.MySQL.Session.Query(query, memberID)
	if err != nil {
		return activities, err
	}
	defer rows.Close()

	for rows.Next() {

		a := MemberActivity{}

		err := rows.Scan(
			&a.ID,
			&a.MemberID,
			&a.Date,
			&a.Description,
			&a.Credit,
			&a.CreditData.Quantity,
			&a.CreditData.UnitName,
			&a.CreditData.UnitCredit,
			&a.Category.ID,
			&a.Category.Name,
			&a.Category.Description,
			&a.Activity.ID,
			&a.Activity.Code,
			&a.Activity.Name,
			&a.Activity.Description,
			&a.Type.ID,
			&a.Type.Name,
		)
		if err != nil {
			fmt.Println(err)
		}

		activities = append(activities, a)
	}

	return activities, nil
}

// UpdateMemberActivityDoc updates the JSON-formatted activity record in the Doc DB (MongoDB)
func UpdateMemberActivityDoc(a *MemberActivity, w *sync.WaitGroup) {

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
func AddMemberActivity(a MemberActivityInput) (int, error) {

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
	(member_id, ce_activity_id, ce_activity_type_id, evidence, created_at, updated_at,
	activity_on, quantity, points_per_unit, description)
	VALUES("%v", "%v", "%v", "%v", NOW(), NOW(), "%v", "%v", "%v", "%v")`
	query = fmt.Sprintf(query, a.MemberID, a.ActivityID, a.TypeID, a.Evidence, a.Date, a.Quantity, a.UnitCredit, a.Description)

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

	return int(id), nil
}

// UpdateMemberActivity updates an existing member activity record in the MySQL db
func UpdateMemberActivity(a MemberActivityInput) error {

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
	ce_activity_type_id= "%v",
	evidence= "%v",
	updated_at = NOW(),
	activity_on = "%v",
	quantity= "%v",
	points_per_unit= "%v",
	description = "%v"
	WHERE id = %v
	LIMIT 1`
	query = fmt.Sprintf(query, a.ActivityID, a.TypeID, a.Evidence, a.Date, a.Quantity, a.UnitCredit, a.Description, a.ID)
	_, err = datastore.MySQL.Session.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// DuplicateMemberActivity returns the id of a duplicate member activity, or 0 if not found
func DuplicateMemberActivity(a MemberActivityInput) int {

	var dupId int

	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return dupId
	}

	query := `SELECT id FROM ce_m_activity WHERE
	member_id = "%v" AND
	ce_activity_id = "%v" AND
	ce_activity_type_id = "%v" AND
	activity_on = "%v" AND
	description = "%v"
	LIMIT 1`
	query = fmt.Sprintf(query, a.MemberID, a.ActivityID, a.TypeID, a.Date, a.Description)

	row := datastore.MySQL.Session.QueryRow(query)
	row.Scan(&dupId)

	return dupId
}

// Save a recurring activity
//func (a *members.RecurringActivity) Save() error {
//
//	return nil
//}
