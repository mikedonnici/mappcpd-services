package server

import (
	"fmt"
	"net/http"

	"github.com/cardiacsociety/web-services/internal/organisation"

	"github.com/cardiacsociety/web-services/internal/qualification"
	"github.com/cardiacsociety/web-services/internal/speciality"
	"github.com/gorilla/mux"
)

// Qualifications fetches list of Qualifications
func Qualifications(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	xq, err := qualification.All(DS)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Data retrieved from " + DS.MySQL.Desc}
	p.Data = xq
	m := make(map[string]interface{})
	m["count"] = len(xq)
	m["description"] = "List of Qualifications"
	p.Meta = m
	p.Send(w)
}

// Specialities fetches list of Specialities (areas of interest)
func Specialities(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	xq, err := speciality.All(DS)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Data retrieved from " + DS.MySQL.Desc}
	p.Data = xq
	m := make(map[string]interface{})
	m["count"] = len(xq)
	m["description"] = "List of Specialities"
	p.Meta = m
	p.Send(w)
}

// Organisations fetches list of Organisations and can include a typeId on the url.
func Organisations(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	v := mux.Vars(r)
	// endpoint .../organisations/ with no type returns 404, so this will never run
	if v["type"] == "" {
		p.Message = Message{http.StatusBadRequest, "failed", " organisation type not specified"}
		p.Send(w)
		return
	}

	var typeID int
	switch v["type"] {
	case "councils":
		typeID = 1
	case "groups", "group", "workinggroups", "workinggroup":
		typeID = 2
	case "committees", "committee":
		typeID = 3
	case "education", "educational", "university", "universities", "tertiary":
		typeID = 8
	case "board":
		typeID = 5
	case "other":
		typeID = 6
	case "institutes", "institute", "hospitals", "hospital":
		typeID = 7
	case "foundations", "foundation":
		typeID = 9
	case "government":
		typeID = 10
	}

	xo, err := organisation.ByTypeID(DS, typeID)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Data retrieved from " + DS.MySQL.Desc}
	p.Data = xo
	m := make(map[string]interface{})
	m["count"] = len(xo)
	m["description"] = fmt.Sprintf("List of Organisations of type: '%s'", v["type"])
	p.Meta = m
	p.Send(w)
}
