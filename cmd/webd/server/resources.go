package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/resource"
	"github.com/gorilla/mux"
)

// ResourcesID fetches a single resource from the MySQLConnection db
func ResourcesID(w http.ResponseWriter, req *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)
	// Request - convert id from string to int type
	v := mux.Vars(req)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failed", err.Error()}
	}

	r, err := resource.ByID(DS, id)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
		p.Data = r
		// Sync from MySQLConnection -> MongoDB - runs ina  separate go routine
		resource.SyncResource(DS, r)
	}

	p.Send(w)
}

// ResourcesCollection searches the Resources collection with search criteria POST'd as JSON request body
func ResourcesCollection(w http.ResponseWriter, r *http.Request) {

	// Response
	p := NewResponder(UserAuthToken.Encoded)

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
	res, err = resource.QueryResourcesCollection(DS, q)
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

// ResourcesLatest returns the most recent 'n' resources by createdAt date
func ResourcesLatest(w http.ResponseWriter, r *http.Request) {

	// Response
	p := Payload{}

	// Request - convert id from string to int type
	v := mux.Vars(r)
	n, err := strconv.Atoi(v["n"])
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Grab the latest...
	var q datastore.MongoQuery
	q.Limit = n
	q.Sort = "-createdAt"

	var res []interface{}
	res, err = resource.QueryResourcesCollection(DS, q)
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
