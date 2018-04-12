package handlers

import (
	"database/sql"
	"net/http"

	"github.com/mappcpd/web-services/cmd/webd/rest/router/handlers/responder"
	"github.com/mappcpd/web-services/cmd/webd/rest/router/middleware"
	"github.com/mappcpd/web-services/internal/member"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/member/activity"
)

// MembersProfile fetches a member record by id
func MembersProfile(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Get user id from token
	id := middleware.UserAuthToken.Claims.ID

	// Get the Member record
	m, err := member.MemberByID(id)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
		p.Data = m

		// TODO: remove this when fetching - should only be on update
		member.SyncMember(m)
	}

	p.Send(w)
}

// MembersActivities fetches activity records for a member
func MembersActivities(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	a, err := activity.MemberActivitiesByMemberID(middleware.UserAuthToken.Claims.ID)

	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
	p.Meta = map[string]int{"count": len(a)}
	p.Data = a
	p.Send(w)
}

// MembersEvaluation created reports for each evaluation period
// by gathering the CPD activities within the dates, adding them up, applying caps etc
func MembersEvaluation(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Collect the evaluation periods
	es, err := activity.MemberActivityReports(middleware.UserAuthToken.Claims.ID)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
	p.Meta = map[string]int{"count": len(es)}
	p.Data = es
	p.Send(w)
}

// MembersReports
func MembersReports(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	ce, err := activity.CurrentMemberActivityReport(middleware.UserAuthToken.Claims.ID)

	// todo: testing trigger
	activity.PDFReport(ce)

	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
	p.Data = ce
	p.Send(w)
}
