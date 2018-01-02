package schema

import "github.com/graphql-go/graphql"

// qualification represents a qualification obtained by the member
var qualificationType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "qualification",
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
