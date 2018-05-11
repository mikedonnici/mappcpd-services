package member

import (
	"os"
	"strconv"

	"github.com/mappcpd/web-services/internal/auth"
	"github.com/mappcpd/web-services/internal/platform/jwt"
	"github.com/pkg/errors"
)

// RefreshToken will validate a JWT, and return a fresh token if all ok
func RefreshToken(currentToken string) (string, error) {

	// Validate current token
	ct, err := jwt.Check(currentToken)
	if err != nil {
		return "", errors.New("Cannot refresh token as current token is invalid: - " + err.Error())
	}

	// Make sure the current token has "member" scope to prevent switch from admin token
	if ct.CheckScope("member") == false {
		return "", errors.New("Cannot refresh non-member token")
	}

	// Verify scope directly from database in case scope has changed since the original token was issued
	scope, err := auth.AuthScope(ct.Claims.ID)
	if err != nil {
		return "", errors.New("Scope has changed since the original token was issued - " + err.Error())
	}

	t, err := freshToken(ct.Claims.ID, ct.Claims.Name, scope)
	if err != nil {
		return "", errors.New("Error creating new token -  " + err.Error())
	}

	return t.Encoded, nil
}

// todo move this to graphql server where it is used
// freshJWS is a copy of of the same unexported func in responder.go - the above func is only called
// in the GraphQL server and as this one here replies on env vars this needs to be removed at some point.
// freshToken issues a new token and adds custom claims id (member id) and name (member name) and well as custom scope
func freshToken(id int, name string, scope []string) (jwt.Token, error) {

	var t *jwt.Token

	iss := os.Getenv("MAPPCPD_API_URL")
	key := os.Getenv("MAPPCPD_JWT_SIGNING_KEY")
	ttl, err := strconv.Atoi(os.Getenv("MAPPCPD_JWT_TTL_HOURS"))
	if err != nil {
		return *t, errors.Wrap(err, "Could not create fresh token")
	}

	t, err = jwt.New(iss, key, ttl)
	if err != nil {
		return *t, errors.Wrap(err, "Could not create fresh token")
	}

	t.Claims.ID = id
	t.Claims.Name = name
	t.Claims.Scope = scope

	err = t.Encode()
	if err != nil {
		return *t, errors.Wrap(err, "Could not encode token")
	}

	return *t, nil
}
