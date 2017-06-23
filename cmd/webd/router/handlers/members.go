package handlers

import (
	"database/sql"
	"net/http"

	_json "github.com/mappcpd/web-services/cmd/webd/router/handlers/json"
	_mw "github.com/mappcpd/web-services/cmd/webd/router/handlers/middleware"
	m_ "github.com/mappcpd/web-services/internal/members"
	ds_ "github.com/mappcpd/web-services/internal/platform/datastore"
)

// MembersIDHandler fetches a member record by id
func MembersProfile(w http.ResponseWriter, r *http.Request) {

	p := _json.NewPayload()

	// Get user id from token
	id := _mw.UserAuthToken.Claims.ID

	// Get the Member record
	m, err := m_.MemberByID(id)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = _json.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + ds_.MySQL.Source}
		p.Data = m

		// TODO: remove this when fetching - should only be on update
		m_.SyncMember(m)
	}

	p.Send(w)
}

// MembersActivities fetches activity records for a member
func MembersActivities(w http.ResponseWriter, r *http.Request) {

	p := _json.NewPayload()

	a, err := m_.MemberActivitiesByMemberID(_mw.UserAuthToken.Claims.ID)

	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = _json.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + ds_.MySQL.Source}
	p.Meta = map[string]int{"count": len(a)}
	p.Data = a
	p.Send(w)
}

// MembersEvaluation created reports for each evaluation period
// by gathering the CPD activities within the dates, adding them up, applying caps etc
func MembersEvaluation(w http.ResponseWriter, r *http.Request) {

	p := _json.NewPayload()

	// Collect the evaluation periods
	es, err := m_.EvaluationsByMemberID(_mw.UserAuthToken.Claims.ID)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = _json.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + ds_.MySQL.Source}
	p.Meta = map[string]int{"count": len(es)}
	p.Data = es
	p.Send(w)
}
