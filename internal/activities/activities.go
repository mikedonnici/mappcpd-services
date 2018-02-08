package activities

import (
	"fmt"

	"database/sql"

	"github.com/pkg/errors"

	"github.com/mappcpd/web-services/internal/platform/datastore"
	"runtime"
)

// Activity describes a type of activity, eg online learning. This is NOT the same
// as the category which is a much broader grouping.
type Activity struct {
	ID          int            `json:"id" bson:"id"`
	Code        string         `json:"code" bson:"code"`
	Name        string         `json:"name" bson:"name"`
	Description string         `json:"description" bson:"description"`
	Credit      ActivityCredit `json:"credit" bson:"credit"`
}

// ActivityCredit holds the detail about how the credit is calculated for the activity
type ActivityCredit struct {
	QuantityFixed   bool    `json:"quantityFixed"`
	Quantity        float64 `json:"quantity" bson:"quantity"`
	UnitCode        string  `json:"unitCode" bson:"unitCode"`
	UnitName        string  `json:"unitName" bson:"unitName"`
	UnitDescription string  `json:"unitDescription" bson:"unitDescription"`
	UnitCredit      float64 `json:"unitCredit" bson:"unitCredit"`
}

// ActivityCategory stored details about the category
type ActivityCategory struct {
	ID          int    `json:"id" bson:"id"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

// ActivityType represents a further classification of an Activity into a 'type' of the activity in question.
type ActivityType struct {
	ID   sql.NullInt64 `json:"id" bson:"id"` // can be NULL for old data
	Name string        `json:"name" bson:"name"`
	Activity Activity `json:"activity: bson: "activity"`
}

//type Activities []Activity

// Activities fetches activity (types)
func Activities() ([]Activity, error) {

	var xa []Activity

	query := "SELECT id, ce_activity_unit_id, code, name, description FROM ce_activity WHERE active = 1"

	rows, err := datastore.MySQL.Session.Query(query)
	if err != nil {
		return xa, err
	}
	defer rows.Close()

	for rows.Next() {
		at := Activity{}
		// map ce_activity.ce_activity_unit_id
		var ceActivityUnitID int
		rows.Scan(&at.ID, &ceActivityUnitID, &at.Code, &at.Name, &at.Description)
		at.Credit, err = ActivityCreditData(ceActivityUnitID)
		if err != nil {
			return xa, err
		}
		xa = append(xa, at)
	}

	return xa, nil
}

// ActivityTypes fetches activity types
func ActivityTypes() ([]ActivityType, error) {

	var xat []ActivityType

	query := "SELECT id, name FROM ce_activity_type WHERE active = 1"

	rows, err := datastore.MySQL.Session.Query(query)
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

// ActivityTypesByActivity fetches activity sub-types for the activity designated by activityID
func ActivityTypesByActivity(activityID int) ([]ActivityType, error) {

	var xat []ActivityType

	query := "SELECT id, name FROM ce_activity_type WHERE active = 1 AND ce_activity_id = ?"

	rows, err := datastore.MySQL.Session.Query(query, activityID)
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

// ActivityByID fetches a single activity by id
func ActivityByID(id int) (Activity, error) {

	var a Activity

	// map ce_activity.ce_activity_unit_id
	var ceActivityUnitID int

	query := "SELECT id, ce_activity_unit_id, code, name, description FROM ce_activity WHERE active = 1 AND id = ?"
	err := datastore.MySQL.Session.QueryRow(query, id).Scan(
		&a.ID,
		&ceActivityUnitID,
		&a.Code,
		&a.Name,
		&a.Description,
	)
	if err != nil {
		return a, err
	}

	// Add credit info
	a.Credit, err = ActivityCreditData(ceActivityUnitID)
	if err != nil {
		return a, err
	}

	return a, nil
}

// ActivityByActivityTypeID fetches a single activity by activity type id
func ActivityByActivityTypeID(activityTypeID int) (Activity, error) {

	var a Activity

	// get the activity id
	var id int
	query := "SELECT ce_activity_id FROM ce_activity_type WHERE id = ?"
	err := datastore.MySQL.Session.QueryRow(query, activityTypeID).Scan(&id)
	if err != nil {

		function, file, line, _ := runtime.Caller(0)
		msg := fmt.Sprintf("File: %s  Function: %s Line: %d", file, runtime.FuncForPC(function).Name(), line)

		return a, errors.Wrap(err, msg)
	}

	return ActivityByID(id)
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
