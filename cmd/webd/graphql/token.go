package graphql

import (
	"os"
	"strconv"

	"github.com/mikedonnici/mappcpd-services/internal/platform/jwt"
	"github.com/pkg/errors"
)

// freshToken will validate a JWT, and return a fresh token if all ok
func freshToken(currentToken string) (string, error) {

	iss := os.Getenv("MAPPCPD_API_URL")
	key := os.Getenv("MAPPCPD_JWT_SIGNING_KEY")
	ttl, err := strconv.Atoi(os.Getenv("MAPPCPD_JWT_TTL_HOURS"))
	if err != nil {
		return "", errors.Wrap(err, "Could not convert hours string to int")
	}

	// Validate current token
	ct, err := jwt.Decode(currentToken, key)
	if err != nil {
		return "", errors.New("Cannot refresh token as current token is invalid: - " + err.Error())
	}

	// Make sure the current token has "member" role
	if ct.Claims.Role != "member" {
		return "", errors.New("Cannot refresh non-member token")
	}

	// custom claims
	c := map[string]interface{}{
		"id":   ct.Claims.ID,
		"name": ct.Claims.Name,
		"role": ct.Claims.Role,
	}

	t, err := jwt.New(iss, key, ttl).CustomClaims(c).Encode()
	if err != nil {
		return "", errors.New("Error creating new token -  " + err.Error())
	}

	return t.Encoded, nil
}
