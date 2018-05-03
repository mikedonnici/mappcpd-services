package handlers

import (
	"fmt"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"

	"github.com/mappcpd/web-services/cmd/webd/rest/router/handlers/responder"
	"github.com/mappcpd/web-services/cmd/webd/rest/router/middleware"
	"github.com/mappcpd/web-services/internal/organisation"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// AllOrganisations handles requests for Organisation records
func AllOrganisations(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	l, err := organisation.All()
	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Description}
	p.Data = l
	m := make(map[string]interface{})
	m["count"] = len(l)
	m["description"] = "All active organisations"
	p.Meta = m
	p.Send(w)
}

// OrganisationByID handles requests for a single Organisation record
func OrganisationByID(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	o, err := organisation.ByID(id)
	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Description}
	p.Data = o
	m := make(map[string]interface{})
	m["description"] = fmt.Sprintf("Retrieved Organisation id %v - %s ", id, o.Name)
	p.Meta = m
	p.Send(w)
}
