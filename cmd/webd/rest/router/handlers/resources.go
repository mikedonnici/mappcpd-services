package handlers

import (
	"strconv"

	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	_json "github.com/mappcpd/web-services/cmd/webd/rest/router/handlers/responder"
	mw_ "github.com/mappcpd/web-services/cmd/webd/rest/router/middleware"
	ds_ "github.com/mappcpd/web-services/internal/platform/datastore"
	r_ "github.com/mappcpd/web-services/internal/resources"
)

// ResourcesID fetches a single resource from the MySQLConnection db
func ResourcesID(w http.ResponseWriter, req *http.Request) {

	p := _json.New(mw_.UserAuthToken.Token)
	// Request - convert id from string to int type
	v := mux.Vars(req)
	id, err := strconv.ParseInt(v["id"], 10, 0)
	if err != nil {
		p.Message = _json.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	r, err := r_.ResourceByID(id)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = _json.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + ds_.MySQL.Source}
		p.Data = r
		// Sync from MySQLConnection -> MongoDB
		r_.SyncResource(r)
	}

	p.Send(w)
}

// ResourcesCollection searches the Resources collection with search criteria POST'd as JSON request body
func ResourcesCollection(w http.ResponseWriter, r *http.Request) {

	// Response
	p := _json.New(mw_.UserAuthToken.Token)

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
	res, err = r_.QueryResourcesCollection(q)
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

// ResourcesLatest returns the most recent 'n' resources by createdAt date
func ResourcesLatest(w http.ResponseWriter, r *http.Request) {

	// Response
	p := _json.Payload{}

	// Request - convert id from string to int type
	v := mux.Vars(r)
	n, err := strconv.Atoi(v["n"])
	if err != nil {
		p.Message = _json.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Grab the latest...
	var q ds_.MongoQuery
	q.Limit = n
	q.Sort = "-createdAt"

	var res []interface{}
	res, err = r_.QueryResourcesCollection(q)
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
