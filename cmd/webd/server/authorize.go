package server

import (
	"net/http"
	"os"

	"github.com/cardiacsociety/web-services/internal/platform/jwt"
)

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

	// Create an Encoded value from the token string
	at, err := jwt.Decode(t, os.Getenv("MAPPCPD_JWT_SIGNING_KEY"))
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
