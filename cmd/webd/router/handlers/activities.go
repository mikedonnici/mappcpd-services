package handlers

import (
	"fmt"
	"strconv"
	"time"

	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/mappcpd/web-services/cmd/webd/router/handlers/responder"
	"github.com/mappcpd/web-services/cmd/webd/router/middleware"
	"github.com/mappcpd/web-services/internal/activities"
	//"github.com/mappcpd/web-services/internal/attachments"
	"github.com/mappcpd/web-services/internal/attachments"
	"github.com/mappcpd/web-services/internal/fileset"
	"github.com/mappcpd/web-services/internal/members"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// Activities fetches list of activity types
func Activities(w http.ResponseWriter, _ *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	al, err := activities.ActivityList()
	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
	p.Data = al
	m := make(map[string]interface{})
	m["count"] = len(al)
	m["description"] = "This is a list of Activity types for creating lists etc. The typeId is required for creating new Activity records"
	p.Meta = m
	p.Send(w)
}

// ActivitiesID fetches a single activity type by ID
func ActivitiesID(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	a, err := activities.ActivityByID(id)
	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
	p.Data = a
	m := make(map[string]interface{})
	m["description"] = "The typeId must included when creating new Activity records"
	p.Meta = m
	p.Send(w)
}

// Activities fetches a single activity record by id
func MembersActivitiesID(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Response
	a, err := members.MemberActivityByID(id)
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// Authorization - need  owner of the record
	if middleware.UserAuthToken.Claims.ID != a.MemberID {
		p.Message = responder.Message{http.StatusUnauthorized, "failed", "Token does not belong to the owner of resource"}
		p.Send(w)
		return
	}

	// All good
	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
	p.Data = a
	p.Send(w)
}

// MembersActivitiesAdd adds a new activity for the logged in member
func MembersActivitiesAdd(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Decode JSON body into NewActivity value
	a := members.MemberActivityRow{}
	a.MemberID = middleware.UserAuthToken.Claims.ID
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		msg := "Error decoding JSON: " + err.Error() + ". Check the format of request body."
		p.Message = responder.Message{http.StatusBadRequest, "failure", msg}
		p.Send(w)
		return
	}

	aid, err := members.AddMemberActivity(a)
	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failure", err.Error()}
		p.Send(w)
		return
	}

	// Fetch the new record for return
	ar, err := members.MemberActivityByID(int(aid))
	if err != nil {
		msg := "Could not fetch the new record"
		p.Message = responder.Message{http.StatusInternalServerError, "failure", msg + " " + err.Error()}
		p.Send(w)
		return
	}

	msg := fmt.Sprintf("Added a new activity (id: %v) for member (id: %v)", aid, middleware.UserAuthToken.Claims.ID)
	p.Message = responder.Message{http.StatusCreated, "success", msg}
	p.Data = ar
	p.Send(w)
}

// MembersActivitiesUpdate updates an existing activity for the logged in member.
// First we fetch the existing record into an Activity, and then replace the update fields with
// // new values - this will be validated in the same way as a new activity and can also
// update one to many fields.
func MembersActivitiesUpdate(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Get activity id from path... and make it an int
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Fetch the original activity record
	a, err := members.MemberActivityRowByID(id)
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// Authorization - need  owner of the record
	if middleware.UserAuthToken.Claims.ID != a.MemberID {
		p.Message = responder.Message{http.StatusUnauthorized, "failed", "Token does not belong to the owner of resource"}
		p.Send(w)
		return
	}

	// activity update posted in JSON body
	au := members.MemberActivityRow{}
	err = json.NewDecoder(r.Body).Decode(&au)
	if err != nil {
		msg := "Error decoding JSON: " + err.Error() + ". Check the format of request body."
		p.Message = responder.Message{http.StatusBadRequest, "failure", msg}
		p.Send(w)
		return
	}

	// Merge the original into the new record to fill in any blanks. The merge package
	// will only overwrite 'zero' values, so the updates are kept, and the nil values
	// back filled with the original values
	err = mergo.Merge(&au, a)
	if err != nil {
		fmt.Println(err)
	}

	// Update the activity record
	err = members.UpdateMemberActivity(au)
	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failure", err.Error()}
		p.Send(w)
		return
	}

	// Fetch the updated record for response
	ar, err := members.MemberActivityByID(id)
	if err != nil {
		msg := "Could not fetch the updated record"
		p.Message = responder.Message{http.StatusInternalServerError, "failure", msg + " " + err.Error()}
		p.Send(w)
		return
	}

	msg := fmt.Sprintf("Updated activity (id: %v) for member (id: %v)", id, middleware.UserAuthToken.Claims.ID)
	p.Message = responder.Message{http.StatusOK, "success", msg}
	p.Data = ar
	p.Send(w)
}

// MembersActivitiesRecurring fetches the member's recurring activities (if any) stored in MongoDB
func MembersActivitiesRecurring(w http.ResponseWriter, _ *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	ra, err := members.MemberRecurring(middleware.UserAuthToken.Claims.ID)
	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failed", "Failed to initialise a value of type MemberRecurring -" + err.Error()}
		p.Send(w)
		return
	}

	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MongoDB.Source}
	p.Meta = map[string]int{"count": len(ra.Activities)}
	p.Data = ra
	p.Send(w)
}

// MembersActivitiesRecurringAdd adds a new recurring activity to the array in the Recurring doc that belongs to the member.
// Note that this function reads and writes only to MongoDB
func MembersActivitiesRecurringAdd(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Get user id from token
	id := middleware.UserAuthToken.Claims.ID

	// Fetch the recurring activity doc for this user first
	ra, err := members.MemberRecurring(id)
	if err != nil {
		msg := "MembersActivitiesRecurringAdd() Failed to initialise a value of type Recurring -" + err.Error()
		fmt.Println(msg)
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	ra.UpdatedAt = time.Now()

	// Decode the new activity from POST body...
	b := members.RecurringActivity{}
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		msg := "MembersActivitiesRecurringAdd() failed to decode body -" + err.Error()
		fmt.Println(msg)
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	b.ID = bson.NewObjectId()
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()

	// Add the new recurring activity to the list...
	ra.Activities = append(ra.Activities, b)

	// ... and save
	err = ra.Save()
	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MongoDB.Source}
	p.Meta = map[string]int{"count": len(ra.Activities)}
	p.Data = ra
	p.Send(w)
}

// MembersActivitiesRecurringRemove removes a recurring activity from the Recurring doc. Not it is not removing a
// doc in the collection, only one element from the array of recurring activities in the doc that belongs to the member
func MembersActivitiesRecurringRemove(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Get user id from token
	id := middleware.UserAuthToken.Claims.ID

	// Fetch the recurring activity doc for this user first
	ra, err := members.MemberRecurring(id)
	if err != nil {
		msg := "MembersActivitiesRecurringAdd() Failed to initialise a value of type Recurring -" + err.Error()
		fmt.Println(msg)
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Remove the recurring activity identified by the _id on url...
	_id := mux.Vars(r)["_id"]

	err = ra.RemoveActivity(_id)
	if err == mgo.ErrNotFound {
		msg := "No activity was found with id " + _id + " - it may have been already deleted"
		p.Message = responder.Message{http.StatusNotFound, "failure", msg + "... data retrieved from " + datastore.MongoDB.Source}

	} else if err != nil {
		msg := "An error occured - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failure", msg + "... data retrieved from " + datastore.MongoDB.Source}
	} else {
		p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MongoDB.Source}
	}

	p.Meta = map[string]int{"count": len(ra.Activities)}
	p.Data = ra
	p.Send(w)
}

// MembersActivitiesRecurringRecorder records a member activity based on a recurring activity.
// It creates a new member activity and then increments the next scheduled date for the recurring activity.
// If ?slip=1 is passed on the url then it will
func MembersActivitiesRecurringRecorder(w http.ResponseWriter, r *http.Request) {

	p := responder.Payload{}

	// Get the member's recurring activities. Strictly speaking we don't need the member id to do this
	// as we can select the document based on the recurring activity id. However, this ensures that the recurring
	// activity belongs to the member - however slim the chances of guessing an ObjectID!
	id := middleware.UserAuthToken.Claims.ID
	ra, err := members.MemberRecurring(id)
	if err != nil {
		msg := "MembersActivitiesRecurringAdd() Failed to initialise a value of type Recurring -" + err.Error()
		fmt.Println(msg)
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Record (or skip) the target activity (_id on url), and increment the schedule
	_id := mux.Vars(r)["_id"]
	q := r.URL.Query()
	// ?skip=anything will do...
	if len(q["skip"]) > 0 {
		fmt.Println("Skip recurring activity...")
		err = ra.Skip(_id)
	} else {
		fmt.Println("Record recurring activity...")
		err = ra.Record(_id)
	}

	if err != nil {
		p.Message = responder.Message{http.StatusNotFound, "failed", "Could not record or skip recurring activity with id " + _id + " - " + err.Error()}
		p.Meta = map[string]int{"count": len(ra.Activities)}
		p.Data = ra
		p.Send(w)
		return
	}

	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MongoDB.Source}
	p.Meta = map[string]int{"count": len(ra.Activities)}
	p.Data = ra
	p.Send(w)
}

// MembersActivitiesAttachmentRequest handles request for a signed URL to upload an attachment for a CPD activity
func MembersActivitiesAttachmentRequest(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Return the URL Query string params for the caller's convenience, and signedURL
	upload := struct {
		Volume        string `json:"volume"`
		Path          string `json:"path"`
		FileName      string `json:"fileName"`
		FileType      string `json:"fileType"`
		SignedRequest string `json:"signedRequest"`
	}{
		FileName: r.FormValue("filename"),
		FileType: r.FormValue("filetype"),
	}

	// Check we have required query params
	if upload.FileName == "" || upload.FileType == "" {
		msg := "Problems with query params, should have: ?filename=___&filetype=___"
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Check logged in member owns the activity record
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	a, err := members.MemberActivityByID(id)
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// Authorization - need  owner of the record
	if middleware.UserAuthToken.Claims.ID != a.MemberID {
		p.Message = responder.Message{http.StatusUnauthorized, "failed", "Token does not belong to the owner of resource"}
		p.Send(w)
		return
	}

	// Have taken some 'load' off the client - rather than the client having to know what file sets we have we can let
	// the API work it out. We know that a CPD activity attachment will be registered in the ce_m_activity_attachment table
	// so lookup the current file set for that entity
	fs, err := fileset.New("ce_m_activity_attachment")
	if err != nil {
		msg := "Could not determine the storage information for activity attachments - " + err.Error()
		p.Message = responder.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	upload.Path = fs.Path
	upload.Volume = fs.Volume

	// passed all required checks so ok to get a signed request
	url, err := attachments.S3PutRequest(upload.Path, upload.Volume)
	if err != nil {
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	upload.SignedRequest = url

	p.Message = responder.Message{http.StatusOK, "success", "Signed request in data.signedRequest."}
	p.Data = upload
	p.Send(w)
}

// MembersActivitiesAttachmentAdd registers an uploaded file in the database. It creates an association between the
// uploaded file and the relevant database entity thus creating the 'attachment'.
// Todo... this needs to be simplified as we don't need to pass the entity name or id in the POSt body for member
func MembersActivitiesAttachmentAdd(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Get activity id from path... and make it an int
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = responder.Message{http.StatusBadRequest, "failed", err.Error()}
		p.Send(w)
		return
	}

	// The only thing we will need to do with the activity id is verify that it belongs
	// the the same member that owns the JWT
	// todo look up activity and verify owner...
	fmt.Println("Look up activity id", id, " and verify the owner")

	// decode post body into a struct
	var a attachments.Attachment

	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		msg := "Could not decode json in request body - " + err.Error()
		p.Message = responder.Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}
	fmt.Println(a)

	if err := a.Register(); err != nil {
		msg := "Error registering attachment - " + err.Error()
		p.Message = responder.Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	p.Message = responder.Message{http.StatusOK, "success", "Attachment registered"}
	p.Data = a
	p.Send(w)

}
