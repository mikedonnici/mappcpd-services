package router

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/mappcpd/web-services/cmd/webd/router/handlers"
	"github.com/mappcpd/web-services/cmd/webd/router/middleware"
)

func generalSubRouter() *mux.Router {

	// Middleware for General sub-router just need a valid token
	// as these are used by both admin and member scope
	r := mux.NewRouter().StrictSlash(true)

	// general routes
	general := r.PathPrefix(v1GeneralBase).Subrouter()

	// Activity (types)
	general.Methods("GET").Path("/activities").HandlerFunc(handlers.Activities)
	general.Methods("GET").Path("/activities/{id:[0-9]+}").HandlerFunc(handlers.ActivitiesID)

	// Resources
	general.Methods("GET").Path("/resources/{id:[0-9]+}").HandlerFunc(handlers.ResourcesID)
	general.Methods("POST").Path("/resources").HandlerFunc(handlers.ResourcesCollection)
	general.Methods("GET").Path("/resources/latest/{n:[0-9]+}").HandlerFunc(handlers.ResourcesLatest)

	// Modules
	general.Methods("GET").Path("/modules/{id:[0-9]+}").HandlerFunc(handlers.ModulesID)
	general.Methods("POST").Path("/modules").HandlerFunc(handlers.ModulesCollection)

	// Attachments
	general.Methods("OPTIONS").Path("/attachments/putrequest").HandlerFunc(handlers.Preflight)
	general.Methods("GET").Path("/attachments/putrequest").HandlerFunc(handlers.S3PutRequest)

	return general
}

// generalMiddleware applies required middleware to 'general' endpoints
func generalMiddleware(r *mux.Router) *negroni.Negroni {

	// Recovery from panic
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false // don't print the stack

	n := negroni.New()
	n.Use(recovery)
	n.Use(negroni.HandlerFunc(middleware.ValidateToken))
	n.Use(negroni.NewLogger())
	n.Use(negroni.Wrap(r))

	return n
}
