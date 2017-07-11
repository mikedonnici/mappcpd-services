package handlers

import (
	"fmt"
	"strconv"

	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mappcpd/web-services/cmd/webd/router/handlers/responder"
	"github.com/mappcpd/web-services/cmd/webd/router/middleware"
	"github.com/mappcpd/web-services/internal/attachments"
	"github.com/mappcpd/web-services/internal/generic"
	"github.com/mappcpd/web-services/internal/members"
	"github.com/mappcpd/web-services/internal/notes"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/resources"
)

func AdminTest(w http.ResponseWriter, _ *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)
	p.Message = responder.Message{http.StatusOK, "success", "Hi Admin!"}
	p.Send(w)
}

// AdminMembersSearch searches the document database based on GET request
// Taking advantage of the complexity of Mongo queries and pipes is
// difficult using URI parameters. So this is an attempt however will also
// implement a POST version below to allow for a complete JSON query doc
// // to be submitted. Being totally RESTful is not as important  as this
// API is for DB access at this stage.
func AdminMembersSearch(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	var err error
	var query map[string]interface{}
	var projection map[string]interface{}

	// Query
	query, err = queryParams(r.FormValue("q"))
	if err != nil {
		if err != nil {
			p.Message = responder.Message{
				http.StatusBadRequest,
				"failed",
				err.Error(),
			}
			p.Send(w)
			return
		}
	}

	// Projection
	projection = projectParams(r.FormValue("p"))

	// limit
	limit := 0
	if len(r.FormValue("l")) > 0 {
		limit, err = strconv.Atoi(r.FormValue("l"))
		if err != nil {
			p.Message = responder.Message{
				http.StatusBadRequest,
				"failed",
				err.Error(),
			}
			p.Send(w)
			return
		}
	}

	// Run the query...
	var res []interface{}
	if limit > 0 {
		res, err = members.DocMembersLimit(query, projection, limit)
	} else {
		res, err = members.DocMembersAll(query, projection)
	}

	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MongoDB.Source}
	c := len(res)
	p.Meta = responder.DocMeta{c, query, projection}
	p.Data = res
	p.Send(w)
}

// AdminMembersSearchPost uses POST body to specify the search criteria. May not be ReSTful
// but is easier to pass a query as JSON doc in body. Could (at some stage) store the
// POSTed query and return a URL to fetch it. This way it follows ReSTful principles
// and the query can be kept for later / cached?
func AdminMembersSearchPost(w http.ResponseWriter, r *http.Request) {

	// create a binding struct for the JSON request body
	// ie. this is what we are expecting
	type Find struct {
		Query      map[string]interface{} `json:"query"`
		Projection map[string]interface{} `json:"projection"`
		Limit      int                    `json:"limit"`
	}

	p := responder.New(middleware.UserAuthToken.Token)

	// Pull the JSON body out of the request
	decoder := json.NewDecoder(r.Body)
	var f Find
	err := decoder.Decode(&f)
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failure", errMessageDecodeJSON}
		p.Send(w)
		return
	}

	// Limit query
	var res []interface{}
	if f.Limit > 0 {
		res, err = members.DocMembersLimit(f.Query, f.Projection, f.Limit)
	} else {
		res, err = members.DocMembersAll(f.Query, f.Projection)
	}
	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MongoDB.Source}
	c := len(res)
	p.Meta = responder.DocMeta{c, f.Query, f.Projection}
	p.Data = res
	p.Send(w)
}

// AdminMembersUpdate will update a member record by processing the JSON body
func AdminMembersUpdate(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Response
	m, err := members.MemberByID(id)

	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
	default:

		// Pull the JSON body out of the request
		decoder := json.NewDecoder(r.Body)
		var j map[string]interface{}
		err = decoder.Decode(&j)
		if err != nil {
			p.Message = responder.Message{http.StatusBadRequest, "failure", err.Error()}
			p.Send(w)
			return
		}

		fmt.Printf("%T %s\n", j, j)
		fmt.Printf("j[id] %T %v\n", j["id"], j["id"])
		fmt.Printf("m[id] %T %v\n", m.ID, m.ID)

		// As a small sanity check make sure the id on the url
		// matches the id passed in the JSON body
		if j["id"] == "" {
			p.Message = responder.Message{http.StatusBadRequest, "failed", "MySQLConnection row id must be included in the JSON body"}
			p.Send(w)
			return
		}
		// need type assertion as j["id"] is float64 when decoded from JSON
		jid := int(j["id"].(float64))
		fmt.Printf("%v %T - %v %T", m.ID, m.ID, jid, jid)
		if m.ID != jid {
			p.Message = responder.Message{http.StatusBadRequest, "failed", "ID on the request URL does not match the ID in the Body"}
			p.Send(w)
			return
		}

		err := members.UpdateMember(j)
		if err != nil {
			p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
			p.Send(w)
			return
		}

		m, _ = members.MemberByID(id) // Re-fetch
		members.SyncMember(m)         // Sync to doc db

		p.Message = responder.Message{http.StatusOK, "success", "MySQLConnection record updated and copied to MongoDB"}
		p.Data = m
	}

	p.Send(w)
}

// AdminMembersNotes fetches all Notes belonging to a Member
func AdminMembersNotes(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Response
	ns, err := notes.NotesByMemberID(id)
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
		p.Data = ns
	}

	p.Send(w)
}

// AdminNotes fetches a single Note record by Note ID
func AdminNotes(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Response
	d, err := notes.NoteById(id)
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
		p.Data = d
	}

	p.Send(w)
}

// AdminMembersID fetches a member record from the MySQLConnection DB, by id
func AdminMembersID(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Get the Member record
	m, err := members.MemberByID(id)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
		p.Data = m
		members.SyncMember(m)
	}

	p.Send(w)
}

// AdminMembersIDListHandler fetches a list of all member ids from MySQL
func AdminIDList(w http.ResponseWriter, req *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Request - requires at least the 't' query to specify the table name
	// and can have the option 'f' as raw HTML filter
	t := req.FormValue("t")
	if t == "" {
		p.Message = responder.Message{http.StatusBadRequest, "failed", "Requires ?t=[table_name], optional &f=[sql_filter]"}
		p.Send(w)
		return
	}

	// Optional filter
	f := req.FormValue("f")

	// Get the Member record
	ii, err := generic.GetIDs(t, f)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = responder.Message{http.StatusOK, "success", "List of ids from table: " + t + ", db: " + datastore.MySQL.Source}
		p.Meta = map[string]int{"count": len(ii)}
		p.Data = ii
	}

	p.Send(w)
}

// AdminBatchResourcesPost will upload a set of resource records to MySQL
func AdminBatchResourcesPost(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Handling batch resources upload...")

	// Expecting a JSON body with a single 'data' field containing array of Resources to be inserted
	type batch struct {
		Data resources.Resources `json:"data"`
	}
	b := batch{}

	p := responder.New(middleware.UserAuthToken.Token)

	// Pull the JSON body out of the request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failure", errMessageDecodeJSON}
		p.Send(w)
		return
	}
	defer r.Body.Close()

	// In our return data we can store the results for each attempt. Even though .Save() may
	// return an error it could be something minor such as a missing url, which is no
	// reason to stop processing the batch of resources
	var data = struct {
		Failures   map[string]string
		SuccessIDs []int64
	}{}
	data.Failures = make(map[string]string)

	// Store any problems records in meta
	failCount := 0
	successCount := 0

	// Range over .Data
	for _, v := range b.Data {
		r := resources.Resource{}
		r = v
		id, err := r.Save()
		if err != nil {
			data.Failures[r.Name] = err.Error()
			failCount += 1
			continue
		}
		// otherwise, add the id to the list
		data.SuccessIDs = append(data.SuccessIDs, id)
		successCount += 1
	}

	p.Message = responder.Message{http.StatusOK, "success", "Batch completed - see Failures for errors"}
	p.Meta = map[string]int{
		"failed":     failCount,
		"successful": successCount,
	}
	p.Data = data
	p.Send(w)
}

// AdminAttachmentAdd registers a file attachment for a database entity specified in the request body
func AdminAttachmentAdd(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Create an attachment
	var a attachments.Attachment

	// Set the admin id from the token
	a.UserID = middleware.UserAuthToken.Claims.ID

	// Decode post body into attachment fields
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		msg := "Could not decode json in request body - " + err.Error()
		p.Message = responder.Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}
	fmt.Println(a)

	// Check if attachment already exists before registering
	id, err := a.Exists()
	if err != nil {
		msg := "Error checking for duplicate attachment - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	if id > 0 {
		msg := "An attachment with the same file name is already registered. The sanitized file name is %s and the " +
			"attachment record is %s.id = %d. If this is NOT a duplicate file try changing the original file name."
		msg = fmt.Sprintf(msg, a.CleanFilename, a.EntityName, id)
		p.Message = responder.Message{http.StatusConflict, "failed", msg}
		p.Send(w)
		return
	}

	// Register the attachment
	if err := a.Register(); err != nil {
		msg := "Error registering attachment - " + err.Error()
		p.Message = responder.Message{http.StatusForbidden, "failed", msg}
		p.Send(w)
		return
	}

	p.Message = responder.Message{http.StatusOK, "success", "Attachment registered"}
	p.Data = a
	p.Send(w)

}
