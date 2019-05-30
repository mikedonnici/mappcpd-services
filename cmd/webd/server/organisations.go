package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cardiacsociety/web-services/internal/organisation"
	"github.com/gorilla/mux"
)

// AllOrganisations handles requests for Organisation records
func AllOrganisations(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	l, err := organisation.All(DS)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Data = l
	m := make(map[string]interface{})
	m["count"] = len(l)
	m["description"] = "All active organisations"
	p.Meta = m
	p.Send(w)
}

// OrganisationByID handles requests for a single Organisation record
func OrganisationByID(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failed", err.Error()}
	}

	o, err := organisation.ByID(DS, id)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Data = o
	m := make(map[string]interface{})
	m["description"] = fmt.Sprintf("Retrieved Organisation id %v - %s ", id, o.Name)
	p.Meta = m
	p.Send(w)
}
