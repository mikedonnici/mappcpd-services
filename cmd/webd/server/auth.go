package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/cardiacsociety/web-services/internal/auth"
	"github.com/cardiacsociety/web-services/internal/platform/jwt"
)

// AuthMemberLogin handles a authenticates a user by login and password, against
// the db. Scope can also be passed in for admin access.
func AuthMemberLogin(w http.ResponseWriter, r *http.Request) {

	// create a binding struct for the JSON request body
	// ie. this is what we are expecting -CAPS for field names!!!
	type Auth struct {
		Login    string   `json:"login"`
		Password string   `json:"password"`
		Scope    []string `json:"scope"`
	}
	a := Auth{}

	// Response
	p := Payload{}

	// Pull the JSON body out of the request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&a)
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failure", err.Error()}
		p.Send(w)
		return
	}

	// AuthMember returns ID and Name which we pass to the token generator
	id, name, err := auth.AuthMember(DS, a.Login, a.Password)
	if err != nil {
		msg := err.Error()
		if err == sql.ErrNoRows {
			msg = "Login failed"
		}
		p.Message = Message{http.StatusUnauthorized, "failure", msg}
		p.Send(w)
		return
	}

	at, err := freshToken(id, name, "member")
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failure", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = Message{http.StatusOK, "success", "Authentication successful!"}
	p.Data = at
	p.Send(w)
}

// AuthMemberCheckHandler handles a GET request that will verify the JSON Web Encoded
func AuthMemberCheckHandler(w http.ResponseWriter, r *http.Request) {

	p := Payload{}

	// Get the token from the auth header, 'Bearer' seems useless but this is an OAuth2 standard
	// Authorization: Bearer [jwt]
	a := r.Header.Get("Authorization")
	t, err := jwt.FromHeader(a)
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failure", err.Error()}
		p.Send(w)
		return
	}

	jt, err := jwt.Decode(t, os.Getenv("MAPPCPD_JWT_SIGNING_KEY"))
	if err != nil {
		p.Message = Message{http.StatusUnauthorized, "failure", "Authorization failed: " + err.Error()}
		p.Send(w)
		return
	}

	p.Message = Message{http.StatusOK, "success", "Authorized: token is valid"}
	p.Data = jt
	p.Send(w)
}

// MembersToken handles a GET request which validates the current token
// and issue a fresh one, so the consumer can update it at their end
func MembersToken(w http.ResponseWriter, r *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Get the token from the auth header, 'Bearer' seems useless but this is an OAuth2 standard
	// Authorization: Bearer [jwt]
	a := r.Header.Get("Authorization")
	t, err := jwt.FromHeader(a)
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failure", err.Error()}
		p.Send(w)
		return
	}

	// Decode current token first
	at, err := jwt.Decode(t, os.Getenv("MAPPCPD_JWT_SIGNING_KEY"))
	if err != nil {
		p.Message = Message{http.StatusUnauthorized, "failure", "Cannot refresh token as current token is invalid: " + err.Error()}
		p.Send(w)
		return
	}

	// Make sure the current token has "member" scope to prevent switch from admin token
	if at.Claims.Role != "member" {
		p.Message = Message{http.StatusUnauthorized, "failure", "Cannot refresh non-member token"}
		p.Send(w)
		return
	}

	nt, err := freshToken(at.Claims.ID, at.Claims.Name, at.Claims.Role)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failure", err.Error()}
		p.Send(w)
		return
	}

	// All clear
	p.Message = Message{http.StatusOK, "success", "Current token is valid, fresh token supplied in data.new.token"}

	// Data payload will contain the current and a fresh token
	data := make(map[string]jwt.Token)
	data["current"] = at
	data["new"] = nt
	p.Data = data
	p.Send(w)
}

// AuthAdminLogin handles a authenticates an admin user by login and password, against
// the db. Requires an explicit 'scope' property requesting admin access.
func AuthAdminLogin(w http.ResponseWriter, r *http.Request) {

	// create a binding struct for the JSON request body
	// ie. this is what we are expecting -CAPS for field names!!!
	type Auth struct {
		Login    string   `json:"login"`
		Password string   `json:"password"`
		Scope    []string `json:"scope"`
	}
	a := Auth{}

	// Response
	p := Payload{}

	// Pull the JSON body out of the request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&a)
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failure", err.Error()}
		p.Send(w)
		return
	}

	// PostAdminAuth returns ID and Name which we pass to the token generator
	id, name, err := auth.AdminAuth(DS, a.Login, a.Password)
	if err != nil {
		msg := err.Error()
		if err == sql.ErrNoRows {
			msg = "Login failed"
		}
		p.Message = Message{http.StatusUnauthorized, "failure", msg}
		p.Send(w)
		return
	}

	at, err := freshToken(id, name, "admin")
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failure", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = Message{http.StatusOK, "success", "Authentication successful!"}
	p.Data = at
	p.Send(w)
}

// AuthAdminRefreshHandler handles a GET request which validates the current token
// and issues a fresh one so the consumer can extend validity to the maximum time.
// The only difference between this func and GetAuthRefresh is the function call to set scope claims.
func AuthAdminRefreshHandler(w http.ResponseWriter, r *http.Request) {

	p := Payload{}

	// Get the token from the auth header, 'Bearer' seems useless but this is an OAuth2 standard
	// Authorization: Bearer [jwt]
	a := r.Header.Get("Authorization")
	t, err := jwt.FromHeader(a)
	if err != nil {
		p.Message = Message{http.StatusBadRequest, "failure", err.Error()}
		p.Send(w)
		return
	}

	// Decode current token first
	at, err := jwt.Decode(t, os.Getenv("MAPPCPD_JWT_SIGNING_KEY"))
	if err != nil {
		p.Message = Message{http.StatusUnauthorized, "failure", "Cannot refresh token as current token is invalid: " + err.Error()}
		p.Send(w)
		return
	}

	// Make sure the current token has admin scope to prevent a normal user token upgrading to admin!
	if at.Claims.Role != "admin" {
		p.Message = Message{http.StatusUnauthorized, "failure", "Cannot refresh non-admin token"}
		p.Send(w)
		return
	}

	nt, err := freshToken(at.Claims.ID, at.Claims.Name, "admin")
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failure", err.Error()}
		p.Send(w)
		return
	}

	// All clear
	p.Message = Message{http.StatusOK, "success", "Current token is valid, fresh token supplied in data.new.token"}

	// Data payload will contain the current and a fresh token
	data := make(map[string]jwt.Token)
	data["current"] = at
	data["new"] = nt
	p.Data = data
	p.Send(w)
}
