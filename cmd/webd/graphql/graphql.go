package graphql

import (
	"fmt"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/mappcpd/web-services/cmd/webd/graphql/schema/queries"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

func Start(port string) {

	datastore.Connect()

	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"member": queries.Member,
				//"member":  queries.Member,
				//"members": queries.Members,
			},
		})

	//rootMutation := graphql.NewObject(
	//	graphql.ObjectConfig{
	//		Name:        "RootMutation",
	//		Description: "...",
	//		Fields: graphql.Fields{
	//			"addVessel":   mutations.AddVessel,
	//			"addVoyage":   mutations.AddVoyage,
	//			"addPosition": mutations.AddPosition,
	//		},
	//	})

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: rootQuery,
			//Mutation: rootMutation,
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

	http.Handle("/graphql", h)
	fmt.Println("GraphQL server listening at", os.Getenv("MAPPCPD_API_URL")+":"+port+"/graphql")
	http.ListenAndServe(":"+port, nil)
}