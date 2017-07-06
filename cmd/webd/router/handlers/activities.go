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

	_json "github.com/mappcpd/web-services/cmd/webd/router/handlers/json"
	_mw "github.com/mappcpd/web-services/cmd/webd/router/middleware"
	a_ "github.com/mappcpd/web-services/internal/activities"
	//"github.com/mappcpd/web-services/internal/attachments"
	"github.com/mappcpd/web-services/internal/attachments"
	m_ "github.com/mappcpd/web-services/internal/members"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// Activities fetches list of activity types
func Activities(w http.ResponseWriter, _ *http.Request) {

	p := _json.NewPayload(_mw.UserAuthToken.Token)

	al, err := a_.ActivityList()
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
	p.Data = al
	m := make(map[string]interface{})
	m["count"] = len(al)
	m["description"] = "This is a list of Activity types for creating lists etc. The typeId is required for creating new Activity records"
	p.Meta = m
	p.Send(w)
}

// ActivitiesID fetches a single activity type by ID
func ActivitiesID(w http.ResponseWriter, r *http.Request) {

	p := _json.NewPayload(_mw.UserAuthToken.Token)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = _json.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	a, err := a_.ActivityByID(id)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
	p.Data = a
	m := make(map[string]interface{})
	m["description"] = "The typeId must included when creating new Activity records"
	p.Meta = m
	p.Send(w)
}

// Activities fetches a single activity record by id
func MembersActivitiesID(w http.ResponseWriter, r *http.Request) {

	p := _json.NewPayload(_mw.UserAuthToken.Token)

	// Request - convert id from string to int type
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = _json.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Response
	a, err := m_.MemberActivityByID(id)
	switch {
	case err == sql.ErrNoRows:
		p.Message = _json.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// Authorization - need  owner of the record
	if _mw.UserAuthToken.Claims.ID != a.MemberID {
		p.Message = _json.Message{http.StatusUnauthorized, "failed", "Token does not belong to the owner of resource"}
		p.Send(w)
		return
	}

	// All good
	p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
	p.Data = a
	p.Send(w)
}

// MembersActivitiesAdd adds a new activity for the logged in member
func MembersActivitiesAdd(w http.ResponseWriter, r *http.Request) {

	p := _json.NewPayload(_mw.UserAuthToken.Token)

	// Decode JSON body into NewActivity value
	a := m_.MemberActivityRow{}
	a.MemberID = _mw.UserAuthToken.Claims.ID
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		msg := "Error decoding JSON: " + err.Error() + ". Check the format of request body."
		p.Message = _json.Message{http.StatusBadRequest, "failure", msg}
		p.Send(w)
		return
	}

	aid, err := m_.AddMemberActivity(a)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failure", err.Error()}
		p.Send(w)
		return
	}

	// Fetch the new record for return
	ar, err := m_.MemberActivityByID(int(aid))
	if err != nil {
		msg := "Could not fetch the new record"
		p.Message = _json.Message{http.StatusInternalServerError, "failure", msg + " " + err.Error()}
		p.Send(w)
		return
	}

	msg := fmt.Sprintf("Added a new activity (id: %v) for member (id: %v)", aid, _mw.UserAuthToken.Claims.ID)
	p.Message = _json.Message{http.StatusCreated, "success", msg}
	p.Data = ar
	p.Send(w)
}

// MembersActivitiesUpdate updates an existing activity for the logged in member.
// First we fetch the existing record into an Activity, and then replace the update fields with
// // new values - this will be validated in the same way as a new activity and can also
// update one to many fields.
func MembersActivitiesUpdate(w http.ResponseWriter, r *http.Request) {

	p := _json.NewPayload(_mw.UserAuthToken.Token)

	// Get activity id from path... and make it an int
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = _json.Message{http.StatusBadRequest, "failed", err.Error()}
	}

	// Fetch the original activity record
	a, err := m_.MemberActivityRowByID(id)
	switch {
	case err == sql.ErrNoRows:
		p.Message = _json.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// Authorization - need  owner of the record
	if _mw.UserAuthToken.Claims.ID != a.MemberID {
		p.Message = _json.Message{http.StatusUnauthorized, "failed", "Token does not belong to the owner of resource"}
		p.Send(w)
		return
	}

	// activity update posted in JSON body
	au := m_.MemberActivityRow{}
	err = json.NewDecoder(r.Body).Decode(&au)
	if err != nil {
		msg := "Error decoding JSON: " + err.Error() + ". Check the format of request body."
		p.Message = _json.Message{http.StatusBadRequest, "failure", msg}
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
	err = m_.UpdateMemberActivity(au)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failure", err.Error()}
		p.Send(w)
		return
	}

	// Fetch the updated record for response
	ar, err := m_.MemberActivityByID(id)
	if err != nil {
		msg := "Could not fetch the updated record"
		p.Message = _json.Message{http.StatusInternalServerError, "failure", msg + " " + err.Error()}
		p.Send(w)
		return
	}

	msg := fmt.Sprintf("Updated activity (id: %v) for member (id: %v)", id, _mw.UserAuthToken.Claims.ID)
	p.Message = _json.Message{http.StatusOK, "success", msg}
	p.Data = ar
	p.Send(w)
}

// MembersActivitiesRecurring fetches the member's recurring activities (if any) stored in MongoDB
func MembersActivitiesRecurring(w http.ResponseWriter, _ *http.Request) {

	p := _json.NewPayload(_mw.UserAuthToken.Token)

	ra, err := m_.MemberRecurring(_mw.UserAuthToken.Claims.ID)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", "Failed to initialise a value of type MemberRecurring -" + err.Error()}
		p.Send(w)
		return
	}

	p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MongoDB.Source}
	p.Meta = map[string]int{"count": len(ra.Activities)}
	p.Data = ra
	p.Send(w)
}

// MembersActivitiesRecurringAdd adds a new recurring activity to the array in the Recurring doc that belongs to the member.
// Note that this function reads and writes only to MongoDB
func MembersActivitiesRecurringAdd(w http.ResponseWriter, r *http.Request) {

	p := _json.NewPayload(_mw.UserAuthToken.Token)

	// Get user id from token
	id := _mw.UserAuthToken.Claims.ID

	// Fetch the recurring activity doc for this user first
	ra, err := m_.MemberRecurring(id)
	if err != nil {
		msg := "MembersActivitiesRecurringAdd() Failed to initialise a value of type Recurring -" + err.Error()
		fmt.Println(msg)
		p.Message = _json.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}
	ra.UpdatedAt = time.Now()

	// Decode the new activity from POST body...
	b := m_.RecurringActivity{}
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		msg := "MembersActivitiesRecurringAdd() failed to decode body -" + err.Error()
		fmt.Println(msg)
		p.Message = _json.Message{http.StatusInternalServerError, "failed", msg}
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
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MongoDB.Source}
	p.Meta = map[string]int{"count": len(ra.Activities)}
	p.Data = ra
	p.Send(w)
}

// MembersActivitiesRecurringRemove removes a recurring activity from the Recurring doc. Not it is not removing a
// doc in the collection, only one element from the array of recurring activities in the doc that belongs to the member
func MembersActivitiesRecurringRemove(w http.ResponseWriter, r *http.Request) {

	p := _json.NewPayload(_mw.UserAuthToken.Token)

	// Get user id from token
	id := _mw.UserAuthToken.Claims.ID

	// Fetch the recurring activity doc for this user first
	ra, err := m_.MemberRecurring(id)
	if err != nil {
		msg := "MembersActivitiesRecurringAdd() Failed to initialise a value of type Recurring -" + err.Error()
		fmt.Println(msg)
		p.Message = _json.Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	// Remove the recurring activity identified by the _id on url...
	_id := mux.Vars(r)["_id"]

	err = ra.RemoveActivity(_id)
	if err == mgo.ErrNotFound {
		msg := "No activity was found with id " + _id + " - it may have been already deleted"
		p.Message = _json.Message{http.StatusNotFound, "failure", msg + "... data retrieved from " + datastore.MongoDB.Source}

	} else if err != nil {
		msg := "An error occured - " + err.Error()
		p.Message = _json.Message{http.StatusInternalServerError, "failure", msg + "... data retrieved from " + datastore.MongoDB.Source}
	} else {
		p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MongoDB.Source}
	}

	p.Meta = map[string]int{"count": len(ra.Activities)}
	p.Data = ra
	p.Send(w)
}

// MembersActivitiesRecurringRecorder records a member activity based on a recurring activity.
// It creates a new member activity and then increments the next scheduled date for the recurring activity.
// If ?slip=1 is passed on the url then it will
func MembersActivitiesRecurringRecorder(w http.ResponseWriter, r *http.Request) {

	p := _json.Payload{}

	// Get the member's recurring activities. Strictly speaking we don't need the member id to do this
	// as we can select the document based on the recurring activity id. However, this ensures that the recurring
	// activity belongs to the member - however slim the chances of guessing an ObjectID!
	id := _mw.UserAuthToken.Claims.ID
	ra, err := m_.MemberRecurring(id)
	if err != nil {
		msg := "MembersActivitiesRecurringAdd() Failed to initialise a value of type Recurring -" + err.Error()
		fmt.Println(msg)
		p.Message = _json.Message{http.StatusInternalServerError, "failed", msg}
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
		p.Message = _json.Message{http.StatusNotFound, "failed", "Could not record or skip recurring activity with id " + _id + " - " + err.Error()}
		p.Meta = map[string]int{"count": len(ra.Activities)}
		p.Data = ra
		p.Send(w)
		return
	}

	p.Message = _json.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MongoDB.Source}
	p.Meta = map[string]int{"count": len(ra.Activities)}
	p.Data = ra
	p.Send(w)
}

// MembersActivitiesAttachmentAdd registers a new attachment in the database. It has nothing to do with the actual upload.
func MembersActivitiesAttachmentAdd(w http.ResponseWriter, r *http.Request) {

	p := _json.NewPayload(_mw.UserAuthToken.Token)

	// Get activity id from path... and make it an int
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		p.Message = _json.Message{http.StatusBadRequest, "failed", err.Error()}
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
		p.Message = _json.Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}
	fmt.Println(a)

	if err := a.Register(); err != nil {
		msg := "Error registering attachment - " + err.Error()
		p.Message = _json.Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	p.Message = _json.Message{http.StatusOK, "success", "Attachment registered"}
	p.Data = a
	p.Send(w)

}
