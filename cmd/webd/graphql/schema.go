package graphql

import (
	"github.com/graphql-go/graphql"
)

// CreateSchema generates the GraphQL schema starting with the root nodes
func CreateSchema() (graphql.Schema, error) {

	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Query",
			Description: "Root query",
			Fields: graphql.Fields{
				"member":     Query,
				"activities": ActivitiesQueryField,
				"events":     EventsQueryField,
			},
		})

	rootMutation := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Mutation",
			Description: "Root mutation",
			Fields: graphql.Fields{
				"member": Mutation,
			},
		})

	cfg := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}

	return graphql.NewSchema(cfg)
}
