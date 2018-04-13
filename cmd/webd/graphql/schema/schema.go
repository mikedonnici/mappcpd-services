package schema

import (
	"github.com/graphql-go/graphql"
)

// Create generates the GraphQL schema starting with the root nodes
func Create() (graphql.Schema, error) {

	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Query",
			Description: "Root query",
			Fields: graphql.Fields{
				"localMember": memberQueryField,
				"activities":  activitiesQueryField,
				"events":      eventsQueryField,
			},
		})

	rootMutation := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Mutation",
			Description: "...",
			Fields: graphql.Fields{
				"localMember": memberMutationField,
			},
		})

	cfg := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}

	return graphql.NewSchema(cfg)
}
