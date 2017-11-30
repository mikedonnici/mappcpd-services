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
		"_id": &graphql.Field{
			Type: graphql.String,
		},
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"active": &graphql.Field{
			Type: graphql.Boolean,
		},
		"profile": profile,
		//"activities": activities,
	},
})


var profile = &graphql.Field{
	Name:        "Profile",
	Description: "Fetch member profile",
	Type:        Profile,
	// No args as we will extract the id from the token
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		src := p.Source.(data.MemberViewer)
		// todo ... security check here
		return data.GetMemberProfile(src.ID)
	},
}

//var activities = &graphql.Field{
//	Name:        "Activities",
//	Description: "Fetches member activities ",
//	Type:        Profile,
//	// Todo - add args to filter the list in some way
//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
//		src := p.Source.(data.MemberViewer)
//		fmt.Println(src.ID)
//		// todo ... security check here
//		return data.GetMemberProfile(src.ID)
//	},
//}
