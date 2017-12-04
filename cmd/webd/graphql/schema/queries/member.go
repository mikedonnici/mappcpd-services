package queries

import (
	"github.com/graphql-go/graphql"
	"github.com/mappcpd/web-services/cmd/webd/graphql/data"
	"github.com/mappcpd/web-services/cmd/webd/graphql/schema/types"
	"github.com/mappcpd/web-services/internal/platform/jwt"
)

// member query is a top-level 'viewer' query field that ensures data is restricted to the member (user)
// identified by the token.
var Member = &graphql.Field{
	Name:        "Member",
	Description: "Viewer query for a member that requires a valid token",
	Type:        types.Member,
	Args: graphql.FieldConfigArgument{
		"token": &graphql.ArgumentConfig{
			Type:        &graphql.NonNull{OfType: graphql.String},
			Description: "Valid JSON web token",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		token, ok := p.Args["token"].(string)
		if ok {
			// Validate the token, and extract the member id
			at, err := jwt.Check(token)
			if err != nil {
				return nil, err
			}
			//fmt.Println(at.Claims)
			id := at.Claims.ID
			// At this point we have a valid token from which we've extracted an id.
			// As a final step we can verify that the id is a valid user in the system,
			// for example, that it is active. Although this is a bit redundant for each request?
			return data.GetMember(id)
		}
		return nil, nil
	},
}

// MemberActivity query field fetches a single activity that belongs to a member
var MemberActivity = &graphql.Field{
	Name:        "MemberActivity",
	Description: "Fetches a member activity.",
	Type:        types.Activity,
	Args: graphql.FieldConfigArgument{
		"token": &graphql.ArgumentConfig{
			Type:        &graphql.NonNull{OfType: graphql.String},
			Description: "Valid JSON web token",
		},
		"activityId": &graphql.ArgumentConfig{
			Type:        &graphql.NonNull{OfType: graphql.Int},
			Description: "ID of the member activity",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		var memberID, activityID int

		t, ok := p.Args["token"].(string)
		if ok {
			at, err := jwt.Check(t)
			if err != nil {
				return nil, err
			}
			memberID = at.Claims.ID

			activityID, ok = p.Args["activityId"].(int)
			if ok {
				return data.GetMemberActivity(memberID, activityID)
			}
		}

		return nil, nil
	},
}

