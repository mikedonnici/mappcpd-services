package graphql

import (
	"fmt"
	"os"

	"net/http"

	"github.com/rs/cors"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/mappcpd/web-services/cmd/webd/graphql/schema"
	"github.com/mappcpd/web-services/cmd/webd/graphql/schema/events"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// Start fires up the GraphQL server
func Start(port string) {

	datastore.Connect()

	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Query",
			Description: "Root query",
			Fields: graphql.Fields{
				"memberUser": schema.MemberUser,
				"activities": schema.Activities,
				"events": events.Query,
			},
		})

	rootMutation := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Mutation",
			Description: "...",
			Fields: graphql.Fields{
				"memberUser": schema.MemberUserInput,
			},
		})

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    rootQuery,
			Mutation: rootMutation,
		},
	)
	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Wrap handler with CORS to handle preflight requests
	ch := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler(h)

	http.Handle("/graphql", ch)
	fmt.Println("GraphQL server listening at", os.Getenv("MAPPCPD_API_URL")+":"+port+"/graphql")
	http.ListenAndServe(":"+port, nil)
}
