package mutations

import (
	"github.com/graphql-go/graphql"
	"github.com/mappcpd/web-services/cmd/webd/graphql/schema/types"

)


// AddMemberActivity records a new activity for a member
var AddMemberActivity = &graphql.Field{
	Name:        "AddMemberActivity",
	Description: "Add a member activity",
	Type:        types.Activity,
	Args: graphql.FieldConfigArgument{
		"position": &graphql.ArgumentConfig{
			Type:        types.PositionInput,
			Description: "A position object",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		position, ok := p.Args["position"].(map[string]interface{})
		if ok {
			dp := data.Position{}
			dp.Unpack(position)
			return data.AddPosition(dp), nil
		}
		return nil, nil
	},
}
