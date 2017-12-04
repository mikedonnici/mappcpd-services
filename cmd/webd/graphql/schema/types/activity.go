package types

import (
	"github.com/graphql-go/graphql"
)

// Activity represents a Member activity record (not activity type record)
var Activity = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Activity",
	Description: "An activity record belonging to a member",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "id of the member activity record",
		},
		"date": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "The date of the activity",
		},
		"credit": &graphql.Field{
			Type:        graphql.Float,
			Description: "Value or credit for the activity",
		},
		"categoryId": &graphql.Field{
			Type:        graphql.String,
			Description: "The activity category id",
		},
		"category": &graphql.Field{
			Type:        graphql.String,
			Description: "The top-level category of the activity",
		},
		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "The type of activity",
		},
		"typeId": &graphql.Field{
			Type:        graphql.String,
			Description: "The activity type id",
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "The specifics of the activity described by the member",
		},
	},
})
