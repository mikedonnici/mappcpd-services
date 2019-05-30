package activity

import (
	"database/sql"
	"fmt"
	"runtime"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/utility"
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
	MaxCredit     float64 `json:"maxCredit" bson:"maxCredit"`
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
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

// All fetches active Activity records from the specified datastore - used for testing
func All(ds datastore.Datastore) ([]Activity, error) {
	return activityList(ds)
}

// Types fetches the activity types from the specified datastore - used for testing
func Types(ds datastore.Datastore, activityID int) ([]Type, error) {
	return activityTypes(ds, activityID)
}

// ByID fetches an activity from the specified datastore - used for testing
func ByID(ds datastore.Datastore, id int) (Activity, error) {
	return activityByID(ds, id)
}

// ByTypeID fetches an activity by type id from the specified store - used for testing
func ByTypeID(ds datastore.Datastore, typeID int) (Activity, error) {
	return activityByTypeID(ds, typeID)
}

// CreditPerUnit retrieves the credit per unit for an activity from the specified store - used for testing
func CreditPerUnit(ds datastore.Datastore, activityID int) (float64, error) {
	return activityCreditPerUnit(ds, activityID)
}

func activityList(ds datastore.Datastore) ([]Activity, error) {

	var xa []Activity

	q := queries["select-activities"] + " WHERE a.active = 1"
	rows, err := ds.MySQL.Session.Query(q)
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

func activityByID(ds datastore.Datastore, id int) (Activity, error) {

	var a Activity

	q := queries["select-activities"] + ` WHERE a.id = ? LIMIT 1`
	rows, err := ds.MySQL.Session.Query(q, id) // not using .QueryRow so can share scanActivity func
	if err != nil {
		return a, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanActivity(rows)
	}

	return a, nil
}

func activityByTypeID(ds datastore.Datastore, id int) (Activity, error) {

	var activityID int
	query := "SELECT ce_activity_id FROM ce_activity_type WHERE id = ?"
	err := ds.MySQL.Session.QueryRow(query, id).Scan(&activityID)
	if err != nil {
		function, file, line, ok := runtime.Caller(0)
		msg := utility.ErrorLocationMessage(function, file, line, ok, true)
		return Activity{}, errors.Wrap(err, msg)
	}

	return activityByID(ds, activityID)
}

func activityTypes(ds datastore.Datastore, activityID int) ([]Type, error) {
	var xat []Type

	q := fmt.Sprintf(queries["select-activity-types"], activityID)
	rows, err := ds.MySQL.Session.Query(q)
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

func activityCreditPerUnit(ds datastore.Datastore, id int) (float64, error) {
	var c float64
	query := `SELECT points_per_unit FROM ce_activity WHERE active = 1 AND id = ?`
	err := ds.MySQL.Session.QueryRow(query, id).Scan(&c)
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
