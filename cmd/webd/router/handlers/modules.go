package handlers

import (
	"strconv"

	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	_json "github.com/mappcpd/web-services/cmd/webd/router/handlers/json"
	m_ "github.com/mappcpd/web-services/internal/modules"
	ds_ "github.com/mappcpd/web-services/internal/platform/datastore"
)

// ModulesID fetches a single resource from the MySQLConnection db
func ModulesID(w http.ResponseWriter, req *http.Request) {

	p := _json.Payload{}
	// Request - convert id from string to int type
	v := mux.Vars(req)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = _json.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	m, err := m_.ModuleByID(id)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = _json.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + ds_.MySQL.Source}
		p.Data = m
		m_.SyncModule(m)
	}

	p.Send(w)
}

// ModulesCollection searches the Modules collection with search criteria POST'd as JSON request body
func ModulesCollection(w http.ResponseWriter, r *http.Request) {

	// Response
	p := _json.Payload{}

	// Pull the JSON body out of the request
	decoder := json.NewDecoder(r.Body)
	var q ds_.MongoQuery
	err := decoder.Decode(&q)
	if err != nil {
		p.Message = _json.Message{http.StatusBadRequest, "failure", errMessageDecodeJSON}
		p.Send(w)
		return
	}

	var res []interface{}
	res, err = m_.QueryModulesCollection(q)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + ds_.MongoDB.Source}
	c := len(res)
	p.Meta = _json.MongoMeta{c, q}
	p.Data = res
	p.Send(w)
}
