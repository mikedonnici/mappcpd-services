package graphql

import (
	"github.com/cardiacsociety/web-services/internal/activity"
	"github.com/graphql-go/graphql"
)

// activityTypeMap is a local version of activities.Type, to remove to the sql.NullInt64
type activityTypeMap struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// activitiesData returns a list of activity types
func activitiesData() ([]activity.Activity, error) {
	return activity.All(DS)
}

// activityTypesData returns sub types for an activity
func activityTypesData(activityID int) ([]activity.Type, error) {
	return activity.Types(DS, activityID)
}

// ActivitiesQueryField resolves queries for activities (activity types)
var ActivitiesQueryField = &graphql.Field{
	Description: "Fetches a list of activity types.",
	Type:        graphql.NewList(activityQueryObject),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return activitiesData()
	},
}

// activityQueryObject defines the fields (properties) of an activity
var activityQueryObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "activity",
	Description: "Activity describes a group of related activity types. This is the entity that includes the credit " +
		"value and caps for the activity (types) contained within.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of the activity",
		},
		"code": &graphql.Field{
			Type:        graphql.String,
			Description: "The code representing the activity",
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the activity",
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "A description of the activity",
		},
		"categoryId": &graphql.Field{
			Type:        graphql.Int,
			Description: "ID of the category to which the activity belongs",
		},
		"categoryName": &graphql.Field{
			Type:        graphql.String,
			Description: "ReportName of the category to which the activity belongs",
		},
		"unitId": &graphql.Field{
			Type:        graphql.Int,
			Description: "ID of the unit record used to measure the activity",
		},
		"unitName": &graphql.Field{
			Type:        graphql.String,
			Description: "ReportName of the unit used to measure the activity",
		},
		"creditPerUnit": &graphql.Field{
			Type:        graphql.Float,
			Description: "CPD credit per per unit of activity",
		},
		"types": activityTypesQueryField,
	},
})

// activityTypesQueryField resolves queries for activity types
var activityTypesQueryField = &graphql.Field{
	Description: "Fetches a list of activity types.",
	Type:        graphql.NewList(activityTypeQueryObject),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// get the activity id from the parent (activity) object
		// note .Desc is interface{} which can assert to activity
		activityID := p.Source.(activity.Activity).ID
		types, err := activityTypesData(activityID)
		if err != nil {
			return nil, nil
		}

		var xat []activityTypeMap
		for _, v := range types {
			at := activityTypeMap{}
			if v.ID > 0 {
				at.ID = v.ID
				at.Name = v.Name
				xat = append(xat, at)
			}
		}

		return xat, nil
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
