package types

import "github.com/graphql-go/graphql"

// Qualification represents a qualification obtained by the member
var Qualification = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Qualification",
	Description: "An academic qualification obtained by the member",
	Fields: graphql.Fields{
		"code": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"year": &graphql.Field{
			Type: graphql.String,
		},
	},
})
