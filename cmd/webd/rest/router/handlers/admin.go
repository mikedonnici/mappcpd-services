package handlers

import (
	"fmt"
	"strconv"

	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mappcpd/web-services/cmd/webd/rest/router/handlers/responder"
	"github.com/mappcpd/web-services/cmd/webd/rest/router/middleware"
	"github.com/mappcpd/web-services/internal/attachments"
	"github.com/mappcpd/web-services/internal/fileset"
	"github.com/mappcpd/web-services/internal/generic"
	"github.com/mappcpd/web-services/internal/member"
	"github.com/mappcpd/web-services/internal/notes"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/platform/s3"
	"github.com/mappcpd/web-services/internal/resources"
)

// AdminTest is a test endpoint
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
		res, err = member.DocMembersLimit(query, projection, limit)
	} else {
		res, err = member.DocMembersAll(query, projection)
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
		res, err = member.DocMembersLimit(f.Query, f.Projection, f.Limit)
	} else {
		res, err = member.DocMembersAll(f.Query, f.Projection)
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
	m, err := member.MemberByID(int(id))

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

		err := member.UpdateMember(j)
		if err != nil {
			p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
			p.Send(w)
			return
		}

		m, _ = member.MemberByID(int(id)) // Re-fetch
		member.SyncMember(m)              // Sync to doc db

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
	ns, err := notes.MemberNotes(id)
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Description}
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
	d, err := notes.NoteByID(id)
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Description}
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
	m, err := member.MemberByID(int(id))
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Description}
		p.Data = m
		member.SyncMember(m)
	}

	p.Send(w)
}

// AdminIDList fetches a list of all member ids from MySQL
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
		p.Message = responder.Message{http.StatusOK, "success", "List of ids from table: " + t + ", db: " + datastore.MySQL.Description}
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
		SuccessIDs []int
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
			failCount++
			continue
		}
		// otherwise, add the id to the list
		data.SuccessIDs = append(data.SuccessIDs, id)
		successCount++
	}

	p.Message = responder.Message{http.StatusOK, "success", "Batch completed - see Failures for errors"}
	p.Meta = map[string]int{
		"failed":     failCount,
		"successful": successCount,
	}
	p.Data = data
	p.Send(w)
}

// AdminNotesAttachmentRequest handles a request for a signed url to upload a notes attachment
func AdminNotesAttachmentRequest(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	upload := struct {
		SignedRequest  string `json:"signedRequest"`
		VolumeFilePath string `json:"volumeFilePath"`
		FileName       string `json:"fileName"`
		FileType       string `json:"fileType"`
	}{
		FileName: r.FormValue("filename"),
		FileType: r.FormValue("filetype"),
	}

	// Check we have required query params
	if upload.FileName == "" || upload.FileType == "" {
		msg := "Problems with query params, should have: ?filename=___&filetype=___"
		p.Message = responder.Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	// This is admin so don't need the owner of the record however still check that the record exists
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		msg := "Missing or malformed id in url path - " + err.Error()
		p.Message = responder.Message{http.StatusBadRequest, "failed", msg}
	}

	_, err = notes.NoteByID(id)
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("No note found with id %d -", id) + err.Error()
		p.Message = responder.Message{http.StatusNotFound, "failed", msg}
		p.Send(w)
		return
	case err != nil:
		msg := "Database error - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Get current fileset for note attachments
	fs, err := fileset.NoteAttachment()
	if err != nil {
		msg := "Could not determine the storage information for note attachments - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Build FULL file path or 'key' in S3 parlance
	filePath := fs.Path + strconv.Itoa(id) + "/" + upload.FileName

	// Prepend the volume name to pass back to the client for subsequent file registration
	upload.VolumeFilePath = fs.Volume + filePath

	// get a signed request
	url, err := s3.PutRequest(filePath, fs.Volume)
	if err != nil {
		msg := "Error getting a signed request for upload " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	upload.SignedRequest = url

	p.Message = responder.Message{http.StatusOK, "success", "Signed request in data.signedRequest."}
	p.Data = upload
	p.Send(w)
}

// AdminNotesAttachmentRegister registers a file attachment for a note.
func AdminNotesAttachmentRegister(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	a := attachments.New()
	a.UserID = middleware.UserAuthToken.Claims.ID

	// Get the entity ID from URL path... This is admin so validate record exists but not ownership
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		msg := "Error getting id from url path - " + err.Error()
		p.Message = responder.Message{http.StatusBadRequest, "failed", msg}
	}
	_, err = notes.NoteByID(id)
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("No note found with id %d -", id) + err.Error()
		p.Message = responder.Message{http.StatusNotFound, "failed", msg}
		p.Send(w)
		return
	case err != nil:
		msg := "Database error - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	a.EntityID = id

	// Decode post body fields: "cleanFilename" and "cloudyFilename" into Attachment.
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		msg := "Could not decode json in request body - " + err.Error()
		p.Message = responder.Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	// Get current fileset for note attachments
	fs, err := fileset.NoteAttachment()
	if err != nil {
		msg := "Could not determine the storage information for note attachments - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	a.FileSet = fs

	// Register the attachment
	if err := a.Register(); err != nil {
		msg := "Error registering attachment - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	p.Message = responder.Message{http.StatusOK, "success", "Attachment registered"}
	p.Data = a
	p.Send(w)
}

// AdminResourcesAttachmentRequest handles a request for a signed url to upload a resource attachment
func AdminResourcesAttachmentRequest(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	upload := struct {
		SignedRequest  string `json:"signedRequest"`
		VolumeFilePath string `json:"volumeFilePath"`
		FileName       string `json:"fileName"`
		FileType       string `json:"fileType"`
	}{
		FileName: r.FormValue("filename"),
		FileType: r.FormValue("filetype"),
	}

	// Check we have required query params
	if upload.FileName == "" || upload.FileType == "" {
		msg := "Problems with query params, should have: ?filename=___&filetype=___"
		p.Message = responder.Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	// This is admin so don't need the owner of the record however still check that the record exists
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		msg := "Missing or malformed id in url path - " + err.Error()
		p.Message = responder.Message{http.StatusBadRequest, "failed", msg}
	}

	_, err = resources.ResourceByID(id)
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("No resource found with id %d -", id) + err.Error()
		p.Message = responder.Message{http.StatusNotFound, "failed", msg}
		p.Send(w)
		return
	case err != nil:
		msg := "Database error - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Get current fileset for note attachments
	fs, err := fileset.ResourceAttachment()
	if err != nil {
		msg := "Could not determine the storage information for resource attachments - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Build FULL file path or 'key' in S3 parlance
	filePath := fs.Path + strconv.Itoa(id) + "/" + upload.FileName

	// Prepend the volume name to pass back to the client for subsequent file registration
	upload.VolumeFilePath = fs.Volume + filePath

	// get a signed request
	url, err := s3.PutRequest(filePath, fs.Volume)
	if err != nil {
		msg := "Error getting a signed request for upload " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	upload.SignedRequest = url

	p.Message = responder.Message{http.StatusOK, "success", "Signed request in data.signedRequest."}
	p.Data = upload
	p.Send(w)
}

// AdminResourcesAttachmentRegister registers a file attachment for a resource. If ?thumbnail=1 is passed on the
// url then the resource file is designated as a thumbnail by setting thumbnail flag to 1 in db.
func AdminResourcesAttachmentRegister(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	a := attachments.New()
	a.UserID = middleware.UserAuthToken.Claims.ID

	// Get the entity ID from URL path... This is admin so validate record exists but not ownership
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		msg := "Error getting id from url path - " + err.Error()
		p.Message = responder.Message{http.StatusBadRequest, "failed", msg}
	}
	_, err = resources.ResourceByID(id)
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("No resource found with id %d -", id) + err.Error()
		p.Message = responder.Message{http.StatusNotFound, "failed", msg}
		p.Send(w)
		return
	case err != nil:
		msg := "Database error - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	a.EntityID = id

	// Decode post body fields: "cleanFilename" and "cloudyFilename" into Attachment.
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		msg := "Could not decode json in request body - " + err.Error()
		p.Message = responder.Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	// Get current fileset for resource attachments
	fs, err := fileset.ResourceAttachment()
	if err != nil {
		msg := "Could not determine the storage information for resource attachments - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	a.FileSet = fs

	// Check if it is a thumbnail
	var flag string
	if r.FormValue("thumbnail") == "1" {
		flag = "thumbnail"
	}

	// Register the attachment
	if err := a.Register(flag); err != nil {
		msg := "Error registering attachment - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	p.Message = responder.Message{http.StatusOK, "success", "Attachment registered"}
	p.Data = a
	p.Send(w)
}
