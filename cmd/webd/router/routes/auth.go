package routes

import (
	"github.com/gorilla/mux"

	"github.com/mappcpd/web-services/cmd/webd/router/handlers"
)

// AuthSubRouter sets up a router for auth with no middleware
func AuthSubRouter(prefix string) *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	auth := r.PathPrefix(prefix).Subrouter()
	auth.Methods("POST").Path("/member").HandlerFunc(handlers.AuthMemberLogin)
	auth.Methods("POST").Path("/admin").HandlerFunc(handlers.AuthAdminLogin)

	return auth
}
