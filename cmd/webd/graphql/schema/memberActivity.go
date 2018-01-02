package schema

import (
	"github.com/graphql-go/graphql"
)

// memberActivity represents a Member memberActivity record (not memberActivity type record)
var memberActivityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "memberActivity",
	Description: "An activity record belonging to a member. This is an instance on an activity *type* that is recorded " +
		"by a member as having been completed on a particular date, with additional information such as a detailed description",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of the member activity record",
		},
		"date": &graphql.Field{
			Type:        graphql.String,
			Description: "The date the activity was undertaken, as string format 'YYYY-MM-DD'.",
		},
		"dateTime": &graphql.Field{
			Type: graphql.DateTime,
			Description: "The date the activity was undertaken. Note only a date string is required, eg '2017-12-07' and " +
				"any time information is discarded. This field returns the date in RFC3339 format with the time set " +
				"to 00:00:00 UTC to facilitate date ordering and other date-related operations.",
		},
		"credit": &graphql.Field{
			Type:        graphql.Float,
			Description: "Value or credit for the memberActivity",
		},
		"categoryId": &graphql.Field{
			Type:        graphql.Int,
			Description: "The memberActivity category id",
		},
		"category": &graphql.Field{
			Type:        graphql.String,
			Description: "The top-level category of the memberActivity",
		},
		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "The type of memberActivity",
		},
		"typeId": &graphql.Field{
			Type:        graphql.Int,
			Description: "The memberActivity type id",
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "The specifics of the memberActivity described by the member",
		},
	},
})

// memberActivityInput is an input object type used as an argument for adding / updating a memberActivity
var memberActivityInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "memberActivityInput",
	Description: "An input object type used as an argument for adding / updating a memberActivity",
	Fields: graphql.InputObjectConfigFieldMap{

		// optional member activity id - if supplied then it is an update
		"id": &graphql.InputObjectFieldConfig{
			Type:        graphql.Int,
			Description: "Optional id of the member memberActivity record - if supplied then will update existing.",
		},

		"date": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.String},
			Description: "The date of the memberActivity",
		},

		"credit": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.Float},
			Description: "Value or credit for the memberActivity",
		},

		"typeId": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.Int},
			Description: "The memberActivity type id",
		},

		"description": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.String},
			Description: "The specifics of the memberActivity described by the member",
		},
	},
})
