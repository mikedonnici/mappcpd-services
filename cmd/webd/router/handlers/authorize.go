package handlers

import (
	"log"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"

	_json "github.com/mappcpd/web-services/cmd/webd/router/handlers/responder"
	j_ "github.com/mappcpd/web-services/internal/platform/jwt"
)

// AuthorizeScope checks token claims 'scope' field for one or more string values and returns
// false if any are missing.
// TODO: Authorize maybe should be some kind of middleware once we implement sub routers
func AuthorizeScope(w http.ResponseWriter, r *http.Request, s ...string) bool {

	p := _json.Payload{}

	// Get token string from header
	a := r.Header.Get("Authorization")
	t, err := j_.FromHeader(a)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return false
	}

	// Create an AuthToken value from the token string
	at, err := j_.Check(t)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
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
// record linked to a member id has been fetched, and we want to check that the owner is making
// the request for it.
func AuthorizeID(w http.ResponseWriter, r *http.Request, mid int) bool {

	p := _json.Payload{}

	// Get token string from header
	a := r.Header.Get("Authorization")
	t, err := j_.FromHeader(a)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return false
	}

	// Create an AuthToken value from the token string
	at, err := j_.Check(t)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return false
	}

	// Now check the member ID passed in matches the token id
	if mid != at.Claims.ID {
		return false
	}

	return true
}
