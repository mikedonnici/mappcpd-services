package types

import "github.com/graphql-go/graphql"

var Member = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Member",
	Description: "A Member",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"firstName": &graphql.Field{
			Type: graphql.String,
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
		},
	},
})
