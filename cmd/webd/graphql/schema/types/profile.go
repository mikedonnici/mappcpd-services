package types

import "github.com/graphql-go/graphql"

var Profile = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Profile",
	Description: "A Member profile",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"firstName": &graphql.Field{
			Type: graphql.String,
		},
		"middleNames": &graphql.Field{
			Type: graphql.String,
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
		},
		"postNominal": &graphql.Field{
			Type: graphql.String,
		},
		"dateOfBirth": &graphql.Field{
			Type: graphql.String,
		},
		"contact": &graphql.Field{
			Type: Contact,
		},
	},
})
