package mutations

import (
	"github.com/graphql-go/graphql"
	"github.com/mappcpd/web-services/cmd/webd/graphql/data"
	"github.com/mappcpd/web-services/cmd/webd/graphql/schema/types"
	"github.com/mappcpd/web-services/internal/platform/jwt"
)

var Member = &graphql.Field{
	Name:        "Member",
	Description: "Viewer mutation for a member that requires a valid token",
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

// AddMemberActivity records a new activity for a member
var AddMemberActivity = &graphql.Field{
	Name:        "AddMemberActivity",
	Description: "Add a member activity",
	Type:        types.Activity,
	Args: graphql.FieldConfigArgument{
		"memberActivity": &graphql.ArgumentConfig{
			Type:        types.MemberActivityInput,
			Description: "A position object",
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
