package router

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	_h "github.com/mappcpd/web-services/cmd/webd/router/handlers"
	_mw "github.com/mappcpd/web-services/cmd/webd/router/handlers/middleware"
)

func memberSubRouter() *mux.Router {

	// Middleware for Members sub-router
	r := mux.NewRouter().StrictSlash(true)

	// members routes
	members := r.PathPrefix(v1MemberBase).Subrouter()
	members.Methods("GET").Path("/").HandlerFunc(_h.Index)
	members.Methods("GET").Path("/token").HandlerFunc(_h.MembersToken)
	members.Methods("OPTIONS").Path("/token").HandlerFunc(_h.Preflight)
	members.Methods("GET").Path("/profile").HandlerFunc(_h.MembersProfile)
	members.Methods("GET").Path("/activities").HandlerFunc(_h.MembersActivities)
	members.Methods("GET").Path("/activities/{id:[0-9]+}").HandlerFunc(_h.MembersActivitiesID)
	members.Methods("POST").Path("/activities").HandlerFunc(_h.MembersActivitiesAdd)
	members.Methods("PUT").Path("/activities/{id:[0-9]+}").HandlerFunc(_h.MembersActivitiesUpdate)
	members.Methods("GET").Path("/activities/recurring").HandlerFunc(_h.MembersActivitiesRecurring)
	members.Methods("POST").Path("/activities/recurring").HandlerFunc(_h.MembersActivitiesRecurringAdd)

	members.Methods("OPTIONS").Path("/activities/recurring/{_id}").HandlerFunc(_h.Preflight)
	members.Methods("DELETE").Path("/activities/recurring/{_id}").HandlerFunc(_h.MembersActivitiesRecurringRemove)

	members.Methods("OPTIONS").Path("/activities/recurring/{_id}/recorder").HandlerFunc(_h.Preflight)
	members.Methods("POST").Path("/activities/recurring/{_id}/recorder").HandlerFunc(_h.MembersActivitiesRecurringRecorder)

	members.Methods("GET").Path("/evaluations").HandlerFunc(_h.MembersEvaluation)

	return members
}

// memberMiddleware wraps the member sub router with appropriate middleware
func memberMiddleware(r *mux.Router) *negroni.Negroni {

	// Recovery from panic
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false // don't print the stack

	n := negroni.New()
	n.Use(recovery)
	n.Use(negroni.HandlerFunc(_mw.ValidateToken))
	n.Use(negroni.HandlerFunc(_mw.MemberScope))
	n.Use(negroni.NewLogger())
	n.Use(negroni.Wrap(r))

	return n
}
