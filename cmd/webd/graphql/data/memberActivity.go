package data

import (
	"errors"
	"fmt"
	"time"

	"github.com/mappcpd/web-services/internal/members"
	"github.com/mappcpd/web-services/internal/utility"
)

// MemberActivity is a simpler representation of the member activity than the nested one in the current REST api.
type MemberActivity struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Credit      float64   `json:"credit"`
	CategoryID  int       `json:"categoryId"`
	Category    string    `json:"category"`
	TypeID      int       `json:"typeId"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
}

// GetMemberActivities fetches activities for a member. By default it returns the entire set,
// ordered by activity date desc. Some filters have been added here for the caller's convenience.
func GetMemberActivities(memberID int, filter map[string]interface{}) ([]MemberActivity, error) {
	var xa []MemberActivity

	// This returns a nested struct which is simplified below.
	xma, err := members.MemberActivitiesByMemberID(memberID)

	// Set up date filters
	from, okFrom := filter["from"].(time.Time)
	to, okTo := filter["to"].(time.Time)
	if okFrom && okTo {
		if from.After(to) {
			return xa, errors.New("from date cannot be after to date")
		}
	}

	for _, v := range xma {

		// Apply date filters, skip to next iteration if the data is outside the range
		if okFrom {
			if v.DateISO.Before(from) {
				continue
			}
		}
		if okTo {
			if v.DateISO.After(to) {
				continue
			}
		}

		// Passed through date filters, add the record to our simplified struct
		a := MemberActivity{
			ID:          v.ID,
			Date:        v.DateISO,
			Credit:      v.Credit,
			CategoryID:  v.Category.ID,
			Category:    v.Category.Name,
			TypeID:      v.Activity.ID,
			Type:        v.Activity.Name,
			Description: v.Description,
		}
		xa = append(xa, a)
	}


	// Although less efficient, apply 'last' n filter last - otherwise it cannot be used in conjunction with
	// the date filters.
	last, ok := filter["last"].(int)
	if ok {
		// Activities are returned in reverse order so returning the 'last' n items, ie the most *recent*, means
		// slicing from the index 0. If n is greater than the total, just return the total.
		if last < len(xma) {
			xa = xa[:last]
		}
	}


	return xa, err
}

// Unpack an object into a value of type MemberActivity
func (ma *MemberActivity) Unpack(obj map[string]interface{}) error {
	if val, ok := obj["id"].(int); ok {
		ma.ID = val
	}
	if val, ok := obj["date"].(string); ok {
		d, err := utility.DateStringToTime(val)
		if err != nil {
			return err
		}
		ma.Date = d
	}
	if val, ok := obj["credit"].(float64); ok {
		ma.Credit = val
	}
	if val, ok := obj["categoryId"].(int); ok {
		ma.CategoryID = int(val)
	}
	if val, ok := obj["typeId"].(int); ok {
		ma.TypeID = int(val)
	}
	if val, ok := obj["description"].(string); ok {
		ma.Description = val
	}

	return nil
}

// GetMemberActivity fetches a single activities by id.
// It verifies that the activity is owned by the member by memberID.
func GetMemberActivity(memberID, activityID int) (MemberActivity, error) {

	var a MemberActivity

	// This returns a nested struct which we can simplify
	ma, err := members.MemberActivityByID(activityID)
	if err != nil {
		return a, err
	}

	// Verify owner match
	if ma.MemberID != memberID {
		msg := fmt.Sprintf("MemberActivity with id %v does not belong to member with id %v", activityID, memberID)
		return a, errors.New(msg)
	}

	a.ID = ma.ID
	a.Date = ma.DateISO
	a.Credit = ma.Credit
	a.CategoryID = ma.Category.ID
	a.Category = ma.Category.Name
	a.TypeID = ma.Activity.ID
	a.Type = ma.Activity.Name
	a.Description = ma.Description

	return a, nil
}

// AddMemberActivity adds a member activity
func AddMemberActivity(memberID int, memberActivity MemberActivity) (MemberActivity, error) {

	// Create the required type for the insert
	// todo: add evidence and categoryId
	ma := members.MemberActivityRow{
		MemberID:    memberID,
		ActivityID:  memberActivity.TypeID,
		Date:        memberActivity.Date.String(),
		Quantity:    memberActivity.Credit,
		Description: memberActivity.Description,
	}

	// A return value for the new record
	var mar MemberActivity

	// This just returns the new record id, so re-fetch the member activity record
	// so that all the fields are populated for the response.
	newID, err := members.AddMemberActivity(ma)
	if err != nil {
		return mar, err
	}

	return GetMemberActivity(memberID, newID)

}

// UpdateMemberActivity adds a member activity
func UpdateMemberActivity(memberID int, memberActivity MemberActivity) (MemberActivity, error) {

	// Create the required type for the insert
	ma := members.MemberActivityRow{
		MemberID:    memberID,
		ID:          memberActivity.ID,     // id of the activity instance
		ActivityID:  memberActivity.TypeID, // id of the activity type
		Date:        memberActivity.Date.String(),
		Quantity:    memberActivity.Credit,
		Description: memberActivity.Description,
	}

	// A return value for the new record
	var mar MemberActivity

	// This just returns an error so re-fetch the member activity record
	// so that all the fields are populated for the response.
	err := members.UpdateMemberActivity(ma)
	if err != nil {
		return mar, err
	}

	return GetMemberActivity(memberID, ma.ID)

}
