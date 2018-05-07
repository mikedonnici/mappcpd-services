package activity

import (
	"database/sql"
	"fmt"
	"runtime"

	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/utility"
	"github.com/pkg/errors"
)

// Activity describes a group of related activity types. This is the entity that includes the credit value
// and caps for the activity (types) contained within.
type Activity struct {
	ID            int     `json:"id" bson:"id"`
	Code          string  `json:"code" bson:"code"`
	Name          string  `json:"name" bson:"name"`
	Description   string  `json:"description" bson:"description"`
	CategoryID    int     `json:"categoryId" bson:"categoryId"`
	CategoryName  string  `json:"categoryName" bson:"categoryName"`
	UnitID        int     `json:"unitId" bson:"unitId"`
	UnitName      string  `json:"unitName" bson:"unitName"`
	CreditPerUnit float64 `json:"creditPerUnit" bson:"creditPerUnit"`
	MaxCredit     float64 `json:"creditPerUnit" bson:"creditPerUnit"`
	// Credit       Credit `json:"credit" bson:"credit"`
}

// Category is the broadest grouping of activity and is purely descriptive
type Category struct {
	ID          int    `json:"id" bson:"id"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

// Credit holds the detail about how the credit is calculated for the activity
// todo remove this and flatten into the Activity type - too complex
type Credit struct {
	QuantityFixed   bool    `json:"quantityFixed"`
	Quantity        float64 `json:"quantity" bson:"quantity"`
	UnitCode        string  `json:"unitCode" bson:"unitCode"`
	UnitName        string  `json:"unitName" bson:"unitName"`
	UnitDescription string  `json:"unitDescription" bson:"unitDescription"`
	UnitCredit      float64 `json:"unitCredit" bson:"unitCredit"`
}

// Type represents a specific form, or example, of an Activity, Where Category is the broadest
// descriptive attribute, Type is the most specific. However, it is also purely descriptive as the numbers all
// occur in the Activity entity.
type Type struct {
	ID   sql.NullInt64 `json:"id" bson:"id"` // can be NULL for old data
	Name string        `json:"name" bson:"name"`
}

// All fetches active Activity records
func All() ([]Activity, error) {
	return activityList(datastore.MySQL)
}

// AllStore fetches active Activity records from the specified datastore - used for testing
func AllStore(conn datastore.MySQLConnection) ([]Activity, error) {
	return activityList(conn)
}

// Types fetches the activity types
func Types(activityID int) ([]Type, error) {
	return activityTypes(activityID, datastore.MySQL)
}

// Types fetches the activity types from the specified datastore - used for testing
func TypesStore(activityID int, conn datastore.MySQLConnection) ([]Type, error) {
	return activityTypes(activityID, conn)
}

// ByID fetches an activity
func ByID(id int) (Activity, error) {
	return activityByID(id, datastore.MySQL)
}

// ByIDStore fetches an activity from the specified datastore - used for testing
func ByIDStore(id int, conn datastore.MySQLConnection) (Activity, error) {
	return activityByID(id, conn)
}

// ByTypeID fetches an activity by type id
func ByTypeID(typeID int) (Activity, error) {
	return activityByTypeID(typeID, datastore.MySQL)
}

// ByTypeID fetches an activity by type id from the specified store - used for testing
func ByTypeIDStore(typeID int, conn datastore.MySQLConnection) (Activity, error) {
	return activityByTypeID(typeID, conn)
}

// CreditPerUnit gets the credit value, per unit (eg hour, item) for an activity
func CreditPerUnit(activityID int) (float64, error) {
	return activityCreditPerUnit(activityID, datastore.MySQL)
}

// CreditPerUnitStore retrieves the credit per unit for an activity from the specified store - used for testing
func CreditPerUnitStore(activityID int, conn datastore.MySQLConnection) (float64, error) {
	return activityCreditPerUnit(activityID, conn)
}

func activityList(conn datastore.MySQLConnection) ([]Activity, error) {

	var xa []Activity

	q := Queries["select-activities"] + " WHERE a.active = 1"
	rows, err := conn.Session.Query(q)
	if err != nil {
		return xa, err
	}
	defer rows.Close()

	for rows.Next() {
		a, err := scanActivity(rows)
		if err != nil {
			return xa, err
		}
		xa = append(xa, a)
	}

	return xa, nil
}

func activityByID(id int, conn datastore.MySQLConnection) (Activity, error) {

	var a Activity

	q := Queries["select-activities"] + ` WHERE a.id = ? LIMIT 1`
	rows, err := conn.Session.Query(q, id) // not using .QueryRow so can share scanActivity func
	if err != nil {
		return a, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanActivity(rows)
	}

	return a, nil
}

func activityByTypeID(id int, conn datastore.MySQLConnection) (Activity, error) {

	var activityID int
	query := "SELECT ce_activity_id FROM ce_activity_type WHERE id = ?"
	err := conn.Session.QueryRow(query, id).Scan(&activityID)
	if err != nil {
		function, file, line, ok := runtime.Caller(0)
		msg := utility.ErrorLocationMessage(function, file, line, ok, true)
		return Activity{}, errors.Wrap(err, msg)
	}

	return activityByID(activityID, conn)
}

func activityTypes(activityID int, conn datastore.MySQLConnection) ([]Type, error) {

	var xat []Type

	query := "SELECT id, name FROM ce_activity_type WHERE active = 1 AND ce_activity_id = ?"
	rows, err := conn.Session.Query(query, activityID)
	if err != nil {
		return xat, err
	}
	defer rows.Close()

	for rows.Next() {
		at := Type{}
		err := rows.Scan(&at.ID, &at.Name)
		if err != nil {
			fmt.Println(err)
		}
		xat = append(xat, at)
	}

	return xat, nil
}

func activityCreditPerUnit(id int, conn datastore.MySQLConnection) (float64, error)  {
	var c float64
	query := `SELECT points_per_unit FROM ce_activity WHERE active = 1 AND id = ?`
	err := conn.Session.QueryRow(query, id).Scan(&c)
	return c, err
}

func scanActivity(rows *sql.Rows) (Activity, error) {
	a := Activity{}
	err := rows.Scan(
		&a.ID,
		&a.Code,
		&a.Name,
		&a.Description,
		&a.CategoryID,
		&a.CategoryName,
		&a.UnitID,
		&a.UnitName,
		&a.CreditPerUnit,
		&a.MaxCredit,
	)
	return a, err
}
