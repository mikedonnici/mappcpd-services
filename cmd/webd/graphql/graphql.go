package graphql

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/graphql-go/handler"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/rs/cors"
)

// Store represents the global datastore passed to internal packages by the handlers
var DS datastore.Datastore

// Start fires up the GraphQL server
func Start(port string, ds datastore.Datastore) {

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

	http.Handle("/graphql", ch)
	host := strings.Join(strings.Split(os.Getenv("MAPPCPD_API_URL"), ":")[:2], "")
	fmt.Println("GraphQL server listening at", host+":"+port+"/graphql")
	http.ListenAndServe(":"+port, nil)
}
