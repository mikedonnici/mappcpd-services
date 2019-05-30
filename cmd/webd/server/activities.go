package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cardiacsociety/web-services/internal/activity"
	"github.com/cardiacsociety/web-services/internal/attachments"
	"github.com/cardiacsociety/web-services/internal/cpd"
	"github.com/cardiacsociety/web-services/internal/fileset"
	"github.com/cardiacsociety/web-services/internal/platform/s3"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Activities fetches list of activity types
func Activities(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	al, err := activity.All(DS)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Data = al
	m := make(map[string]interface{})
	m["count"] = len(al)
	m["description"] = "This is a list of Activity types for creating lists etc. The typeId is required for creating new Activity records"
	p.Meta = m
	p.Send(w)
}

// ActivitiesID fetches a single activity type by ID
func ActivitiesID(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failed", err.Error()}
	}

	a, err := activity.ByID(DS, id)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Data = a
	m := make(map[string]interface{})
	m["description"] = "The typeId must included when creating new Activity records"
	p.Meta = m
	p.Send(w)
}

// MembersActivitiesID fetches a single activity record by id
func MembersActivitiesID(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Response
	a, err := cpd.ByID(DS, int(id))
	switch {
	case err == sql.ErrNoRows:
		p.Message = Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// Authorization - need  owner of the record
	if UserAuthToken.Claims.ID != a.MemberID {
		p.Message = Message{http.StatusUnauthorized, "failed", "Encoded does not belong to the owner of resource"}
		p.Send(w)
		return
	}

	// All good
	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Data = a
	p.Send(w)
}

// MembersActivitiesAdd adds a new activity for the logged in member
func MembersActivitiesAdd(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Decode JSON body into ActivityAttachment value
	a := cpd.Input{}
	a.MemberID = UserAuthToken.Claims.ID
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		msg := "Error decoding JSON: " + err.Error() + ". Decode the format of request body."
		p.Message = Message{http.StatusBadRequest, "failure", msg}
		p.Send(w)
		return
	}

	aid, err := cpd.Add(DS, a)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failure", err.Error()}
		p.Send(w)
		return
	}

	// Fetch the new record for return
	ar, err := cpd.ByID(DS, int(aid))
	if err != nil {
		msg := "Could not fetch the new record"
		p.Message = Message{http.StatusInternalServerError, "failure", msg + " " + err.Error()}
		p.Send(w)
		return
	}

	msg := fmt.Sprintf("Added a new activity (id: %v) for member (id: %v)", aid, UserAuthToken.Claims.ID)
	p.Message = Message{http.StatusCreated, "success", msg}
	p.Data = ar
	p.Send(w)
}

// MembersActivitiesUpdate updates an existing activity for the logged in member.
// First we fetch the existing record into an Activity, and then replace the update fields with
// new values - this will be validated in the same way as a new activity and can also
// update one to many fields.
func MembersActivitiesUpdate(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Get activity id from path... and make it an int
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Fetch the original activity record
	a, err := cpd.ByID(DS, int(id))
	switch {
	case err == sql.ErrNoRows:
		p.Message = Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// Authorization - need  owner of the record
	if UserAuthToken.Claims.ID != a.MemberID {
		p.Message = Message{http.StatusUnauthorized, "failed", "Encoded does not belong to the owner of resource"}
		p.Send(w)
		return
	}

	// Original activity - from above we have a CPD but need a subset of this - ie Input
	oa := cpd.Input{
		ID:          a.ID,
		MemberID:    a.MemberID,
		ActivityID:  a.Activity.ID,
		Evidence:    false,
		Date:        a.Date,
		Quantity:    a.CreditData.Quantity,
		UnitCredit:  a.CreditData.UnitCredit,
		Description: a.Description,
	}

	// new activity - ie, updated version posted in JSON body
	na := cpd.Input{}
	err = json.NewDecoder(r.Body).Decode(&na)
	if err != nil {
		msg := "Error decoding JSON: " + err.Error() + ". Decode the format of request body."
		p.Message = Message{http.StatusBadRequest, "failure", msg}
		p.Send(w)
		return
	}

	// Merge the original into the new record to fill in any blanks. The merge package
	// will only overwrite 'zero' values, so the updates are kept, and the nil values
	// back filled with the original values
	fmt.Println("Original:", oa)
	fmt.Println("New:", na)
	err = mergo.Merge(&na, oa)
	if err != nil {
		fmt.Println("Error merging activity fields: ", err)
	}
	fmt.Println("Original:", oa)
	fmt.Println("New:", na)

	// Update the activity record
	err = cpd.Update(DS, na)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failure", err.Error()}
		p.Send(w)
		return
	}

	// updated record - fetch for response
	ur, err := cpd.ByID(DS, int(id))
	if err != nil {
		msg := "Could not fetch the updated record"
		p.Message = Message{http.StatusInternalServerError, "failure", msg + " " + err.Error()}
		p.Send(w)
		return
	}

	msg := fmt.Sprintf("Updated activity (id: %v) for member (id: %v)", id, UserAuthToken.Claims.ID)
	p.Message = Message{http.StatusOK, "success", msg}
	p.Data = ur
	p.Send(w)
}

// MembersActivitiesRecurring fetches the member's recurring activities (if any) stored in MongoDB
func MembersActivitiesRecurring(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	ra, err := cpd.MemberRecurring(DS, UserAuthToken.Claims.ID)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", "Failed to initialise a value of type MemberRecurring -" + err.Error()}
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Meta = map[string]int{"count": len(ra.Activities)}
	p.Data = ra
	p.Send(w)
}

// MembersActivitiesRecurringAdd adds a new recurring activity to the array in the Recurring doc that belongs to the member.
// Note that this function reads and writes only to MongoDB
func MembersActivitiesRecurringAdd(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Get user id from token
	id := UserAuthToken.Claims.ID

	// Fetch the recurring activity doc for this user first
	ra, err := cpd.MemberRecurring(DS, id)
	if err != nil {
		msg := "MembersActivitiesRecurringAdd() Failed to initialise a value of type Recurring -" + err.Error()
		fmt.Println(msg)
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	ra.UpdatedAt = time.Now()

	// Decode the new activity from POST body...
	b := cpd.RecurringActivity{}
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		msg := "MembersActivitiesRecurringAdd() failed to decode body -" + err.Error()
		fmt.Println(msg)
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	b.ID = bson.NewObjectId()
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()

	// Add the new recurring activity to the list...
	ra.Activities = append(ra.Activities, b)

	// ... and save
	err = ra.Save(DS)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Meta = map[string]int{"count": len(ra.Activities)}
	p.Data = ra
	p.Send(w)
}

// MembersActivitiesRecurringRemove removes a recurring activity from the Recurring doc. Not it is not removing a
// doc in the collection, only one element from the array of recurring activities in the doc that belongs to the member
func MembersActivitiesRecurringRemove(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Get user id from token
	id := UserAuthToken.Claims.ID

	// Fetch the recurring activity doc for this user first
	ra, err := cpd.MemberRecurring(DS, id)
	if err != nil {
		msg := "MembersActivitiesRecurringAdd() Failed to initialise a value of type Recurring -" + err.Error()
		fmt.Println(msg)
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Remove the recurring activity identified by the _id on url...
	_id := mux.Vars(r)["_id"]

	err = ra.RemoveActivity(DS, _id)
	if err == mgo.ErrNotFound {
		msg := "No activity was found with id " + _id + " - it may have been already deleted"
		p.Message = Message{http.StatusNotFound, "failure", msg}

	} else if err != nil {
		msg := "An error occured - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failure", msg}
	} else {
		p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	}

	p.Meta = map[string]int{"count": len(ra.Activities)}
	p.Data = ra
	p.Send(w)
}

// MembersActivitiesRecurringRecorder records a member activity based on a recurring activity.
// It creates a new member activity and then increments the next scheduled date for the recurring activity.
// If ?slip=1 is passed on the url then it will
func MembersActivitiesRecurringRecorder(w http.ResponseWriter, r *http.Request) {

	p := Payload{}

	// Get the member's recurring activities. Strictly speaking we don't need the member id to do this
	// as we can select the document based on the recurring activity id. However, this ensures that the recurring
	// activity belongs to the member - however slim the chances of guessing an ObjectID!
	id := UserAuthToken.Claims.ID
	ra, err := cpd.MemberRecurring(DS, id)
	if err != nil {
		msg := "MembersActivitiesRecurringAdd() Failed to initialise a value of type Recurring -" + err.Error()
		fmt.Println(msg)
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// CPD (or skip) the target activity (_id on url), and increment the schedule
	_id := mux.Vars(r)["_id"]
	q := r.URL.Query()
	// ?skip=anything will do...
	if len(q["skip"]) > 0 {
		fmt.Println("Skip recurring activity...")
		err = ra.Skip(DS, _id)
	} else {
		fmt.Println("CPD recurring activity...")
		err = ra.Record(DS, _id)
	}

	if err != nil {
		p.Message = Message{http.StatusNotFound, "failed", "Could not record or skip recurring activity with id " + _id + " - " + err.Error()}
		p.Meta = map[string]int{"count": len(ra.Activities)}
		p.Data = ra
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Meta = map[string]int{"count": len(ra.Activities)}
	p.Data = ra
	p.Send(w)
}

// MembersActivitiesAttachmentRequest handles request for a signed URL to upload an attachment for a CPD activity
func MembersActivitiesAttachmentRequest(w http.ResponseWriter, r *http.Request) {

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
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Decode logged in member owns the activity record
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		msg := "Missing or malformed id in url path - " + err.Error()
		p.Message = Message{http.StatusBadRequest, "failed", msg}
	}

	a, err := cpd.ByID(DS, int(id))
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("No activity found with id %d -", id) + err.Error()
		p.Message = Message{http.StatusNotFound, "failed", msg}
		p.Send(w)
		return
	case err != nil:
		msg := "Database error - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Authorization - need  owner of the record
	if UserAuthToken.Claims.ID != a.MemberID {
		p.Message = Message{http.StatusUnauthorized, "failed", "Encoded does not belong to the owner of resource"}
		p.Send(w)
		return
	}

	// Get current fileset for activity attachments
	fs, err := fileset.ActivityAttachment(DS)
	if err != nil {
		msg := "Could not determine the storage information for activity attachments - " + err.Error()
		p.Message = Message{http.StatusBadRequest, "failed", msg}
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

// MembersActivitiesAttachmentRegister registers an uploaded file in the database.
func MembersActivitiesAttachmentRegister(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	a := attachments.New()
	// not required for this type of attachment but stick it on for good measure :)
	a.UserID = UserAuthToken.Claims.ID

	// Get the entity ID from URL path... This is admin so validate record exists but not ownership
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		msg := "Error getting id from url path - " + err.Error()
		p.Message = Message{http.StatusBadRequest, "failed", msg}
		p.Data = a
		p.Send(w)
		return
	}
	activity, err := cpd.ByID(DS, int(id))
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("No activity found with id %d -", id) + err.Error()
		p.Message = Message{http.StatusNotFound, "failed", msg}
		p.Data = a
		p.Send(w)
		return
	case err != nil:
		msg := "Database error - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Data = a
		p.Send(w)
		return
	}
	// CHECK OWNER!!
	if UserAuthToken.Claims.ID != activity.MemberID {
		p.Message = Message{http.StatusUnauthorized, "failed", "Encoded does not belong to the owner of this resource"}
		p.Data = a
		p.Send(w)
		return
	}
	a.EntityID = int(id)

	// Decode post body fields: "cleanFilename" and "cloudyFilename" into Attachment
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		msg := "Could not decode json in request body - " + err.Error()
		p.Message = Message{http.StatusBadRequest, "failed", msg}
		p.Data = a
		p.Send(w)
		return
	}

	// Get current fileset for activity attachments
	fs, err := fileset.ActivityAttachment(DS)
	if err != nil {
		msg := "Could not determine the storage information for activity attachments - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Data = a
		p.Send(w)
		return
	}
	a.FileSet = fs

	// Register the attachment
	if err := a.Register(DS); err != nil {
		msg := "Error registering attachment - " + err.Error()
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Data = a
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Attachment registered"}
	p.Data = a
	p.Send(w)
}
