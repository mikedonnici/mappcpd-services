package cpd

import (
	"fmt"
	"log"
	"time"

	"github.com/mappcpd/web-services/internal/activity"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/nleof/goyesql"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

var queries = goyesql.MustParseFile("queries.sql")

// MemberActivity represents an instance of an activity recorded by a member - ie a CPD diary entry
type MemberActivity struct {
	ID          int                       `json:"id" bson:"id"`
	MemberID    int                       `json:"memberId" bson:"memberId"`
	CreatedAt   time.Time                 `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time                 `json:"updatedAt" bson:"updatedAt"`
	Date        string                    `json:"date" bson:"date"`
	DateISO     time.Time                 `json:"dateISO" bson:"dateISO"`
	Credit      float64                   `json:"credit" bson:"credit"`
	Description string                    `json:"description" bson:"description"`
	Evidence    bool                      `json:"evidence" bson:"evidence"`
	Category    activity.ActivityCategory `json:"category" bson:"category"`
	Activity    activity.Activity         `json:"activity" bson:"activity"`
	Type        activity.ActivityType     `json:"type" bson:"type"`
	CreditData  activity.ActivityCredit   `json:"creditData" bson:"creditData"`
}

// MemberActivityInput contains fields required to add / update a Member Activity.
type MemberActivityInput struct {
	ID          int     `json:"ID"`
	MemberID    int     `json:"memberId"`
	ActivityID  int     `json:"activityId" validate:"required,min=1"`
	TypeID      int     `json:"typeId" validate:"required,min=1"`
	Date        string  `json:"date" validate:"required"`
	Quantity    float64 `json:"quantity" validate:"required"`
	UnitCredit  float64 `json:"unitCredit"`
	Description string  `json:"description" validate:"required"`
	Evidence    bool    `json:"evidence"`
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

	// evidence is 0 or 1 in the database, we want a boolean
	var evidence int

	query := queries["select-member-activity"] + ` WHERE cma.id = ?`
	err := datastore.MySQL.Session.QueryRow(query, id).Scan(
		&a.ID,
		&a.MemberID,
		&a.Date,
		&a.Description,
		&evidence,
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

	if evidence == 1 {
		a.Evidence = true
	}

	a.DateISO, err = time.Parse("2006-01-02", a.Date)
	if err != nil {
		log.Printf("Error creating ISODate: %s", err.Error())
	}

	return &a, nil
}

// MemberActivitiesByMemberID fetches activities for a particular member
func MemberActivitiesByMemberID(memberID int) ([]MemberActivity, error) {

	memberActivities := MemberActivities{}

	query := queries["select-member-activity"] + ` WHERE member_id = ? ORDER BY activity_on DESC`
	rows, err := datastore.MySQL.Session.Query(query, memberID)
	if err != nil {
		return memberActivities, err
	}
	defer rows.Close()

	for rows.Next() {

		a := MemberActivity{}

		err := rows.Scan(
			&a.ID,
			&a.MemberID,
			&a.Date,
			&a.Description,
			&a.Evidence,
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

		memberActivities = append(memberActivities, a)
	}

	return memberActivities, nil
}

// MemberActivitiesQuery allows for any filter clause
func MemberActivitiesQuery(sqlClause string) ([]MemberActivity, error) {

	memberActivities := MemberActivities{}

	query := queries["select-member-activity"] + ` ` + sqlClause
	rows, err := datastore.MySQL.Session.Query(query)
	if err != nil {
		return memberActivities, err
	}
	defer rows.Close()

	for rows.Next() {

		a := MemberActivity{}

		err := rows.Scan(
			&a.ID,
			&a.MemberID,
			&a.Date,
			&a.Description,
			&a.Evidence,
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

		memberActivities = append(memberActivities, a)
	}

	return memberActivities, nil
}

// AddMemberActivity inserts a new member activity in the MySQL db and returns the new id on success.
func AddMemberActivity(a MemberActivityInput) (int, error) {

	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return 0, err
	}

	// Look up the credit-per-unit for this type of activity...
	uc, err := activity.ActivityUnitCredit(a.ActivityID)
	if err != nil {
		return 0, err
	}
	a.UnitCredit = uc

	// evidence is passed in as bool but in the database stored as 0/1
	var evidence int
	if a.Evidence == true {
		evidence = 1
	}

	query := `INSERT INTO ce_m_activity
	(member_id, ce_activity_id, ce_activity_type_id, evidence, created_at, updated_at,
	activity_on, quantity, points_per_unit, description)
	VALUES("%v", "%v", "%v", "%v", NOW(), NOW(), "%v", "%v", "%v", "%v")`
	query = fmt.Sprintf(query, a.MemberID, a.ActivityID, a.TypeID, evidence, a.Date, a.Quantity, a.UnitCredit, a.Description)

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
	uc, err := activity.ActivityUnitCredit(a.ActivityID)
	if err != nil {
		return err
	}
	a.UnitCredit = uc

	// evidence is passed in as bool but in the database stored as 0/1
	var evidence int
	if a.Evidence == true {
		evidence = 1
	}

	query := `UPDATE ce_m_activity SET ce_activity_id= "%v", ce_activity_type_id= "%v", evidence= "%v",
    updated_at = NOW(), activity_on = "%v", quantity= "%v", points_per_unit= "%v", description = "%v"
    WHERE id = %v LIMIT 1`
	query = fmt.Sprintf(query, a.ActivityID, a.TypeID, evidence, a.Date, a.Quantity, a.UnitCredit, a.Description, a.ID)
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

	query := `SELECT id FROM ce_m_activity WHERE member_id = "%v" AND ce_activity_id = "%v" AND 
	ce_activity_type_id = "%v" AND activity_on = "%v" AND description = "%v" LIMIT 1`
	query = fmt.Sprintf(query, a.MemberID, a.ActivityID, a.TypeID, a.Date, a.Description)

	row := datastore.MySQL.Session.QueryRow(query)
	row.Scan(&dupId)

	return dupId
}

// DeleteMemberActivity deletes a member activity record
func DeleteMemberActivity(memberID, activityID int) error {

	query := `DELETE FROM ce_m_activity WHERE member_id = %d AND id = %d LIMIT 1`
	query = fmt.Sprintf(query, memberID, activityID)
	_, err := datastore.MySQL.Session.Exec(query)
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
