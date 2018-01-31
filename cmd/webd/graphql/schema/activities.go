package schema

import (
	"github.com/graphql-go/graphql"

	"github.com/mappcpd/web-services/internal/activities"
)

// Activity is a trimmer version of an activities.Activity
type Activity struct {
	ID          int    `json:"id" bson:"id"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

// GetActivityTypes returns a list of activity types
func GetActivityTypes() ([]Activity, error) {

	var xat []Activity

	xa, err := activities.ActivityList()
	if err != nil {
		return nil, err
	}

	// stick into the 'flatter' value type
	for _, a := range xa {
		at := Activity{}
		at.ID = a.ID
		at.Code = a.Code
		at.Name = a.Name
		at.Description = a.Description
		xat = append(xat, at)
	}

	return xat, nil
}

// Activities query field fetches activity types
var ActivitiesQuery = &graphql.Field{
	Description: "Fetches a list of activity types.",
	Type:        graphql.NewList(activity),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return GetActivityTypes()
	},
}

// activity (object) defines the fields (properties) of an activity
var activity = graphql.NewObject(graphql.ObjectConfig{
	Name: "activity",
	Description: "Represents a type of activity that can be recorded by a member (memberActivity). " +
		"This query should be used to create select lists, etc.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of the activity type, required when adding member activities",
		},
		"code": &graphql.Field{
			Type:        graphql.String,
			Description: "The code representing the activity type",
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the activity type - use for select lists etc",
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "A more detailed description of the activity type",
		},
	},
})
