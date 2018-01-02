package schema

import "github.com/graphql-go/graphql"

// location represents one or more contact locations pertaining to a member
var memberLocationType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "location",
	Description: "A contact location belonging to a member",
	Fields: graphql.Fields{
		"order": &graphql.Field{
			Type: graphql.Int,
		},
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"address": &graphql.Field{
			Type: graphql.String,
		},
		"city": &graphql.Field{
			Type: graphql.String,
		},
		"state": &graphql.Field{
			Type: graphql.String,
		},
		"postcode": &graphql.Field{
			Type: graphql.String,
		},
		"country": &graphql.Field{
			Type: graphql.String,
		},
		"phone": &graphql.Field{
			Type: graphql.String,
		},
		"fax": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
	},
})
