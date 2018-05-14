package jwt_test

import (
	"fmt"
	"reflect"
	"time"

	//"fmt"
	//"reflect"
	"strings"
	"testing"

	"github.com/mappcpd/web-services/internal/platform/jwt"
	"github.com/mappcpd/web-services/testdata"
)

const issuer = "TestTokenIssuer"
const signingKey = "testTokenSigningKey"
const ttlHours = 4
const userID = 1
const userName = "Mike Donnici"
const userRole = "Member"

var helper = testdata.NewHelper()

func TestCreateTokenEmptyIssuer(t *testing.T) {
	var gotErr bool
	_, err := jwt.New("", signingKey, ttlHours).Encode()
	if err != nil {
		t.Log(err)
		gotErr = true
	}
	msg := "Should get an error when initialising a token with an empty issuer"
	helper.MessageResult(t, msg, true, gotErr)
}

func TestCreateTokenEmptySigningKey(t *testing.T) {
	var gotErr bool
	_, err := jwt.New(issuer, "", ttlHours).Encode()
	if err != nil {
		t.Log(err)
		gotErr = true
	}
	msg := "Should get an error when initialising a token with an empty signing key"
	helper.MessageResult(t, msg, true, gotErr)
}

func TestCreateTokenEmptyTTLHours(t *testing.T) {
	var gotErr bool
	_, err := jwt.New(issuer, signingKey, 0).Encode()
	if err != nil {
		t.Log(err)
		gotErr = true
	}
	msg := "Should get an error when initialising a token with zero TTL hours"
	helper.MessageResult(t, msg, true, gotErr)
}

func TestCreateToken(t *testing.T) {
	var gotErr bool
	_, err := jwt.New(issuer, signingKey, ttlHours).Encode()
	if err != nil {
		t.Log(err)
		gotErr = true
	}
	msg := "Should NOT get an error when initialising a token with proper args"
	helper.MessageResult(t, msg, false, gotErr)
}

func TestEncodedTokenFormat(t *testing.T) {
	tk, _ := jwt.New(issuer, signingKey, ttlHours).Encode()
	msg := "Encoded token string should be in the format of aaa.bbb.ccc"
	l := len(strings.Split(tk.Encoded, "."))
	helper.MessageResult(t, msg, 3, l)
}

func TestTTL(t *testing.T) {
	tk, _ := jwt.New(issuer, signingKey, ttlHours).Encode()
	msg := fmt.Sprintf("ExpiresAt should be %v hours in the future", ttlHours)
	unixHours := time.Now().Unix() / 3600
	expHours := tk.Claims.ExpiresAt / 3600
	helper.MessageResult(t, msg, 4, int(expHours-unixHours))
}

func TestInvalidToken(t *testing.T) {
	tk1, _ := jwt.New(issuer, signingKey, ttlHours).Encode()
	msg := "First token should be valid"
	helper.MessageResult(t, msg, true, tk1.Valid())

	// Second token with a different key
	tk2, _ := jwt.New(issuer, signingKey+"x", ttlHours).Encode()
	msg = "Second token should be valid"
	helper.MessageResult(t, msg, true, tk2.Valid())

	// Switch token strings and both should be invalid
	tk1.Encoded, tk2.Encoded = tk2.Encoded, tk1.Encoded
	msg = "First token now has incorrect signing key and should be INVALID"
	helper.MessageResult(t, msg, false, tk1.Valid())
	msg = "Second token now has incorrect signing key and should be INVALID"
	helper.MessageResult(t, msg, false, tk2.Valid())
}

func TestDecodeTokenString(t *testing.T) {

	// custom claims
	c := map[string]interface{}{
		"id":   userID,
		"name": userName,
		"role": userRole,
	}

	tk1, err := jwt.New(issuer, signingKey, ttlHours).CustomClaims(c).Encode()
	if err != nil {
		t.Fatalf("Error creating token: %s", err)
	}

	msg := "Initial token should be valid"
	helper.MessageResult(t, msg, true, tk1.Valid())

	tk2, err := jwt.Decode(tk1.Encoded, signingKey)
	if err != nil {
		t.Fatalf("Token string could not be decoded - %s", err)
	}

	deq := reflect.DeepEqual(tk1, tk2)
	msg = "Token string decoded into a Token value should deeply equal the initial Token value"
	helper.MessageResult(t, msg, true, deq)

}
