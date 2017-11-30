package queries

import (
	"github.com/graphql-go/graphql"
	"github.com/mappcpd/web-services/cmd/webd/graphql/data"
	"github.com/mappcpd/web-services/cmd/webd/graphql/schema/types"
	"fmt"
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
			// todo - validate token and extract id from token
			fmt.Println("Validate and extract id from token:", token)
			id := 501
			return data.GetMemberViewer(id)
		}
		return nil, nil
	},
}
