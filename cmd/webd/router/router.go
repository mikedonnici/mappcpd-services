package router

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	_h "github.com/mappcpd/web-services/cmd/webd/router/handlers"
)

const (
	v1AuthBase    = "/v1/auth"
	v1MemberBase  = "/v1/m"
	v1AdminBase   = "/v1/a"
	v1GeneralBase = "/v1/g"
	v1ReportBase  = "/v1/r"
)

func Start() {

	// Router
	r := mux.NewRouter()

	r.Methods("OPTIONS").HandlerFunc(_h.Preflight)

	// Ping and preflight, no middleware required
	r.Methods("GET").Path("/").HandlerFunc(_h.Index)
	//r.Methods("OPTIONS").Path("/").HandlerFunc(_h.Preflight)

	// Auth sub-router, no middleware required
	rAuth := authSubRouter()
	r.PathPrefix(v1AuthBase).Handler(rAuth)

	// Admin sub-router and middleware
	rAdmin := adminSubRouter()                          // add router...
	rAdminMiddleware := adminMiddleware(rAdmin)         // ...plus middleware...
	r.PathPrefix(v1AdminBase).Handler(rAdminMiddleware) // ...and add to main router

	// Reports sub-router, todo: add middleware to reports router
	rReports := reportSubRouter()
	r.PathPrefix(v1ReportBase).Handler(rReports)

	// Member sub-router
	rMember := memberSubRouter()
	rMemberMiddleware := memberMiddleware(rMember)
	r.PathPrefix(v1MemberBase).Handler(rMemberMiddleware)

	// General sub-router
	rGeneral := generalSubRouter()
	rGeneralMiddleware := generalMiddleware(rGeneral)
	r.PathPrefix(v1GeneralBase).Handler(rGeneralMiddleware)

	// Specify port when env var is not set - Heroku sets dynamically so cannot include in .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// CORS handler - needed to add OptionsPassThrough for preflight requests which use OPTIONS http method
	//handler := cors.Default().Handler(r)
	// Todo... tighten this up - not sure if needed  with preflightHandler??
	handler := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"*"},
		AllowedHeaders:     []string{"*"},
		OptionsPassthrough: true,
	}).Handler(r)

	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
