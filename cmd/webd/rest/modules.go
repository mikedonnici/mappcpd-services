package rest

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mappcpd/web-services/internal/modules"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// ModulesID fetches a single resource from the MySQLConnection db
func ModulesID(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Token)
	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failed", err.Error()}
	}

	m, err := modules.ModuleByID(DS, id)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
		p.Data = m
		modules.SyncModule(DS, m)
	}

	p.Send(w)
}

// ModulesCollection searches the Modules collection with search criteria POST'd as JSON request body
func ModulesCollection(w http.ResponseWriter, r *http.Request) {

	// Response
	p := NewResponder(UserAuthToken.Token)

	// Pull the JSON body out of the request
	decoder := json.NewDecoder(r.Body)
	var q datastore.MongoQuery
	err := decoder.Decode(&q)
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failure", errMessageDecodeJSON}
		p.Send(w)
		return
	}

	var res []interface{}
	res, err = modules.QueryModulesCollection(DS, q)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	c := len(res)
	p.Meta = MongoMeta{c, q}
	p.Data = res
	p.Send(w)
}