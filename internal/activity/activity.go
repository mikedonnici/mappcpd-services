package activity

import (
	"fmt"
	"runtime"

	"database/sql"

	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/utility"
	"github.com/nleof/goyesql"
	"github.com/pkg/errors"
)

// ActivityCategory is the broadest grouping of activity and is purely descriptive
type ActivityCategory struct {
	ID          int    `json:"id" bson:"id"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

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
	// Credit       ActivityCredit `json:"credit" bson:"credit"`
}

// ActivityCredit holds the detail about how the credit is calculated for the activity
// todo remove this and flatten into the Activity type - too complex
type ActivityCredit struct {
	QuantityFixed   bool    `json:"quantityFixed"`
	Quantity        float64 `json:"quantity" bson:"quantity"`
	UnitCode        string  `json:"unitCode" bson:"unitCode"`
	UnitName        string  `json:"unitName" bson:"unitName"`
	UnitDescription string  `json:"unitDescription" bson:"unitDescription"`
	UnitCredit      float64 `json:"unitCredit" bson:"unitCredit"`
}

// ActivityType represents a specific form, or example, of an Activity, Where ActivityCategory is the broadest
// descriptive attribute, ActivityType is the most specific. However, it is also purely descriptive as the numbers all
// occur in the Activity entity.
type ActivityType struct {
	ID   sql.NullInt64 `json:"id" bson:"id"` // can be NULL for old data
	Name string        `json:"name" bson:"name"`
	//Activity Activity      `json:"activity" bson:"activity"`
}

var queries = goyesql.MustParseFile("queries.sql")

// All fetches active Activity records
func All() ([]Activity, error) {
	return activityList(datastore.MySQL)
}

// AllStore fetches active Activity records from the specified datastore - used for testing
func AllStore(conn datastore.MySQLConnection) ([]Activity, error) {
	return activityList(conn)
}

// Types fetches the activity types
func Types(activityID int) ([]ActivityType, error) {
	return activityTypes(activityID, datastore.MySQL)
}

// Types fetches the activity types from the specified datastore - used for testing
func TypesStore(activityID int, conn datastore.MySQLConnection) ([]ActivityType, error) {
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

// ActivityByActivityTypeID fetches a single activity by activity type id
func ActivityByActivityTypeID(activityTypeID int) (Activity, error) {

	var a Activity

	// get the activity id
	var id int
	query := "SELECT ce_activity_id FROM ce_activity_type WHERE id = ?"
	err := datastore.MySQL.Session.QueryRow(query, activityTypeID).Scan(&id)
	if err != nil {
		function, file, line, ok := runtime.Caller(0)
		msg := utility.ErrorLocationMessage(function, file, line, ok, true)
		return a, errors.Wrap(err, msg)
	}

	return ByID(id)
}

// ActivityUnitCredit gets the credit value, per unit (eg hour, item) for a particular
// type of activity. For example, attendance at a workshop may be  measured in units
// of 'hours', each of which is worth 1 CPD credit point. It received the
// id of the activity (type) and returns the value as a float.
// Note that it will also return an error if the activity (type) is not active
func ActivityUnitCredit(id int) (float64, error) {

	var p float64
	query := `SELECT points_per_unit FROM ce_activity
		  WHERE active = 1 AND id = ?`
	err := datastore.MySQL.Session.QueryRow(query, id).Scan(&p)
	if err != nil {
		return p, err
	}

	return p, nil
}

// ActivityCreditData gets the values for the ActivityCredit properties
// for a particular activity type. This describes all of the information about
// the way an activity is credited - units, points per unit, etc.
//
// It receives an argument that is the id of the activity unit record, that is,
// from the ce_activity_unit table.
func ActivityCreditData(activityUnitID int) (ActivityCredit, error) {

	u := ActivityCredit{}
	u.QuantityFixed = false

	// Coalesce any NULL-able fields
	query := `SELECT
		COALESCE(name, ''),
		COALESCE(description, ''),
	    specify_quantity
		FROM ce_activity_unit
		WHERE id = ?`

	// temp map the specify_quantity field
	var specifyQuantity int

	err := datastore.MySQL.Session.QueryRow(query, activityUnitID).Scan(
		&u.UnitName,
		&u.UnitDescription,
		&specifyQuantity,
	)

	// MySQL table has a flag specify_quantity that tells the software if the user is allowed to input a quantity.
	// If set to zero them the unit / item is measures as a 'single item' or thing, without a quanity. For example,
	// publishing a paper - a single event.
	if specifyQuantity == 0 {
		u.QuantityFixed = true
	}

	return u, err
}

// Save an activity to MySQL
func (a Activity) Save() error {
	//query :=
	return nil
}

func activityList(conn datastore.MySQLConnection) ([]Activity, error) {

	var xa []Activity

	q := queries["select-activities"] + " WHERE a.active = 1"
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

	// map ce_activity.ce_activity_unit_id
	//var ceActivityUnitID int

	// Not using .QueryRow even though is only one row - so can share the scanActivity func
	q := queries["select-activities"] + ` WHERE a.id = ? LIMIT 1`
	rows, err := conn.Session.Query(q, id)
	if err != nil {
		return a, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanActivity(rows)
	}

	return a, nil
}

func activityTypes(activityID int, conn datastore.MySQLConnection) ([]ActivityType, error) {

	var xat []ActivityType

	query := "SELECT id, name FROM ce_activity_type WHERE active = 1 AND ce_activity_id = ?"
	rows, err := conn.Session.Query(query, activityID)
	if err != nil {
		return xat, err
	}
	defer rows.Close()

	for rows.Next() {
		at := ActivityType{}
		err := rows.Scan(&at.ID, &at.Name)
		if err != nil {
			fmt.Println(err)
		}
		xat = append(xat, at)
	}

	return xat, nil
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
