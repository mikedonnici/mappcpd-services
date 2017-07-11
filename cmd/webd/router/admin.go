package router

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/mappcpd/web-services/cmd/webd/router/handlers"
	"github.com/mappcpd/web-services/cmd/webd/router/middleware"
)

// adminSubRouter adds end points for admin, and appropriate middleware
func adminSubRouter() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	admin := r.PathPrefix(v1AdminBase).Subrouter()

	admin.Methods("GET").Path("/test").HandlerFunc(handlers.AdminTest)
	admin.Methods("GET").Path("/idlist").HandlerFunc(handlers.AdminIDList)
	admin.Methods("GET").Path("/members").HandlerFunc(handlers.AdminMembersSearch)
	admin.Methods("POST").Path("/members").HandlerFunc(handlers.AdminMembersSearchPost)
	admin.Methods("GET").Path("/members/{id:[0-9]+}").HandlerFunc(handlers.AdminMembersID)
	admin.Methods("POST").Path("/members/{id:[0-9]+}").HandlerFunc(handlers.AdminMembersUpdate)
	admin.Methods("GET").Path("/members/{id:[0-9]+}/notes").HandlerFunc(handlers.AdminMembersNotes)
	admin.Methods("GET").Path("/notes/{id:[0-9]+}").HandlerFunc(handlers.AdminNotes)
	admin.Methods("GET").Path("/organisations").HandlerFunc(handlers.AdminOrganisations)
	admin.Methods("GET").Path("/organisations/{id:[0-9]+}/groups").HandlerFunc(handlers.AdminOrganisationGroups)

	// these routes are available in the 'general' endpoints and are included here just for convenience
	admin.Methods("GET").Path("/resources/{id:[0-9]+}").HandlerFunc(handlers.ResourcesID)
	admin.Methods("POST").Path("/resources").HandlerFunc(handlers.ResourcesCollection)
	admin.Methods("GET").Path("/modules/{id:[0-9]+}").HandlerFunc(handlers.ModulesID)
	admin.Methods("POST").Path("/modules").HandlerFunc(handlers.ModulesCollection)

	// Attachment registration
	admin.Methods("POST").Path("/attachments").HandlerFunc(handlers.AdminAttachmentAdd)

	// Batch routes for bulk uploading
	admin.Methods("POST").Path("/batch/resources").HandlerFunc(handlers.AdminBatchResourcesPost)

	return admin
}

// adminMiddleWare wraps the require middleware handlers around the router passed in
func adminMiddleware(r *mux.Router) *negroni.Negroni {

	// Recovery from panic
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false // don't print the stack

	n := negroni.New()
	n.Use(recovery)
	n.Use(negroni.HandlerFunc(middleware.ValidateToken))
	n.Use(negroni.HandlerFunc(middleware.AdminScope))
	n.Use(negroni.NewLogger())
	n.Use(negroni.Wrap(r))

	return n
}
