package rest

import (
	"log"
	"strconv"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mappcpd/web-services/internal/platform/jwt"
)

// AuthorizeScope checks token claims 'scope' field for one or more string values and returns
// false if any are missing.
// TODO: Authorize maybe should be some kind of middleware once we implement sub routers
func AuthorizeScope(w http.ResponseWriter, r *http.Request, s ...string) bool {

	p := Payload{}

	// Get token string from header
	a := r.Header.Get("Authorization")
	t, err := jwt.FromHeader(a)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return false
	}

	// Create an AuthToken value from the token string
	at, err := jwt.Check(t)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return false
	}

	// Now we can check the Scope claims for the strings passed in, and respond
	// on first negative - ie as soon as we find a missing scope authorization
	for i := range s {
		// "self" is a key word to check that the id on the url matches the id in the token
		if s[i] == "self" {
			// Get the ID off the url...
			v := mux.Vars(r)
			id, err := strconv.Atoi(v["id"])
			if err != nil {
				log.Printf("Authorize() error: %s", err.Error())
				return false
			}
			// Compare with the ID in the token
			if id != at.Claims.ID {
				return false
			}
		} else {
			if at.CheckScope(s[i]) == false {
				return false
			}
		}
	}

	// All clear
	return true
}

// AuthorizeID checks the member id passed in matches the token ID. This is used when a
// request is made that related to a record owned by a member. For example:
// GET /v1/m/activities/1234/attachments is requesting the attachment files for activity '1234'. In order to verify
// that the logged in member owns the record we currently fetch the activity record and compare the member_id with the
// user id in the token.
// todo - faster way to verify owner of an entity
func AuthorizeID(w http.ResponseWriter, r *http.Request, mid int) bool {

	p := Payload{}

	// Get token string from header
	a := r.Header.Get("Authorization")
	t, err := jwt.FromHeader(a)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return false
	}

	// Create an AuthToken value from the token string
	at, err := jwt.Check(t)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return false
	}

	// Now check the member ID passed in matches the token id
	if mid != at.Claims.ID {
		return false
	}

	return true
}
