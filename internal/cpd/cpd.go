package cpd

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"

	"github.com/cardiacsociety/web-services/internal/activity"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// CPD represents an instance of a cpd activity recorded by a member - ie a CPD diary entry
type CPD struct {
	ID          int               `json:"id" bson:"id"`
	MemberID    int               `json:"memberId" bson:"memberId"`
	CreatedAt   time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt" bson:"updatedAt"`
	Date        string            `json:"date" bson:"date"`
	DateISO     time.Time         `json:"dateISO" bson:"dateISO"`
	Credit      float64           `json:"credit" bson:"credit"`
	Description string            `json:"description" bson:"description"`
	Evidence    bool              `json:"evidence" bson:"evidence"`
	Category    activity.Category `json:"category" bson:"category"`
	Activity    activity.Activity `json:"activity" bson:"activity"`
	Type        activity.Type     `json:"type" bson:"type"`
	CreditData  activity.Credit   `json:"creditData" bson:"creditData"`
}

// Input contains fields required to add or update a Member Activity
type Input struct {
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

// ByID fetches a CPD record by id from the specified store - used for testing
func ByID(ds datastore.Datastore, id int) (CPD, error) {
	return cpdByID(ds, id)
}

// ByMemberID fetches all cpd belonging to a member from the specified store - used for testing
func ByMemberID(ds datastore.Datastore, memberID int) ([]CPD, error) {
	return cpdByMemberID(ds, memberID)
}

// Query runs the base cpd query with any filter clause
func Query(ds datastore.Datastore, sqlClause string) ([]CPD, error) {
	return cpdQuery(ds, sqlClause)
}

// Add inserts a new cpd record into the specified datastore, and returns the new id - used for testing
func Add(ds datastore.Datastore, a Input) (int, error) {
	return add(ds, a)
}

// Update updates a cpd record in the specified store - used for testing
func Update(ds datastore.Datastore, a Input) error {
	return update(ds, a)
}

// DuplicateOf returns the id of a duplicate member activity, or 0 if not found - from the specified store
func DuplicateOf(ds datastore.Datastore, a Input) (int, error) {
	return duplicateOf(ds, a)
}

// Delete ensures the record is owned by MemberID before deleting from specified datastore - used for testing
func Delete(ds datastore.Datastore, memberID, activityID int) error {
	return delete(ds, memberID, activityID)
}

func cpdByID(ds datastore.Datastore, id int) (CPD, error) {

	a := CPD{}
	var evidence int // stored as 0/1 in db - translate to bool

	query := Queries["select-member-activity"] + ` WHERE cma.id = ?`
	err := ds.MySQL.Session.QueryRow(query, id).Scan(
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
		return a, errors.Wrap(err, "scan error")
	}

	if evidence == 1 {
		a.Evidence = true
	}

	a.DateISO, err = time.Parse("2006-01-02", a.Date)
	if err != nil {
		log.Printf("Error creating ISODate: %s", err.Error())
	}

	return a, nil
}

func cpdByMemberID(ds datastore.Datastore, id int) ([]CPD, error) {

	var xc []CPD

	query := Queries["select-member-activity"] + ` WHERE member_id = ? ORDER BY activity_on DESC`
	rows, err := ds.MySQL.Session.Query(query, id)
	if err != nil {
		return xc, err
	}
	defer rows.Close()

	for rows.Next() {

		c := CPD{}
		var evidence int // stored as 0/1 in db - translate to bool

		err := rows.Scan(
			&c.ID,
			&c.MemberID,
			&c.Date,
			&c.Description,
			&evidence,
			&c.Credit,
			&c.CreditData.Quantity,
			&c.CreditData.UnitName,
			&c.CreditData.UnitCredit,
			&c.Category.ID,
			&c.Category.Name,
			&c.Category.Description,
			&c.Activity.ID,
			&c.Activity.Code,
			&c.Activity.Name,
			&c.Activity.Description,
			&c.Type.ID,
			&c.Type.Name,
		)
		if err != nil {
			fmt.Println(err)
		}

		if evidence == 1 {
			c.Evidence = true
		}

		xc = append(xc, c)
	}

	return xc, nil
}

func cpdQuery(ds datastore.Datastore, clause string) ([]CPD, error) {

	var xc []CPD

	query := Queries["select-member-activity"] + ` ` + clause
	rows, err := ds.MySQL.Session.Query(query)
	if err != nil {
		return xc, err
	}
	defer rows.Close()

	for rows.Next() {

		c := CPD{}
		var evidence int // stored as 0/1 in db - translate to bool

		err := rows.Scan(
			&c.ID,
			&c.MemberID,
			&c.Date,
			&c.Description,
			&evidence,
			&c.Credit,
			&c.CreditData.Quantity,
			&c.CreditData.UnitName,
			&c.CreditData.UnitCredit,
			&c.Category.ID,
			&c.Category.Name,
			&c.Category.Description,
			&c.Activity.ID,
			&c.Activity.Code,
			&c.Activity.Name,
			&c.Activity.Description,
			&c.Type.ID,
			&c.Type.Name,
		)
		if err != nil {
			fmt.Println(err)
		}

		if evidence == 1 {
			c.Evidence = true
		}

		xc = append(xc, c)
	}

	return xc, nil
}

func add(ds datastore.Datastore, a Input) (int, error) {

	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return 0, err
	}

	// Look up the credit-per-unit for this type of activity...
	uc, err := activity.CreditPerUnit(ds, a.ActivityID)
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

	r, err := ds.MySQL.Session.Exec(query)
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

func update(ds datastore.Datastore, a Input) error {

	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return err
	}

	uc, err := activity.CreditPerUnit(ds, a.ActivityID)
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
	_, err = ds.MySQL.Session.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// delete requires memberID to ensure ownership of the cpd record
func delete(ds datastore.Datastore, memberID, activityID int) error {
	query := `DELETE FROM ce_m_activity WHERE member_id = %d AND id = %d LIMIT 1`
	query = fmt.Sprintf(query, memberID, activityID)
	_, err := ds.MySQL.Session.Exec(query)
	return err
}

func duplicateOf(ds datastore.Datastore, a Input) (int, error) {

	var dupId int

	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return dupId, err
	}

	query := `SELECT id FROM ce_m_activity WHERE member_id = "%v" AND ce_activity_id = "%v" AND 
		ce_activity_type_id = "%v" AND activity_on = "%v" AND description = "%v" LIMIT 1`
	query = fmt.Sprintf(query, a.MemberID, a.ActivityID, a.TypeID, a.Date, a.Description)

	err = ds.MySQL.Session.QueryRow(query).Scan(&dupId)
	if err == sql.ErrNoRows {
		return dupId, nil
	}

	return dupId, err
}
