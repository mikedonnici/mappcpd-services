package rest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mappcpd/web-services/internal/attachments"
	"github.com/mappcpd/web-services/internal/fileset"
	"github.com/mappcpd/web-services/internal/generic"
	"github.com/mappcpd/web-services/internal/member"
	"github.com/mappcpd/web-services/internal/note"
	"github.com/mappcpd/web-services/internal/platform/s3"
	"github.com/mappcpd/web-services/internal/resource"
)

// AdminTest is a test endpoint
func AdminTest(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)
	p.Message = Message{http.StatusOK, "success", "Hi Admin!"}
	p.Send(w)
}

// AdminMembersSearch searches the document database based on GET request
// Taking advantage of the complexity of Mongo queries and pipes is
// difficult using URI parameters. So this is an attempt however will also
// implement a POST version below to allow for a complete JSON query doc
// // to be submitted. Being totally RESTful is not as important  as this
// API is for DB access at this stage.
func AdminMembersSearch(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	var err error
	var query map[string]interface{}
	var projection map[string]interface{}

	// Query
	query, err = queryParams(r.FormValue("q"))
	if err != nil {
		if err != nil {
			p.Message = Message{
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
			p.Message = Message{
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
		res, err = member.DocMembersLimit(DS, query, projection, limit)
	} else {
		res, err = member.DocMembersAll(DS, query, projection)
	}

	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	c := len(res)
	p.Meta = DocMeta{c, query, projection}
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

	p := NewResponder(UserAuthToken.Encoded)

	// Pull the JSON body out of the request
	decoder := json.NewDecoder(r.Body)
	var f Find
	err := decoder.Decode(&f)
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failure", errMessageDecodeJSON}
		p.Send(w)
		return
	}

	// Limit query
	var res []interface{}
	if f.Limit > 0 {
		res, err = member.DocMembersLimit(DS, f.Query, f.Projection, f.Limit)
	} else {
		res, err = member.DocMembersAll(DS, f.Query, f.Projection)
	}
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	c := len(res)
	p.Meta = DocMeta{c, f.Query, f.Projection}
	p.Data = res
	p.Send(w)
}

// AdminMembersNotes fetches all Notes belonging to a Member
func AdminMembersNotes(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Response
	ns, err := note.ByMemberID(DS, id)
	switch {
	case err == sql.ErrNoRows:
		p.Message = Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
		p.Data = ns
	}

	p.Send(w)
}

// AdminNotes fetches a single Note record by Note ID
func AdminNotes(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Response
	d, err := note.ByID(DS, id)
	switch {
	case err == sql.ErrNoRows:
		p.Message = Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
		p.Data = d
	}

	p.Send(w)
}

// AdminMembersID fetches a member record from the MySQLConnection DB, by id
func AdminMembersID(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Get the Member record
	m, err := member.ByID(DS, int(id))
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
		p.Data = m
		member.SyncMember(DS, m)
	}

	p.Send(w)
}

// AdminIDList fetches a list of all member ids from MySQL
func AdminIDList(w http.ResponseWriter, req *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Request - requires at least the 't' query to specify the table name
	// and can have the option 'f' as raw HTML filter
	t := req.FormValue("t")
	if t == "" {
		p.Message = Message{http.StatusBadRequest, "failed", "Requires ?t=[table_name], optional &f=[sql_filter]"}
		p.Send(w)
		return
	}

	// Optional filter
	f := req.FormValue("f")

	// Get the Member record
	ii, err := generic.GetIDs(DS, t, f)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = Message{http.StatusOK, "success", "All of ids from table: " + t + ", db: "}
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
		Data resource.Resources `json:"data"`
	}
	b := batch{}

	p := NewResponder(UserAuthToken.Encoded)

	// Pull the JSON body out of the request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failure", errMessageDecodeJSON}
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
		r := resource.Resource{}
		r = v
		id, err := r.Save(DS)
		if err != nil {
			data.Failures[r.Name] = err.Error()
			failCount++
			continue
		}
		// otherwise, add the id to the list
		data.SuccessIDs = append(data.SuccessIDs, id)
		successCount++
	}

	p.Message = Message{http.StatusOK, "success", "Batch completed - see Failures for errors"}
	p.Meta = map[string]int{
		"failed":     failCount,
		"successful": successCount,
	}
	p.Data = data
	p.Send(w)
}

// AdminNotesAttachmentRequest handles a request for a signed url to upload a notes attachment
func AdminNotesAttachmentRequest(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	upload := struct {
		SignedRequest  string `json:"signedRequest"`
		VolumeFilePath string `json:"volumeFilePath"`
		FileName       string `json:"fileName"`
		FileType       string `json:"fileType"`
	}{
		FileName: r.FormValue("filename"),
		FileType: r.FormValue("filetype"),
	}

	// Decode we have required query params
	if upload.FileName == "" || upload.FileType == "" {
		msg := "Problems with query params, should have: ?filename=___&filetype=___"
		p.Message = Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	// This is admin so don't need the owner of the record however still check that the record exists
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		msg := "Missing or malformed id in url path - " + err.Error()
		p.Message = Message{http.StatusBadRequest, "failed", msg}
	}

	_, err = note.ByID(DS, id)
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("No note found with id %d -", id) + err.Error()
		p.Message = Message{http.StatusNotFound, "failed", msg}
		p.Send(w)
		return
	case err != nil:
		msg := "Database error - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Get current fileset for note attachments
	fs, err := fileset.NoteAttachment(DS)
	if err != nil {
		msg := "Could not determine the storage information for note attachments - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
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
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	upload.SignedRequest = url

	p.Message = Message{http.StatusOK, "success", "Signed request in data.signedRequest."}
	p.Data = upload
	p.Send(w)
}

// AdminNotesAttachmentRegister registers a file attachment for a note.
func AdminNotesAttachmentRegister(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	a := attachments.New()
	a.UserID = UserAuthToken.Claims.ID

	// Get the entity ID from URL path... This is admin so validate record exists but not ownership
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		msg := "Error getting id from url path - " + err.Error()
		p.Message = Message{http.StatusBadRequest, "failed", msg}
	}
	_, err = note.ByID(DS, id)
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("No note found with id %d -", id) + err.Error()
		p.Message = Message{http.StatusNotFound, "failed", msg}
		p.Send(w)
		return
	case err != nil:
		msg := "Database error - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	a.EntityID = id

	// Decode post body fields: "cleanFilename" and "cloudyFilename" into Attachment.
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		msg := "Could not decode json in request body - " + err.Error()
		p.Message = Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	// Get current fileset for note attachments
	fs, err := fileset.NoteAttachment(DS)
	if err != nil {
		msg := "Could not determine the storage information for note attachments - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	a.FileSet = fs

	// Register the attachment
	if err := a.Register(DS); err != nil {
		msg := "Error registering attachment - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Attachment registered"}
	p.Data = a
	p.Send(w)
}

// AdminResourcesAttachmentRequest handles a request for a signed url to upload a resource attachment
func AdminResourcesAttachmentRequest(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	upload := struct {
		SignedRequest  string `json:"signedRequest"`
		VolumeFilePath string `json:"volumeFilePath"`
		FileName       string `json:"fileName"`
		FileType       string `json:"fileType"`
	}{
		FileName: r.FormValue("filename"),
		FileType: r.FormValue("filetype"),
	}

	// Decode we have required query params
	if upload.FileName == "" || upload.FileType == "" {
		msg := "Problems with query params, should have: ?filename=___&filetype=___"
		p.Message = Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	// This is admin so don't need the owner of the record however still check that the record exists
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		msg := "Missing or malformed id in url path - " + err.Error()
		p.Message = Message{http.StatusBadRequest, "failed", msg}
	}

	_, err = resource.ByID(DS, id)
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("No resource found with id %d -", id) + err.Error()
		p.Message = Message{http.StatusNotFound, "failed", msg}
		p.Send(w)
		return
	case err != nil:
		msg := "Database error - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Get current fileset for note attachments
	fs, err := fileset.ResourceAttachment(DS)
	if err != nil {
		msg := "Could not determine the storage information for resource attachments - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
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
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	upload.SignedRequest = url

	p.Message = Message{http.StatusOK, "success", "Signed request in data.signedRequest."}
	p.Data = upload
	p.Send(w)
}

// AdminResourcesAttachmentRegister registers a file attachment for a resource. If ?thumbnail=1 is passed on the
// url then the resource file is designated as a thumbnail by setting thumbnail flag to 1 in db.
func AdminResourcesAttachmentRegister(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	a := attachments.New()
	a.UserID = UserAuthToken.Claims.ID

	// Get the entity ID from URL path... This is admin so validate record exists but not ownership
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		msg := "Error getting id from url path - " + err.Error()
		p.Message = Message{http.StatusBadRequest, "failed", msg}
	}
	_, err = resource.ByID(DS, id)
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("No resource found with id %d -", id) + err.Error()
		p.Message = Message{http.StatusNotFound, "failed", msg}
		p.Send(w)
		return
	case err != nil:
		msg := "Database error - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	a.EntityID = id

	// Decode post body fields: "cleanFilename" and "cloudyFilename" into Attachment.
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		msg := "Could not decode json in request body - " + err.Error()
		p.Message = Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	// Get current fileset for resource attachments
	fs, err := fileset.ResourceAttachment(DS)
	if err != nil {
		msg := "Could not determine the storage information for resource attachments - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	a.FileSet = fs

	// Decode if it is a thumbnail
	var flag string
	if r.FormValue("thumbnail") == "1" {
		flag = "thumbnail"
	}

	// Register the attachment
	if err := a.Register(DS, flag); err != nil {
		msg := "Error registering attachment - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Attachment registered"}
	p.Data = a
	p.Send(w)
}
