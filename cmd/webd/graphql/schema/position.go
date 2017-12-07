package schema

import "github.com/graphql-go/graphql"

// position represents a position held by a member
var positionType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "position",
	Description: "A position or affiliation with a council, committee or group",
	Fields: graphql.Fields{
		"orgCode": &graphql.Field{
			Type: graphql.String,
		},
		"orgName": &graphql.Field{
			Type: graphql.String,
		},
		"code": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"startDate": &graphql.Field{
			Type: graphql.String,
		},
		"endDate": &graphql.Field{
			Type: graphql.String,
		},
	},
})
