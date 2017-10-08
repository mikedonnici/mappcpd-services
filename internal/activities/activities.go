package activities

import (
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// ActivityType describes the type of activity, eg online learning. This is NOT the same
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
	QuantityFixed bool `json:"quantityFixed"`
	Quantity        float32 `json:"quantity" bson:"quantity"`
	UnitCode        string  `json:"unitCode" bson:"unitCode"`
	UnitName        string  `json:"unitName" bson:"unitName"`
	UnitDescription string  `json:"unitDescription" bson:"unitDescription"`
	UnitCredit      float32 `json:"unitCredit" bson:"unitCredit"`
}

// ActivityCategory stored details about the category
type ActivityCategory struct {
	ID          int    `json:"id" bson:"id"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

type Activities []Activity

// ActivityList fetches a list of all the 'active' activity types
func ActivityList() (Activities, error) {

	var ats Activities

	query := "SELECT id, ce_activity_unit_id, code, name, description FROM ce_activity WHERE active = 1"

	rows, err := datastore.MySQL.Session.Query(query)
	if err != nil {
		return ats, err
	}
	defer rows.Close()

	for rows.Next() {
		at := Activity{}
		// map ce_activity.ce_activity_unit_id
		var ceActivityUnitID int
		rows.Scan(&at.ID, &ceActivityUnitID, &at.Code, &at.Name, &at.Description)
		at.Credit, err = ActivityCreditData(ceActivityUnitID)
		if err != nil {
			return ats, err
		}
		ats = append(ats, at)
	}

	return ats, nil
}

// ActivityByID fetches a single activity type by id
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

// ActivityUnitCredit gets the credit value, per unit (eg hour, item) for a particular
// type of activity. For example, attendance at a workshop may be  measured in units
// of 'hours', each of which is worth 1 CPD credit point. It received the
// id of the activity (type) and returns the value as a float.
// Note that it will also return an error if the activity (type) is not active
func ActivityUnitCredit(id int) (float32, error) {

	var p float32
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

	// temp map teh specify_quantity field
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
