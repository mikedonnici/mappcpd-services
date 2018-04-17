package member

import (
	"time"
	"fmt"

	"github.com/pkg/errors"

	"github.com/mappcpd/web-services/internal/member/activity"
	"github.com/mappcpd/web-services/internal/activities"
	"github.com/mappcpd/web-services/internal/date"
)

// activityData is a leaner representation of members.activityData
type activityData struct {
	ID            int       `json:"id"`
	Date          string    `json:"date"`
	DateTime      time.Time `json:"dateTime"`
	Quantity      float64   `json:"quantity"`
	CreditPerUnit float64   `json:"creditPerUnit"`
	Credit        float64   `json:"credit"`
	Description   string    `json:"description"`
	ActivityID    int       `json:"activityId"`
	Activity      string    `json:"activity"`
	CategoryID    int       `json:"categoryId"`
	Category      string    `json:"category"`
	TypeID        int       `json:"typeId"`
	Type          string    `json:"type"`
	// Attachments
	//Attachments []Attachment
	// todo: remove this UploadURL is a signed URL that allows for uploading file attachments
	UploadURL string `json:"uploadUrl"`
}

// activityAttachmentData represents an file associated with a member activity
type activityAttachmentData struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

// activityInputData represents an object for mutating a member activity
type activityInputData struct {
	ID          int     `json:"id"` // ID - if present triggers and update, else record will be added
	Date        string  `json:"date"`
	Quantity    float64 `json:"quantity"`
	Description string  `json:"description"`
	ActivityID  int     `json:"activityId"`
	TypeID      int     `json:"typeId"`
}

// unpack an object into a value of type MemberActivity
func (ma *activityData) unpack(obj map[string]interface{}) error {
	if val, ok := obj["id"].(int); ok {
		ma.ID = val
	}
	if val, ok := obj["date"].(string); ok {
		ma.Date = val
		d, err := date.StringToTime(val)
		if err != nil {
			return err
		}
		ma.DateTime = d
	}
	if val, ok := obj["quantity"].(float64); ok {
		ma.Quantity = val
	}
	if val, ok := obj["creditPerUnit"].(float64); ok {
		ma.CreditPerUnit = val
	}
	if val, ok := obj["credit"].(float64); ok {
		ma.Credit = val
	}
	if val, ok := obj["categoryId"].(int); ok {
		ma.CategoryID = int(val)
	}
	if val, ok := obj["activityId"].(int); ok {
		ma.ActivityID = int(val)
	}
	if val, ok := obj["typeId"].(int); ok {
		ma.TypeID = int(val)
	}
	if val, ok := obj["description"].(string); ok {
		ma.Description = val
	}

	return nil
}

// unpack an object into a value of type MemberActivityInput
func (mai *activityInputData) unpack(obj map[string]interface{}) error {
	if val, ok := obj["id"].(int); ok {
		mai.ID = val
	}
	if val, ok := obj["date"].(string); ok {
		mai.Date = val
	}
	if val, ok := obj["quantity"].(float64); ok {
		mai.Quantity = val
	}
	if val, ok := obj["typeId"].(int); ok {
		mai.TypeID = int(val)
	}
	if val, ok := obj["description"].(string); ok {
		mai.Description = val
	}

	return nil
}

// mapActivitiesData fetches activities for a member and maps to local activityData type.
func mapActivitiesData(memberID int, filter map[string]interface{}) ([]activityData, error) {

	var xa []activityData

	// This returns a nested struct which is simplified below.
	xma, err := activity.MemberActivitiesByMemberID(memberID)

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
		a := activityData{
			ID:            v.ID,
			Date:          v.Date,
			DateTime:      v.DateISO,
			Quantity:      v.CreditData.Quantity,
			CreditPerUnit: v.CreditData.UnitCredit,
			Credit:        v.Credit,
			CategoryID:    v.Category.ID,
			Category:      v.Category.Name,
			ActivityID:    v.Activity.ID,
			Activity:      v.Activity.Name,
			TypeID:        int(v.Type.ID.Int64), // null-able field
			Type:          v.Type.Name,
			Description:   v.Description,
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

// mapActivityData verifies ownership, fetches a member activity by ID, then maps to local activityData type.
func mapActivityData(memberID, memberActivityID int) (activityData, error) {

	var a activityData

	// This returns a nested struct which we can simplify
	ma, err := activity.MemberActivityByID(memberActivityID)
	if err != nil {
		return a, err
	}

	// Verify owner match
	if ma.MemberID != memberID {
		msg := fmt.Sprintf("Member activity (id %v) does not belong to member (id %v)", memberActivityID, memberID)
		return a, errors.New(msg)
	}

	a.ID = ma.ID
	a.Date = ma.Date
	a.DateTime = ma.DateISO
	a.Quantity = ma.CreditData.Quantity
	a.CreditPerUnit = ma.CreditData.UnitCredit
	a.Credit = ma.Credit
	a.CategoryID = ma.Category.ID
	a.Category = ma.Category.Name
	a.ActivityID = ma.Activity.ID
	a.Activity = ma.Activity.Name
	a.TypeID = int(ma.Type.ID.Int64)
	a.Type = ma.Type.Name
	a.Description = ma.Description

	return a, nil
}

// addActivity adds a member activity
func addActivity(memberID int, activityInput activityInputData) (activityData, error) {

	// Create the required type for the insert
	// todo: add evidence and attachment
	ma := activity.MemberActivityInput{
		MemberID:    memberID,
		ActivityID:  activityInput.ActivityID,
		TypeID:      activityInput.TypeID,
		Date:        activityInput.Date,
		Quantity:    activityInput.Quantity,
		Description: activityInput.Description,
	}

	// A return value for the new record
	var mar activityData

	// This just returns the new record id, so re-fetch the member activity record
	// so that all the fields are populated for the response.
	newID, err := activity.AddMemberActivity(ma)
	if err != nil {
		return mar, err
	}

	return mapActivityData(memberID, newID)
}

// updateActivity updates an existing member activity record
func updateActivity(memberID int, activityInput activityInputData) (activityData, error) {

	// Create the required value
	ma := activity.MemberActivityInput{
		ID:          activityInput.ID,
		MemberID:    memberID,
		ActivityID:  activityInput.ActivityID,
		TypeID:      activityInput.TypeID,
		Date:        activityInput.Date,
		Quantity:    activityInput.Quantity,
		Description: activityInput.Description,
	}

	// A return value for the new record
	var mar activityData

	// This just returns an error so re-fetch the member activity record
	// so that all the fields are populated for the response.
	err := activity.UpdateMemberActivity(ma)
	if err != nil {
		return mar, err
	}

	return mapActivityData(memberID, ma.ID)
}

// activityDuplicateID returns the id of a matching member activity, or 0 if not found
func activityDuplicateID(memberID int, activityInput activityInputData) int {

	// Create the required value
	ma := activity.MemberActivityInput{
		ID:          activityInput.ID,
		MemberID:    memberID,
		ActivityID:  activityInput.ActivityID,
		TypeID:      activityInput.TypeID,
		Date:        activityInput.Date,
		Quantity:    activityInput.Quantity,
		Description: activityInput.Description,
	}

	return activity.DuplicateMemberActivity(ma)
}

// activityIDByTypeID returns the activity id for an activity type id
func activityIDByTypeID(activityTypeID int) (int, error) {
	a, err := activities.ActivityByActivityTypeID(activityTypeID)
	return a.ID, err
}
