package handlers

import (
	"strconv"

	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mappcpd/web-services/cmd/webd/rest/router/handlers/responder"
	"github.com/mappcpd/web-services/cmd/webd/rest/router/middleware"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/resources"
)

// ResourcesID fetches a single resource from the MySQLConnection db
func ResourcesID(w http.ResponseWriter, req *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)
	// Request - convert id from string to int type
	v := mux.Vars(req)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	r, err := resources.ResourceByID(id)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
		p.Data = r
		// Sync from MySQLConnection -> MongoDB
		resources.SyncResource(r)
	}

	p.Send(w)
}

// ResourcesCollection searches the Resources collection with search criteria POST'd as JSON request body
func ResourcesCollection(w http.ResponseWriter, r *http.Request) {

	// Response
	p := responder.New(middleware.UserAuthToken.Token)

	// Pull the JSON body out of the request
	decoder := json.NewDecoder(r.Body)
	var q datastore.MongoQuery
	err := decoder.Decode(&q)
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failure", errMessageDecodeJSON}
		p.Send(w)
		return
	}

	var res []interface{}
	res, err = resources.QueryResourcesCollection(q)
	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MongoDB.Source}
	c := len(res)
	p.Meta = responder.MongoMeta{c, q}
	p.Data = res
	p.Send(w)
}

// ResourcesLatest returns the most recent 'n' resources by createdAt date
func ResourcesLatest(w http.ResponseWriter, r *http.Request) {

	// Response
	p := responder.Payload{}

	// Request - convert id from string to int type
	v := mux.Vars(r)
	n, err := strconv.Atoi(v["n"])
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Grab the latest...
	var q datastore.MongoQuery
	q.Limit = n
	q.Sort = "-createdAt"

	var res []interface{}
	res, err = resources.QueryResourcesCollection(q)
	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MongoDB.Source}
	c := len(res)
	p.Meta = responder.MongoMeta{c, q}
	p.Data = res
	p.Send(w)
}
