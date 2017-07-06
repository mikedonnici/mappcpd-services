package router

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/mappcpd/web-services/cmd/webd/router/handlers"
	"github.com/mappcpd/web-services/cmd/webd/router/middleware"
)

func memberSubRouter() *mux.Router {

	// Middleware for Members sub-router
	r := mux.NewRouter().StrictSlash(true)

	// members routes
	members := r.PathPrefix(v1MemberBase).Subrouter()
	members.Methods("GET").Path("/").HandlerFunc(handlers.Index)
	members.Methods("GET").Path("/token").HandlerFunc(handlers.MembersToken)
	members.Methods("OPTIONS").Path("/token").HandlerFunc(handlers.Preflight)
	members.Methods("GET").Path("/profile").HandlerFunc(handlers.MembersProfile)

	members.Methods("GET").Path("/activities").HandlerFunc(handlers.MembersActivities)
	members.Methods("POST").Path("/activities").HandlerFunc(handlers.MembersActivitiesAdd)

	members.Methods("GET").Path("/activities/{id:[0-9]+}").HandlerFunc(handlers.MembersActivitiesID)
	members.Methods("PUT").Path("/activities/{id:[0-9]+}").HandlerFunc(handlers.MembersActivitiesUpdate)
	members.Methods("POST").Path("/activities/{id:[0-9]+}/attachments").HandlerFunc(handlers.MembersActivitiesAttachmentAdd)

	members.Methods("GET").Path("/activities/recurring").HandlerFunc(handlers.MembersActivitiesRecurring)
	members.Methods("POST").Path("/activities/recurring").HandlerFunc(handlers.MembersActivitiesRecurringAdd)

	members.Methods("OPTIONS").Path("/activities/recurring/{_id}").HandlerFunc(handlers.Preflight)
	members.Methods("DELETE").Path("/activities/recurring/{_id}").HandlerFunc(handlers.MembersActivitiesRecurringRemove)

	members.Methods("OPTIONS").Path("/activities/recurring/{_id}/recorder").HandlerFunc(handlers.Preflight)
	members.Methods("POST").Path("/activities/recurring/{_id}/recorder").HandlerFunc(handlers.MembersActivitiesRecurringRecorder)

	members.Methods("GET").Path("/evaluations").HandlerFunc(handlers.MembersEvaluation)

	return members
}

// memberMiddleware wraps the member sub router with appropriate middleware
func memberMiddleware(r *mux.Router) *negroni.Negroni {

	// Recovery from panic
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false // don't print the stack

	n := negroni.New()
	n.Use(recovery)
	n.Use(negroni.HandlerFunc(middleware.ValidateToken))
	n.Use(negroni.HandlerFunc(middleware.MemberScope))
	n.Use(negroni.NewLogger())
	n.Use(negroni.Wrap(r))

	return n
}
