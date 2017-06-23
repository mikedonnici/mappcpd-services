package router

import (
	"github.com/gorilla/mux"

	_h "github.com/mappcpd/web-services/cmd/webd/router/handlers"
)

// reportSubRouter sets up a router for report endpoints - no middleware for now
func reportSubRouter() *mux.Router {

	// TODO - add a KEY for reports
	r := mux.NewRouter().StrictSlash(true)
	reports := r.PathPrefix(v1ReportBase).Subrouter()
	reports.Methods("GET").Path("/test").HandlerFunc(_h.ReportsTest)
	reports.Methods("GET").Path("/modulesbydate").HandlerFunc(_h.ReportsModulesByDate)
	reports.Methods("GET").Path("/pointsbyrecorddate").HandlerFunc(_h.ReportsPointsByRecordDate)
	reports.Methods("GET").Path("/pointsbyactivitydate").HandlerFunc(_h.ReportsPointsByActivityDate)

	return reports
}
