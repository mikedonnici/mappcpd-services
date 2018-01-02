package data

import (
	"github.com/mappcpd/web-services/internal/activities"
)

// ActivityType is a simpler representation of an activity type
type ActivityType struct {
	ID          int    `json:"id" bson:"id"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

// GetActivityTypes returns a list of activity types
func GetActivityTypes() ([]ActivityType, error) {

	var xat []ActivityType

	xa, err := activities.ActivityList()
	if err != nil {
		return nil, err
	}

	// stick into the 'flatter' value type
	for _, a := range xa {
		at := ActivityType{}
		at.ID = a.ID
		at.Code = a.Code
		at.Name = a.Name
		at.Description = a.Description
		xat = append(xat, at)
	}

	return xat, nil
}
