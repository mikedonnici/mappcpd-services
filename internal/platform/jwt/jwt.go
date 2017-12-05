package jwt

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// todo - tokenLifeHours should be configurable via env var
const tokenLifeHours = 168 // 1 week

// The key for signing the JWTs - using the MYSQL_URL string for now so it will be host specific
var signingKey = []byte(os.Getenv("MAPPCPD_MYSQL_URL"))

type AuthToken struct {
	Token     string      `json:"token"`
	IssuedAt  time.Time   `json:"issuedAt"`
	ExpiresAt time.Time   `json:"expiresAt"`
	Claims    TokenClaims `json:"claims"`
}

type TokenClaims struct {
	ID    int64    `json:"id"`
	Name  string   `json:"name"`
	Scope []string `json:"scope"`
	jwt.StandardClaims
}

// CreateJWT creates a JWT
func CreateJWT(id int64, name string, scope []string) (AuthToken, error) {

	// Return token
	at := AuthToken{}

	// Create the Claims
	claims := TokenClaims{
		id,
		name,
		scope,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * tokenLifeHours).Unix(),
			Issuer:    os.Getenv("MAPPCPD_API_URL"),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := newToken.SignedString(signingKey)
	if err != nil {
		return at, err
	}

	at.Token = ss
	at.setDates()
	at.Claims = claims

	return at, nil
}

// Check validates a JSON web token and returns an AuthToken value
func Check(t string) (AuthToken, error) {

	// token for return
	at := AuthToken{Token: t}

	// Custom error
	var TokenError error

	// The jwt library panics if the jwt does not contain 3 '.'s
	// Assume because it splits the string at each period and gets and index
	// out of range if it does not end up with three pieces.
	if len(strings.Split(t, ".")) < 3 {
		TokenError = errors.New("JWT should have 3 parts in format aaaa.bbbb.cccc")
		return at, TokenError
	}

	// Parse the token which sets the Valid field
	// This code lifted from https://godoc.org/github.com/dgrijalva/jwt-go#Parse
	// not 100% sure what's going on :)
	tok, err := jwt.Parse(t, func(tok *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		TokenError = errors.New("Error parsing token: " + err.Error())
		return at, TokenError
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if ok && tok.Valid {

		// Set the Unix dates - YES this is redundant because we re-parse
		// the token in the JWT.setDates() method
		at.setDates()

		// set the values in  AuthToken.Claims .. tricky
		// type problems here almost broke my brain...
		// These are TokenCLaims - ie custom claims
		at.Claims.ID = int64(claims["id"].(float64))
		at.Claims.Name = claims["name"].(string)

		// Scope needs to be a []string but when we unpack the token
		// is is a []interface{} = so use assertion to make the []string
		a := claims["scope"].([]interface{})
		b := make([]string, len(a))
		for i := range b {
			b[i] = a[i].(string)
		}
		at.Claims.Scope = b

		// These are jwt.StandardClaims ...
		at.Claims.ExpiresAt = int64(claims["exp"].(float64))
		at.Claims.IssuedAt = int64(claims["iat"].(float64))
		at.Claims.Issuer = claims["iss"].(string)

		return at, nil
	}

	// ve is some fancy bit shifting thingo, too fancy for me so will
	// return some simple error strings to the caller
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			TokenError = errors.New("Token malformed")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			TokenError = errors.New("Token expired or not active")
		} else {
			TokenError = errors.New("jwt " + err.Error())
		}
	} else {
		TokenError = errors.New("Failed to validate token - " + err.Error())
	}

	return at, TokenError
}

// FromHeader extracts the jwt string from the header Authorization string (a).
// The Header should be in the format: Bearer aaaa.bbbb.cccc
func FromHeader(a string) (string, error) {

	var errMsg error

	t := strings.Fields(a) // splits on any amount of white space
	if len(t) < 2 || t[0] != "Bearer" {
		errMsg = errors.New("Authorization header should be: Bearer [jwt]")
		return "", errMsg
	}

	return strings.TrimSpace(t[1]), nil
}

// setDates is a utility method for the custom JWT type to set the
// IssuedAt and ExpiresAt fields in our custom JWT struct
func (t *AuthToken) setDates() {

	// Parse the Token in our custom type  to get a a a jwt.token value
	// from which to extract the dates we need...
	tok, err := jwt.Parse(t.Token, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		log.Printf("setDates() error parsing token: %s", err.Error())
		return
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if ok && tok.Valid {

		// The dates in Claims are stored as float64
		// we want friendly date strings so need int64 first!
		iat := int64(claims["iat"].(float64))
		exp := int64(claims["exp"].(float64))

		// Get local Time values in Unix format
		iatUnix := time.Unix(iat, 0)
		expUnix := time.Unix(exp, 0)

		// Then tweak the format? - leave for now...

		t.IssuedAt = iatUnix
		t.ExpiresAt = expUnix
	}
}

// CheckScope checks a token has a particular scope string, Received token (t) and
// the string (s) to check for.
func (t AuthToken) CheckScope(s string) bool {

	for _, v := range t.Claims.Scope {
		if v == s {
			return true
		}
	}

	return false
}

func (t AuthToken) String() string {
	return string(t.Token)
}
