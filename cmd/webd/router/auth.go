package router

import (
	"github.com/gorilla/mux"

	_h "github.com/mappcpd/web-services/cmd/webd/router/handlers"
)

// authSubRouter sets up a router for auth with no middleware
func authSubRouter() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	auth := r.PathPrefix(v1AuthBase).Subrouter()
	auth.Methods("POST").Path("/member").HandlerFunc(_h.AuthMemberLogin)
	auth.Methods("POST").Path("/admin").HandlerFunc(_h.AuthAdminLogin)

	return auth
}
