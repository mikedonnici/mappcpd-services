package graphql

import (
	"fmt"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/mappcpd/web-services/cmd/webd/graphql/schema"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/rs/cors"
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
				//"members": queries.Members,
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
	//http.Handle("/graphql", optionsCheck(ch))
	fmt.Println("GraphQL server listening at", os.Getenv("MAPPCPD_API_URL")+":"+port+"/graphql")
	http.ListenAndServe(":"+port, nil)
}

// optionsCheck checks if the request is an OPTIONS type, and if so returns an options handler.
// Otherwise it just calls the
//func optionsCheck(h http.Handler) http.Handler {
//
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if r.Method == "OPTIONS" {
//			optionsHandler(w,r)
//		} else {
//			h.ServeHTTP(w, r)
//		}
//
//	})
//}

// optionsHandler handles an OPTIONS request such as is made by Chrome in preflight requests.
// This is the same as the Preflight() func in the REST server
//func optionsHandler(w http.ResponseWriter, r *http.Request) {
//	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("Handling OPTIONS request")
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
//		w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
//		w.Header().Set("Content-Type", "text/plain")
//		io.WriteString(w, "Cabin crew, please arm doors and crosscheck :)")
//	})
//	h.ServeHTTP(w,r)
//}
