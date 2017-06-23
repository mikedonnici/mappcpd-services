package router

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	_h "github.com/mappcpd/web-services/cmd/webd/router/handlers"
	_mw "github.com/mappcpd/web-services/cmd/webd/router/middleware"
)

// adminSubRouter adds end points for admin, and appropriate middleware
func adminSubRouter() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	admin := r.PathPrefix(v1AdminBase).Subrouter()

	admin.Methods("GET").Path("/test").HandlerFunc(_h.AdminTest)
	admin.Methods("GET").Path("/idlist").HandlerFunc(_h.AdminIDList)
	admin.Methods("GET").Path("/members").HandlerFunc(_h.AdminMembersSearch)
	admin.Methods("POST").Path("/members").HandlerFunc(_h.AdminMembersSearchPost)
	admin.Methods("GET").Path("/members/{id:[0-9]+}").HandlerFunc(_h.AdminMembersID)
	admin.Methods("POST").Path("/members/{id:[0-9]+}").HandlerFunc(_h.AdminMembersUpdate)
	admin.Methods("GET").Path("/members/{id:[0-9]+}/notes").HandlerFunc(_h.AdminMembersNotes)
	admin.Methods("GET").Path("/notes/{id:[0-9]+}").HandlerFunc(_h.AdminNotes)
	admin.Methods("GET").Path("/organisations").HandlerFunc(_h.AdminOrganisations)
	admin.Methods("GET").Path("/organisations/{id:[0-9]+}/groups").HandlerFunc(_h.AdminOrganisationGroups)

	// these routes are available in the 'general' endpoints and are included here just for convenience
	admin.Methods("GET").Path("/resources/{id:[0-9]+}").HandlerFunc(_h.ResourcesID)
	admin.Methods("POST").Path("/resources").HandlerFunc(_h.ResourcesCollection)
	admin.Methods("GET").Path("/modules/{id:[0-9]+}").HandlerFunc(_h.ModulesID)
	admin.Methods("POST").Path("/modules").HandlerFunc(_h.ModulesCollection)

	// Batch routes for bulk uploading
	admin.Methods("POST").Path("/batch/resources").HandlerFunc(_h.AdminBatchResourcesPost)

	return admin
}

// adminMiddleWare wraps the require middleware handlers around the router passed in
func adminMiddleware(r *mux.Router) *negroni.Negroni {

	// Recovery from panic
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false // don't print the stack

	n := negroni.New()
	n.Use(recovery)
	n.Use(negroni.HandlerFunc(_mw.ValidateToken))
	n.Use(negroni.HandlerFunc(_mw.AdminScope))
	n.Use(negroni.NewLogger())
	n.Use(negroni.Wrap(r))

	return n
}
