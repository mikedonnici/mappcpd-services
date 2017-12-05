package mutations

import (
	"github.com/graphql-go/graphql"
	"github.com/mappcpd/web-services/cmd/webd/graphql/schema/types"

	"github.com/mappcpd/web-services/cmd/webd/graphql/data"
)

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

			newId, err := data.AddMemberActivity(501, ma)
			if err != nil {
				return ma, err
			}

			// Return the newly created record

			return data.AddPosition(dp), nil
		}
		return nil, nil
	},
}
