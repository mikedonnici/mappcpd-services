package types

import "github.com/graphql-go/graphql"

// Contact represents a contact 'card' - that is, a single contact record that pertains to a Member.
var Contact = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Contact",
	Description: "A contact record belonging to a member",
	Fields: graphql.Fields{
		"emailPrimary": &graphql.Field{
			Type: graphql.String,
		},
		"emailSecondary": &graphql.Field{
			Type: graphql.String,
		},
		"mobile": &graphql.Field{
			Type: graphql.String,
		},
		"locations": &graphql.Field{
			Type: graphql.NewList(Location),
		},
	},
})
