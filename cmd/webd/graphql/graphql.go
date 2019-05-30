package graphql

import (
	"net/http"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/graphql-go/handler"
	"github.com/rs/cors"
)

// DS represents the global datastore passed to internal packages by the handlers
var DS datastore.Datastore

// Server returns a handler for the GraphQL server
func Server(ds datastore.Datastore) http.Handler {

	DS = ds

	schema, err := CreateSchema()
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

	return ch
}
