package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/mappcpd/web-services/cmd/webd/graphql/data"
)

// Activities query field fetches the list of activity types
var Activities = &graphql.Field{
	Description: "Fetches a list of activity types.",
	Type:        graphql.NewList(activityType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return data.GetActivityTypes()
	},
}
