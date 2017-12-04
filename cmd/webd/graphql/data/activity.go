package data

import (
	"time"
	"errors"
	"fmt"

	"github.com/mappcpd/web-services/internal/members"
)

// Activity is a simpler representation of the member activity than the nested one in the current REST api.
type Activity struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Credit      float32   `json:"credit"`
	CategoryID  int       `json:"categoryId"`
	Category    string    `json:"category"`
	TypeID      int       `json:"typeId"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
}

// GetMemberActivities fetches activities for a member
func GetMemberActivities(memberID int) ([]Activity, error) {
	var xa []Activity

	// This returns a nested struct which we can simplify
	xma, err := members.MemberActivitiesByMemberID(memberID)
	for _, v := range xma {
		a := Activity{
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

	return xa, err
}

// GetMemberActivity fetches a single activities by id.
// It verifies that the activity is owned by the member by memberID.
func GetMemberActivity(memberID, activityID int) (Activity, error) {

	var a Activity

	// This returns a nested struct which we can simplify
	ma, err := members.MemberActivityByID(activityID)
	if err != nil {
		return a, err
	}

	// Verify owner match
	if ma.MemberID != memberID {
		msg := fmt.Sprintf("Activity with id %v does not belong to member with id %v", activityID, memberID)
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
