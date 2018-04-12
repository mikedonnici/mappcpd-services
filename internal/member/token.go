package member

import (
	"github.com/pkg/errors"

	"github.com/mappcpd/web-services/internal/auth"
	"github.com/mappcpd/web-services/internal/platform/jwt"
)

// FreshToken will validate a JWT, and return a fresh token if all ok
func FreshToken(currentToken string) (string, error) {

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
	scopes, err := auth.AuthScope(ct.Claims.ID)
	if err != nil {
		return "", errors.New("Scope has changed since the original token was issued - " + err.Error())
	}

	nt, err := jwt.CreateJWT(ct.Claims.ID, ct.Claims.Name, scopes)
	if err != nil {
		return "", errors.New("Error creating new token -  " + err.Error())
	}

	return nt.String(), nil
}
