package graphql

import (
	"fmt"
	"os"
	"strings"

	"net/http"

	"github.com/graphql-go/handler"
	"github.com/rs/cors"

	"github.com/mappcpd/web-services/cmd/webd/graphql/schema"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// Start fires up the GraphQL server
func Start(port string) {

	// todo: should this even be here? Shouldn't the internal packages handle the connection?
	datastore.Connect()

	schema, err := schema.Create()
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
