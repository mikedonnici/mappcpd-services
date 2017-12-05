package types

import (
	"github.com/graphql-go/graphql"
)

// MemberActivity represents a Member activity record (not activity type record)
var Activity = graphql.NewObject(graphql.ObjectConfig{
	Name:        "MemberActivity",
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

// MemberActivityInput is an input type for a member activity
var MemberActivityInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "MemberActivityInput",
	Description: "Member activity input type",
	Fields: graphql.InputObjectConfigFieldMap{

		// if id is supplied then it is an edit
		"id": &graphql.InputObjectFieldConfig{
			Type:        graphql.Int,
			Description: "Optional id of the member activity record - if supplied then will update existing.",
		},

		"date": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.DateTime},
			Description: "The date of the activity",
		},

		"credit": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.Float},
			Description: "Value or credit for the activity",
		},

		"categoryId": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.Int},
			Description: "The activity category id",
		},

		"typeId": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.Int},
			Description: "The activity type id",
		},

		"description": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.String},
			Description: "The specifics of the activity described by the member",
		},
	},
})
