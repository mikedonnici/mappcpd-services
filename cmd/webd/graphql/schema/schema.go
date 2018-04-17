package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/mappcpd/web-services/cmd/webd/graphql/schema/activity"
	"github.com/mappcpd/web-services/cmd/webd/graphql/schema/events"
	"github.com/mappcpd/web-services/cmd/webd/graphql/schema/member"
)

// Create generates the GraphQL schema starting with the root nodes
func Create() (graphql.Schema, error) {

	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Query",
			Description: "Root query",
			Fields: graphql.Fields{
				"member":     member.Query,
				"activities": activity.ActivitiesQueryField,
				"events":     events.EventsQueryField,
			},
		})

	rootMutation := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Mutation",
			Description: "Root mutation",
			Fields: graphql.Fields{
				"member": member.Mutation,
			},
		})

	cfg := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}

	return graphql.NewSchema(cfg)
}
