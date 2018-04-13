package activities

import (
	"fmt"

	"database/sql"

	"github.com/pkg/errors"

	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/utility"
	"runtime"
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
	CreditPerUnit float32 `json:"creditPerUnit" bson:"creditPerUnit"`
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
	ID       sql.NullInt64 `json:"id" bson:"id"` // can be NULL for old data
	Name     string        `json:"name" bson:"name"`
	Activity Activity      `json:"activity" bson:"activity"`
}

// Activities fetches active activity records
func Activities() ([]Activity, error) {

	var xa []Activity

	q := `SELECT
			a.id AS ActivityID,
			a.code AS ActivityCode,
			a.name AS ActivityName,
			a.description AS ActivityDescription,
			a.ce_activity_category_id AS ActivityCategoryID,
			c.name AS ActivityCategoryName,
			a.ce_activity_unit_id AS ActivityUnitID,
    		u.name AS ActivityUnitName,
    		a.points_per_unit AS CreditPerUnit
		  FROM
			ce_activity a
				LEFT JOIN
			ce_activity_category c ON a.ce_activity_category_id = c.id
				LEFT JOIN
			ce_activity_unit u ON a.ce_activity_unit_id = u.id
		  WHERE
			a.active = 1`

	rows, err := datastore.MySQL.Session.Query(q)
	if err != nil {
		return xa, err
	}
	defer rows.Close()

	for rows.Next() {
		a := Activity{}

		rows.Scan(
			&a.ID,
			&a.Code,
			&a.Name,
			&a.Description,
			&a.CategoryID,
			&a.CategoryName,
			&a.UnitID,
			&a.UnitName,
			&a.CreditPerUnit,
		)

		// More detail about the way the activity is credited
		//// is stored in the Credit field...
		//a.Credit, err = ActivityCreditData(a.UnitID)
		//if err != nil {
		//	return xa, err
		//}

		xa = append(xa, a)
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
	//a.Credit, err = ActivityCreditData(ceActivityUnitID)
	//if err != nil {
	//	return a, err
	//}

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
		function, file, line, ok := runtime.Caller(0)
		msg := utility.ErrorLocationMessage(function, file, line, ok, true)
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
