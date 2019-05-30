package server

import (
	"net/http"

	"github.com/cardiacsociety/web-services/cmd/webd/graphql"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	v1AuthBase    = "/v1/auth"
	v1MemberBase  = "/v1/m"
	v1AdminBase   = "/v1/a"
	v1GeneralBase = "/v1/g"
	v1ReportBase  = "/v1/r"
	graphQLBase = "/graphql"
)

// DS represents the global datastore passed to internal packages by the handlers
var DS datastore.Datastore

// Router returns a http.Handler for all web service endpoints
func Router(ds datastore.Datastore) http.Handler {

	DS = ds

	// Router
	r := mux.NewRouter()

	// Ping and preflight, no middleware required
	r.Methods("GET").Path("/").HandlerFunc(Index)
	r.Methods("OPTIONS").HandlerFunc(Preflight)

	// Auth sub-router, no middleware required
	rAuth := AuthSubRouter(v1AuthBase)
	r.PathPrefix(v1AuthBase).Handler(rAuth)

	// Admin sub-router and middleware
	rAdmin := AdminSubRouter(v1AdminBase)               // add router...
	rAdminMiddleware := AdminMiddleware(rAdmin)         // ...plus middleware...
	r.PathPrefix(v1AdminBase).Handler(rAdminMiddleware) // ...and add to main router

	// Reports sub-router, todo: add middleware to reports router
	rReports := ReportSubRouter(v1ReportBase)
	r.PathPrefix(v1ReportBase).Handler(rReports)

	// Member sub-router
	rMember := MemberSubRouter(v1MemberBase)
	rMemberMiddleware := MemberMiddleware(rMember)
	r.PathPrefix(v1MemberBase).Handler(rMemberMiddleware)

	// General sub-router
	rGeneral := GeneralSubRouter(v1GeneralBase)
	rGeneralMiddleware := GeneralMiddleware(rGeneral)
	r.PathPrefix(v1GeneralBase).Handler(rGeneralMiddleware)

	// GraphQL
	rGraphQL := graphql.Server(ds)
	r.PathPrefix(graphQLBase).Handler(rGraphQL)

	// CORS handler - needed to add OptionsPassThrough for preflight requests which use OPTIONS http method
	//handler := cors.Default().Handler(r)
	// Todo... tighten this up - not sure if needed  with preflightHandler??
	// todo: seem to have sorted this in the graphql handler so can possible remove the Preflight handler
	// in favour of the same set tup in graphql
	handler := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"*"},
		AllowedHeaders:     []string{"*"},
		OptionsPassthrough: true,
	}).Handler(r)

	return handler
}
