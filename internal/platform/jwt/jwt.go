package jwt

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// tokenLifeHours specifies the expiry time of the JWT, specified in env
var tokenLifeHours int

// The key for signing the JWTs
var signingKey = []byte("sg_e1JIskjvosMow6Vjra7N-oFep-vcNUH-2H1bVtDc")

type Token struct {
	issuer     string
	signingKey string
	ttlHours   int
	Encoded    string      `json:"token"`
	IssuedAt   time.Time   `json:"issuedAt"`
	ExpiresAt  time.Time   `json:"expiresAt"`
	Claims     TokenClaims `json:"claims"`
}

type TokenClaims struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Scope []string `json:"scope"`
	jwt.StandardClaims
}

//func init() {
//	var err error
//	envr.New("jwtEnv", []string{"JWT_TTL_HOURS"}).Passive()
//	tokenLifeHours, err = strconv.Atoi(os.Getenv("JWT_TTL_HOURS"))
//	if err != nil {
//		fmt.Println("Error setting tokenLifeHours from env var JWT_TTL_HOURS -", err)
//		fmt.Println("Setting a default value of 48 hours")
//		tokenLifeHours = 48
//	}
//}

// New returns a pointer to a Token
func New(issuer, signingKey string, ttlHours int) (*Token, error) {

	var t Token

	if issuer == "" {
		return &t, errors.New("Issuer cannot be blank")
	}
	if signingKey == "" {
		return &t, errors.New("Signing key cannot be blank")
	}
	if ttlHours < 1 {
		return &t, errors.New("TTL hours must be a positive integer")
	}

	t.issuer = issuer
	t.signingKey = signingKey
	t.ttlHours = ttlHours
	t.setStandardClaims()

	return &t, nil
}

//// NewToken is a convenience function that provides a one-liner to generate a token from id, name and scope
//func NewToken(id int, name string, scope []string) (Token, error) {
//	t := New()
//	t.Claims.ID = id
//	t.Claims.Name = name
//	t.Claims.Scope = scope
//	t.setStandardClaims()
//	err := t.Encode()
//	return *t, err
//}

func (t *Token) Encode() error {
	var err error
	t.Encoded, err = jwt.NewWithClaims(jwt.SigningMethodHS256, t.Claims).SignedString(signingKey)
	return err
}

// setStandardClaims sets the issuedAt and expiresAt fields from the exp and iat timestamps in the token claims.
func (t *Token) setStandardClaims() {

	t.Claims.StandardClaims = jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(t.ttlHours)).Unix(),
		Issuer:    t.issuer,
	}

	// Sets Unix dates at root of struct - for convenience??
	issuedAt := int64(t.Claims.StandardClaims.IssuedAt)
	t.IssuedAt = time.Unix(issuedAt, 0)

	expiresAt := int64(t.Claims.StandardClaims.ExpiresAt)
	t.ExpiresAt = time.Unix(expiresAt, 0)
}

// CheckScope checks a token has a particular scope string, Received token (t) and the string (s) to check for.
func (t *Token) CheckScope(s string) bool {
	for _, v := range t.Claims.Scope {
		if v == s {
			return true
		}
	}
	return false
}

func (t *Token) String() string {
	return string(t.Encoded)
}

// Check validates a JSON web token and returns an Encoded value
func Check(t string) (Token, error) {

	// token for return
	at := Token{Encoded: t}

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
		// the token in the JWT.setStandardClaims() method
		at.setStandardClaims()

		// set the values in  Encoded.Claims .. tricky
		// type problems here almost broke my brain...
		// These are TokenCLaims - ie custom claims
		at.Claims.ID = int(claims["id"].(float64))
		at.Claims.Name = claims["name"].(string)

		// Scope needs to be a []string but when we unpack the token
		// is is a []interface{} = so use assertion to make the []string
		s, ok := claims["scope"]
		if ok {
			a := s.([]interface{})
			b := make([]string, len(a))
			for i := range b {
				b[i] = a[i].(string)
			}
			at.Claims.Scope = b
		}

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
			TokenError = errors.New("Encoded malformed")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Encoded is either expired or not active yet
			TokenError = errors.New("Encoded expired or not active")
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
