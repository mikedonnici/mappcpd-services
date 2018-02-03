package schema

import (
	"github.com/graphql-go/graphql"

	"github.com/mappcpd/web-services/internal/activities"
)

// activity is a trimmer version of an activities.activity
type activity struct {
	ID          int    `json:"id" bson:"id"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

// activitiesData returns a list of activity types
func activitiesData() ([]activity, error) {

	var xla []activity

	xa, err := activities.ActivityList()
	if err != nil {
		return nil, err
	}

	// map to local type
	for _, a := range xa {
		at := activity{}
		at.ID = a.ID
		at.Code = a.Code
		at.Name = a.Name
		at.Description = a.Description
		xla = append(xla, at)
	}

	return xla, nil
}

// activityTypesData returns sub types for an activity
func activityTypesData(activityID int) ([]activities.ActivityType, error) {
	return activities.ActivityTypes(activityID)
}

// activitiesQueryField resolves queries for activities (activity types)
var activitiesQueryField = &graphql.Field{
	Description: "Fetches a list of activity types.",
	Type:        graphql.NewList(activityQueryObject),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return activitiesData()
	},
}

// activityQueryObject defines the fields (properties) of an activity
var activityQueryObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "activity",
	Description: "Represents a type of activity that can be recorded by a member (memberActivityQueryField). " +
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
		"types": activityTypesQueryField,
	},
})

// activityTypesQueryField resolves queries for activities sub-types
var activityTypesQueryField = &graphql.Field{
	Description: "Fetches a list of activity sub-types.",
	Type:        graphql.NewList(activityTypeQueryObject),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// get the activity id from the parent (activity) object
		// note .Source is interface{} which can assert to activity
		id := p.Source.(activity).ID
		types, err := activityTypesData(id)
		if err != nil {
			return nil, nil
		}

		return types, nil
	},
}

// activityTypeQueryObject defines the fields (properties) of an activity sub-type
var activityTypeQueryObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "activityType",
	Description: "Activity sub-types or examples.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of the activity sub-type, required when adding member activities",
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the activity sub-type",
		},
	},
})
