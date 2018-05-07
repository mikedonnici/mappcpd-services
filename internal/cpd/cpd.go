package cpd

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"

	"github.com/mappcpd/web-services/internal/activity"
	"github.com/mappcpd/web-services/internal/platform/datastore"
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

// MemberActivityAttachment contains information about a file attached to a member activity

// ByID fetches a CPD record by id
func ByID(id int) (CPD, error) {
	return cpdByID(id, datastore.MySQL)
}

// ByIDStore fetches a CPD record by id from the specified store - used for testing
func ByIDStore(id int, conn datastore.MySQLConnection) (CPD, error) {
	return cpdByID(id, conn)
}

// ByMemberID fetches all cpd belonging to a member
func ByMemberID(memberID int) ([]CPD, error) {
	return cpdByMemberID(memberID, datastore.MySQL)
}

// ByMemberIDStore fetches all cpd belonging to a member from the specified store - used for testing
func ByMemberIDStore(memberID int, conn datastore.MySQLConnection) ([]CPD, error) {
	return cpdByMemberID(memberID, conn)
}

// Query runs the base cpd query with any filter clause
func Query(sqlClause string) ([]CPD, error) {
	return cpdQuery(sqlClause, datastore.MySQL)
}

// Query runs the base cpd query with any filter clause
func QueryStore(sqlClause string, conn datastore.MySQLConnection) ([]CPD, error) {
	return cpdQuery(sqlClause, conn)
}

// Add inserts a new cpd record and returns the new id
func Add(a Input) (int, error) {
	return add(a, datastore.MySQL)
}

// AddStore inserts a new cpd record into the specified datastore, and returns the new id - used for testing
func AddStore(a Input, conn datastore.MySQLConnection) (int, error) {
	return add(a, conn)
}

// Update updates a cpd record
func Update(a Input) error {
	return update(a, datastore.MySQL)
}

// Update updates a cpd record in the specified store - used for testing
func UpdateStore(a Input, conn datastore.MySQLConnection) error {
	return update(a, conn)
}

// DuplicateOf returns the id of a duplicate member activity, or 0 if not found
func DuplicateOf(a Input) (int, error) {
	return duplicateOf(a, datastore.MySQL)
}

// DuplicateOf returns the id of a duplicate member activity, or 0 if not found - from the specified store
func DuplicateOfStore(a Input, conn datastore.MySQLConnection) (int, error) {
	return duplicateOf(a, conn)
}

// Delete ensures the record is owned by MemberID before deleting
func Delete(memberID, activityID int) error {
	return delete(memberID, activityID, datastore.MySQL)
}

// DeleteStore ensures the record is owned by MemberID before deleting from specified datastore - used for testing
func DeleteStore(memberID, activityID int, conn datastore.MySQLConnection) error {
	return delete(memberID, activityID, conn)
}

func cpdByID(id int, conn datastore.MySQLConnection) (CPD, error) {

	a := CPD{}
	var evidence int // stored as 0/1 in db - translate to bool

	query := Queries["select-member-activity"] + ` WHERE cma.id = ?`
	err := conn.Session.QueryRow(query, id).Scan(
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

func cpdByMemberID(id int, conn datastore.MySQLConnection) ([]CPD, error) {

	var xc []CPD

	query := Queries["select-member-activity"] + ` WHERE member_id = ? ORDER BY activity_on DESC`
	rows, err := conn.Session.Query(query, id)
	if err != nil {
		return xc, err
	}
	defer rows.Close()

	for rows.Next() {

		c := CPD{}

		err := rows.Scan(
			&c.ID,
			&c.MemberID,
			&c.Date,
			&c.Description,
			&c.Evidence,
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

		xc = append(xc, c)
	}

	return xc, nil
}

func cpdQuery(clause string, conn datastore.MySQLConnection) ([]CPD, error) {

	var xc []CPD

	query := Queries["select-member-activity"] + ` ` + clause
	rows, err := conn.Session.Query(query)
	if err != nil {
		return xc, err
	}
	defer rows.Close()

	for rows.Next() {

		c := CPD{}

		err := rows.Scan(
			&c.ID,
			&c.MemberID,
			&c.Date,
			&c.Description,
			&c.Evidence,
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

		xc = append(xc, c)
	}

	return xc, nil
}

func add(a Input, conn datastore.MySQLConnection) (int, error) {

	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return 0, err
	}

	// Look up the credit-per-unit for this type of activity...
	uc, err := activity.CreditPerUnitStore(a.ActivityID, conn)
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

	r, err := conn.Session.Exec(query)
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

func update(a Input, conn datastore.MySQLConnection) error {

	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return err
	}

	uc, err := activity.CreditPerUnitStore(a.ActivityID, conn)
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
	_, err = conn.Session.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// delete requires memberID to ensure ownership of the cpd record
func delete(memberID, activityID int, conn datastore.MySQLConnection) error {
	query := `DELETE FROM ce_m_activity WHERE member_id = %d AND id = %d LIMIT 1`
	query = fmt.Sprintf(query, memberID, activityID)
	_, err := conn.Session.Exec(query)
	return err
}

func duplicateOf(a Input, conn datastore.MySQLConnection) (int, error) {

	var dupId int

	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return dupId, err
	}

	query := `SELECT id FROM ce_m_activity WHERE member_id = "%v" AND ce_activity_id = "%v" AND 
		ce_activity_type_id = "%v" AND activity_on = "%v" AND description = "%v" LIMIT 1`
	query = fmt.Sprintf(query, a.MemberID, a.ActivityID, a.TypeID, a.Date, a.Description)

	err = conn.Session.QueryRow(query).Scan(&dupId)
	if err == sql.ErrNoRows {
		return dupId, nil
	}

	return dupId, err
}
