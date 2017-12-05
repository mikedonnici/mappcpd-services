package types

import (
	"github.com/graphql-go/graphql"
	"github.com/mappcpd/web-services/cmd/webd/graphql/data"
)

// Member type implements the 'viewer' idea outlined here:
// https://medium.com/the-graphqlhub/graphql-and-authentication-b73aed34bbeb
// This is a top level type that contains fields for which the data pertains
// only to the member identified by the token.
var Member = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Member",
	Description: "Fields accessible by a Member user",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"active": &graphql.Field{
			Type: graphql.Boolean,
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
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"mobile": &graphql.Field{
			Type: graphql.String,
		},
		"locations": &graphql.Field{
			Type: graphql.NewList(Location),
		},
		"qualifications": &graphql.Field{
			Type: graphql.NewList(Qualification),
		},
		"positions": &graphql.Field{
			Type: graphql.NewList(Position),
		},

		// these require sub queries to fetch
		"activities": activities,

		// Mutations
		"addActivity": addMemberActivity,
	},
})

// activities field
var activities = &graphql.Field{
	Name:        "Activities",
	Description: "Fetches member activities ",
	Type:        graphql.NewList(Activity),
	// Todo - add args to filter the list in some way
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		src := p.Source.(data.Member)
		return data.GetMemberActivities(src.ID)
	},
}

// addMemberActivity records a new activity for a member
var addMemberActivity = &graphql.Field{
	Name:        "AddMemberActivity",
	Description: "Add a member activity",
	Type:        Activity,
	Args: graphql.FieldConfigArgument{
		"memberActivity": &graphql.ArgumentConfig{
			Type:        MemberActivityInput,
			Description: "A member activity input type",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		maObj, ok := p.Args["memberActivity"].(map[string]interface{})
		if ok {
			ma := data.MemberActivity{}
			ma.Unpack(maObj)

			return data.AddMemberActivity(501, ma)
		}

		return nil, nil
	},
}
