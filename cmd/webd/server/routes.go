package server

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// AuthSubRouter sets up a router for auth with no middleware
func AuthSubRouter(prefix string) *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	auth := r.PathPrefix(prefix).Subrouter()
	auth.Methods("OPTIONS").Path("/").HandlerFunc(Preflight)
	auth.Methods("POST").Path("/member").HandlerFunc(AuthMemberLogin)
	auth.Methods("POST").Path("/admin").HandlerFunc(AuthAdminLogin)

	return auth
}

// AdminSubRouter adds end points for admin, and appropriate middleware
func AdminSubRouter(prefix string) *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	admin := r.PathPrefix(prefix).Subrouter()

	admin.Methods("GET").Path("/test").HandlerFunc(AdminTest)
	admin.Methods("GET").Path("/idlist").HandlerFunc(AdminIDList)
	admin.Methods("GET").Path("/members").HandlerFunc(AdminMembersSearch)
	admin.Methods("POST").Path("/members").HandlerFunc(AdminMembersSearchPost)
	admin.Methods("GET").Path("/members/{id:[0-9]+}").HandlerFunc(AdminMembersID)
	//admin.Methods("POST").Path("/members/{id:[0-9]+}").HandlerFunc(AdminMembersUpdate)
	admin.Methods("GET").Path("/members/{id:[0-9]+}/notes").HandlerFunc(AdminMembersNotes)
	admin.Methods("GET").Path("/notes/{id:[0-9]+}").HandlerFunc(AdminNotes)
	admin.Methods("GET").Path("/organisations").HandlerFunc(AllOrganisations)
	admin.Methods("GET").Path("/organisations/{id:[0-9]+}").HandlerFunc(OrganisationByID)

	// these routes are available in the 'general' endpoints and are included here just for convenience
	admin.Methods("GET").Path("/resources/{id:[0-9]+}").HandlerFunc(ResourcesID)
	admin.Methods("POST").Path("/resources").HandlerFunc(ResourcesCollection)
	admin.Methods("GET").Path("/modules/{id:[0-9]+}").HandlerFunc(ModulesID)
	admin.Methods("POST").Path("/modules").HandlerFunc(ModulesCollection)

	// Note Attachments
	admin.Methods("OPTIONS").Path("/notes/{id:[0-9]+}/attachments/request").HandlerFunc(Preflight)
	admin.Methods("GET").Path("/notes/{id:[0-9]+}/attachments/request").HandlerFunc(AdminNotesAttachmentRequest)
	admin.Methods("PUT").Path("/notes/{id:[0-9]+}/attachments").HandlerFunc(AdminNotesAttachmentRegister)

	// Resource Attachments
	admin.Methods("OPTIONS").Path("/resources/{id:[0-9]+}/attachments/request").HandlerFunc(Preflight)
	admin.Methods("GET").Path("/resources/{id:[0-9]+}/attachments/request").HandlerFunc(AdminResourcesAttachmentRequest)
	admin.Methods("PUT").Path("/resources/{id:[0-9]+}/attachments").HandlerFunc(AdminResourcesAttachmentRegister)

	// Batch routes for bulk uploading
	admin.Methods("POST").Path("/batch/resources").HandlerFunc(AdminBatchResourcesPost)

	// Report routes
	admin.Methods("POST").Path("/reports/application").HandlerFunc(AdminReportApplicationExcel)
	admin.Methods("POST").Path("/reports/member").HandlerFunc(AdminReportMemberExcel)
	admin.Methods("POST").Path("/reports/journal").HandlerFunc(AdminReportMemberJournalExcel)
	admin.Methods("POST").Path("/reports/invoice").HandlerFunc(AdminReportInvoiceExcel)
	admin.Methods("POST").Path("/reports/payment").HandlerFunc(AdminReportPaymentExcel)
	admin.Methods("POST").Path("/reports/position").HandlerFunc(AdminReportPositionExcel)

	// Membership application
	admin.Methods("POST").Path("/applications").HandlerFunc(AdminNewMembershipApplication)
	
	// Lapse members
	admin.Methods("PUT").Path("/lapsedmembers").HandlerFunc(AdminLapseMembers)

	// Notifications
	admin.Methods("POST").Path("/notifications").HandlerFunc(AdminSendNotifications)

	return admin
}

// AdminMiddleware wraps the require middleware around the router passed in
func AdminMiddleware(r *mux.Router) *negroni.Negroni {

	// Recovery from panic
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false // don't print the stack

	n := negroni.New()
	n.Use(recovery)
	n.Use(negroni.HandlerFunc(ValidateToken))
	n.Use(negroni.HandlerFunc(AdminScope))
	n.Use(negroni.NewLogger())
	n.Use(negroni.Wrap(r))

	return n
}

// GeneralSubRouter is a sub router for requests relevant to all users
func GeneralSubRouter(prefix string) *mux.Router {

	// Middleware for General sub-router just need a valid token
	// as these are used by both admin and member scope
	r := mux.NewRouter().StrictSlash(true)

	// general routes
	general := r.PathPrefix(prefix).Subrouter()

	// Activity (types)
	general.Methods("GET").Path("/activities").HandlerFunc(Activities)
	general.Methods("GET").Path("/activities/{id:[0-9]+}").HandlerFunc(ActivitiesID)

	general.Methods("GET").Path("/qualifications").HandlerFunc(Qualifications)
	general.Methods("GET").Path("/specialities").HandlerFunc(Specialities)
	general.Methods("GET").Path("/organisations/{type}").HandlerFunc(Organisations)

	// Resources
	general.Methods("GET").Path("/resources/{id:[0-9]+}").HandlerFunc(ResourcesID)
	general.Methods("POST").Path("/resources").HandlerFunc(ResourcesCollection)
	general.Methods("GET").Path("/resources/latest/{n:[0-9]+}").HandlerFunc(ResourcesLatest)

	// Modules
	general.Methods("GET").Path("/modules/{id:[0-9]+}").HandlerFunc(ModulesID)
	general.Methods("POST").Path("/modules").HandlerFunc(ModulesCollection)

	return general
}

// GeneralMiddleware applies required middleware to 'general' endpoints
func GeneralMiddleware(r *mux.Router) *negroni.Negroni {

	// Recovery from panic
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false // don't print the stack

	n := negroni.New()
	n.Use(recovery)
	n.Use(negroni.HandlerFunc(ValidateToken))
	n.Use(negroni.NewLogger())
	n.Use(negroni.Wrap(r))

	return n
}

// MemberSubRouter is a sub router for endpoints relevant to member user requests.
func MemberSubRouter(prefix string) *mux.Router {

	// Middleware for Members sub-router
	r := mux.NewRouter().StrictSlash(true)

	// members routes
	members := r.PathPrefix(prefix).Subrouter()
	members.Methods("GET").Path("/").HandlerFunc(Index)
	members.Methods("GET").Path("/token").HandlerFunc(MembersToken)
	members.Methods("OPTIONS").Path("/token").HandlerFunc(Preflight)
	members.Methods("GET").Path("/profile").HandlerFunc(MembersProfile)

	members.Methods("GET").Path("/activities").HandlerFunc(MembersActivities)
	members.Methods("POST").Path("/activities").HandlerFunc(MembersActivitiesAdd)

	members.Methods("GET").Path("/activities/{id:[0-9]+}").HandlerFunc(MembersActivitiesID)
	members.Methods("PUT").Path("/activities/{id:[0-9]+}").HandlerFunc(MembersActivitiesUpdate)

	// Attachments
	members.Methods("OPTIONS").Path("/activities/{id:[0-9]+}/attachments/request").HandlerFunc(Preflight)
	members.Methods("GET").Path("/activities/{id:[0-9]+}/attachments/request").HandlerFunc(MembersActivitiesAttachmentRequest)
	// This is idempotent, hence PUT
	members.Methods("PUT").Path("/activities/{id:[0-9]+}/attachments").HandlerFunc(MembersActivitiesAttachmentRegister)

	members.Methods("GET").Path("/activities/recurring").HandlerFunc(MembersActivitiesRecurring)
	members.Methods("POST").Path("/activities/recurring").HandlerFunc(MembersActivitiesRecurringAdd)

	members.Methods("OPTIONS").Path("/activities/recurring/{_id}").HandlerFunc(Preflight)
	members.Methods("DELETE").Path("/activities/recurring/{_id}").HandlerFunc(MembersActivitiesRecurringRemove)

	members.Methods("OPTIONS").Path("/activities/recurring/{_id}/recorder").HandlerFunc(Preflight)
	members.Methods("POST").Path("/activities/recurring/{_id}/recorder").HandlerFunc(MembersActivitiesRecurringRecorder)

	members.Methods("GET").Path("/evaluations").HandlerFunc(MembersEvaluation)

	members.Methods("POST").Path("/notifications").HandlerFunc(MemberSendNotification)

	members.Methods("GET").Path("/reports/cpd/current").HandlerFunc(CurrentActivityReport)
	members.Methods("GET").Path("/reports/cpd/current/emailer").HandlerFunc(EmailCurrentActivityReport)
	members.Methods("GET").Path("/reports//current/responder").HandlerFunc(EmailCurrentActivityReport)

	return members
}

// MemberMiddleware wraps the member sub router with appropriate middleware
func MemberMiddleware(r *mux.Router) *negroni.Negroni {

	// Recovery from panic
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false // don't print the stack

	n := negroni.New()
	n.Use(recovery)
	n.Use(negroni.HandlerFunc(ValidateToken))
	n.Use(negroni.HandlerFunc(MemberScope))
	n.Use(negroni.NewLogger())
	n.Use(negroni.Wrap(r))

	return n
}

// ReportSubRouter sets up a router for report endpoints - no middleware for now
func ReportSubRouter(prefix string) *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	reports := r.PathPrefix(prefix).Subrouter()
	reports.Methods("GET").Path("/test").HandlerFunc(ReportsTest)
	reports.Methods("GET").Path("/modulesbydate").HandlerFunc(ReportsModulesByDate)
	reports.Methods("GET").Path("/pointsbyrecorddate").HandlerFunc(ReportsPointsByRecordDate)
	reports.Methods("GET").Path("/pointsbyactivitydate").HandlerFunc(ReportsPointsByActivityDate)
	reports.Methods("GET").Path("/excel/{id}").HandlerFunc(ReportsExcel)

	return reports
}
