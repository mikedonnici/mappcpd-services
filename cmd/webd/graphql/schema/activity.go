package schema

import "github.com/graphql-go/graphql"

// activityType represents an activity 'type' :)
var activityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "activity",
	Description: "An activity 'type' is a classification of the different types of activities that can be recorded by members.. " +
		"This query should be used for the creation of select lists and the like, and the activityID is required to be sent " +
		"with each request to record an activity, on behalf of a member.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of the member activity record",
		},
		"code": &graphql.Field{
			Type:        graphql.String,
			Description: "The code representing the activity type.",
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the activity type - use for select lists etc.",
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "A more detailed description of the activity type.",
		},
	},
})
