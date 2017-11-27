package routes

import (
	"github.com/gorilla/mux"

	"github.com/mappcpd/web-services/cmd/webd/rest/router/handlers"
)

// ReportSubRouter sets up a router for report endpoints - no middleware for now
func ReportSubRouter(prefix string) *mux.Router {

	// TODO - add a KEY for reports
	r := mux.NewRouter().StrictSlash(true)
	reports := r.PathPrefix(prefix).Subrouter()
	reports.Methods("GET").Path("/test").HandlerFunc(handlers.ReportsTest)
	reports.Methods("GET").Path("/modulesbydate").HandlerFunc(handlers.ReportsModulesByDate)
	reports.Methods("GET").Path("/pointsbyrecorddate").HandlerFunc(handlers.ReportsPointsByRecordDate)
	reports.Methods("GET").Path("/pointsbyactivitydate").HandlerFunc(handlers.ReportsPointsByActivityDate)

	return reports
}
