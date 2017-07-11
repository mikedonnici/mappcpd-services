package handlers

import (
	"fmt"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"

	_json "github.com/mappcpd/web-services/cmd/webd/router/handlers/responder"
	_mw "github.com/mappcpd/web-services/cmd/webd/router/middleware"
	o_ "github.com/mappcpd/web-services/internal/organisations"
	ds_ "github.com/mappcpd/web-services/internal/platform/datastore"
)

// AdminOrganisations handles requests for Organisation records
// Todo there is only one organisation as such so this is retrieving groups.
func AdminOrganisations(w http.ResponseWriter, r *http.Request) {

	p := _json.New(_mw.UserAuthToken.Token)

	l, err := o_.OrganisationsList()
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + ds_.MySQL.Source}
	p.Data = l
	m := make(map[string]interface{})
	m["count"] = len(l)
	m["description"] = "List of organisations..."
	p.Meta = m
	p.Send(w)
}

// AdminOrganisationGroups handles requests for Organisation records
func AdminOrganisationGroups(w http.ResponseWriter, r *http.Request) {

	p := _json.New(_mw.UserAuthToken.Token)

	// Request - convert Organisation id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = _json.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Get groups... todo - we should do this via GetOrganisationByID
	l, err := o_.OrganisationGroupsList(id)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// Grab the organisation name for our message
	o, err := o_.OrganisationByID(id)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + ds_.MySQL.Source}
	p.Data = l
	m := make(map[string]interface{})
	m["count"] = len(l)
	m["description"] = fmt.Sprintf("Retrieved groups for Organisation id %v - %s ", id, o.Name)
	p.Meta = m
	p.Send(w)
}
