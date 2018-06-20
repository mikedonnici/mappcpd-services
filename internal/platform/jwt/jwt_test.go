package jwt_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mikedonnici/mappcpd-services/internal/platform/jwt"
	"github.com/mikedonnici/mappcpd-services/testdata"
	"github.com/matryer/is"
)

const issuer = "TestTokenIssuer"
const signingKey = "testTokenSigningKey"
const ttlHours = 4
const userID = 1
const userName = "Mike Donnici"
const userRole = "Member"

var helper = testdata.NewHelper()

func TestCreateTokenEmptyIssuer(t *testing.T) {
	is := is.New(t)
	_, err := jwt.New("", signingKey, ttlHours).Encode()
	is.True(err != nil) // Should get an error when initialising a token with an empty issuer
}

func TestCreateTokenEmptySigningKey(t *testing.T) {
	is := is.New(t)
	_, err := jwt.New(issuer, "", ttlHours).Encode()
	is.True(err != nil) // Should get an error when initialising a token with an empty signing key
}

func TestCreateTokenEmptyTTLHours(t *testing.T) {
	is := is.New(t)
	_, err := jwt.New(issuer, signingKey, 0).Encode()
	t.Log(err)
	is.True(err != nil) // Should get an error when initialising a token with zero TTL hours
}

func TestCreateToken(t *testing.T) {
	is := is.New(t)
	_, err := jwt.New(issuer, signingKey, ttlHours).Encode()
	t.Log(err)
	is.True(err == nil) // Should NOT get an error when initialising a token with proper args
}

func TestEncodedTokenFormat(t *testing.T) {
	is := is.New(t)
	tk, _ := jwt.New(issuer, signingKey, ttlHours).Encode()
	l := len(strings.Split(tk.Encoded, "."))
	is.Equal(l, 3) // Encoded token string should be in the format of aaa.bbb.ccc
}

func TestTTL(t *testing.T) {
	is := is.New(t)
	tk, _ := jwt.New(issuer, signingKey, ttlHours).Encode()
	hoursToLive := (tk.Claims.ExpiresAt / 3600) - (time.Now().Unix() / 3600)
	is.Equal(hoursToLive, int64(ttlHours)) // Incorrect TTL hours
}

func TestInvalidToken(t *testing.T) {
	is := is.New(t)

	tk1, _ := jwt.New(issuer, signingKey, ttlHours).Encode()
	is.True(tk1.Valid() == true) // First token should be valid

	// Second token with a different key
	tk2, _ := jwt.New(issuer, signingKey+"x", ttlHours).Encode()
	is.True(tk1.Valid() == true) // Second token should be valid

	// Switch token strings and both should be invalid
	tk1.Encoded, tk2.Encoded = tk2.Encoded, tk1.Encoded
	is.True(tk1.Valid() == false) // First token now has incorrect signing key and should be INVALID
	is.True(tk2.Valid() == false) // Second token now has incorrect signing key and should be INVALID
}

func TestDecodeTokenString(t *testing.T) {
	is := is.New(t)

	c := map[string]interface{}{
		"id":   userID,
		"name": userName,
		"role": userRole,
	}

	tk1, err := jwt.New(issuer, signingKey, ttlHours).CustomClaims(c).Encode()
	is.NoErr(err)                // Error creating token
	is.True(tk1.Valid() == true) // Initial token should be valid

	tk2, err := jwt.Decode(tk1.Encoded, signingKey)
	is.NoErr(err)                        // Error decoding token
	is.True(reflect.DeepEqual(tk1, tk2)) // Token and decoded Token should be deeply equal
}

func TestDecodeError(t *testing.T) {
	is := is.New(t)

	s := "bungtoken.stillhas.threeparts"
	_, err := jwt.Decode(s, signingKey)
	is.True(err != nil) // Decoding a fake token with proper format should return an error

	s = "thisIsNotAProperToken"
	_, err = jwt.Decode(s, signingKey)
	is.True(err != nil) // Decoding a fake token WITHOUT proper format should return an error
}

func TestTokenWithFutureDate(t *testing.T) {
	is := is.New(t)

	c := map[string]interface{}{
		"id":   userID,
		"name": userName,
		"role": userRole,
	}

	iat := time.Now().Add(time.Hour * time.Duration(4)) // issue 4 hours in the future
	tk, err := jwt.New(issuer, signingKey, ttlHours).CustomClaims(c).SetTimes(iat).Encode()
	is.NoErr(err)                // Error creating token
	is.True(tk.Valid() == false) // Token with future iat should be INVALID

	expectExpireTime := ttlHours + 4
	expireTime := int(tk.Claims.ExpiresAt/3600) - int(time.Now().Unix()/3600)
	is.True(expectExpireTime == expireTime) // Incorrect expire time
}
