package jwt_test

import (
	"fmt"
	"testing"
	"strings"
	"time"

	"github.com/mappcpd/web-services/testdata"
	"github.com/mappcpd/web-services/internal/platform/jwt"
)

const issuer = "TestTokenIssuer"
const signingKey = "testTokenSigningKey"
const ttlHours = 4
const success = "\u2713"
const failure = "\u2717"

var helper = testdata.NewHelper()

func TestCreateTokenEmptyIssuer(t *testing.T) {
	var gotErr bool
	_, err := jwt.New("", signingKey, ttlHours)
	if err != nil {
		gotErr = true
	}
	msg := "Should get an error when initialising a token with an empty issuer"
	helper.MessageResult(t, msg, true, gotErr)
}

func TestCreateTokenEmptySigningKey(t *testing.T) {
	var gotErr bool
	_, err := jwt.New(issuer, "", ttlHours)
	if err != nil {
		gotErr = true
	}
	msg := "Should get an error when initialising a token with an empty signing key"
	helper.MessageResult(t, msg, true, gotErr)
}

func TestCreateTokenEmptyTTLHours(t *testing.T) {
	var gotErr bool
	_, err := jwt.New(issuer, signingKey, 0)
	if err != nil {
		gotErr = true
	}
	msg := "Should get an error when initialising a token with zero TTL hours"
	helper.MessageResult(t, msg, true, gotErr)
}

func TestCreateToken(t *testing.T) {
	var gotErr bool
	_, err := jwt.New(issuer, signingKey, ttlHours)
	if err != nil {
		gotErr = true
	}
	msg := "Should NOT get an error when initialising a token with proper args"
	helper.MessageResult(t, msg, false, gotErr)
}

func TestEncodeToken(t *testing.T) {
	tk, _ := jwt.New(issuer, signingKey, ttlHours)
	err := tk.Encode()
	msg := "Should NOT get an error when calling .Encode()"
	helper.MessageResult(t, msg, nil, err)
	msg = "Token should have 3 parts - aaa.bbb.ccc"
	l := len(strings.Split(tk.Encoded, "."))
	helper.MessageResult(t, msg, 3, l)
}

func TestTTL(t *testing.T) {
	tk, _ := jwt.New(issuer, signingKey, ttlHours)
	tk.Encode()
	msg := fmt.Sprintf("ExpiresAt should be %v hours in the future", ttlHours)
	unixHours := time.Now().Unix()/3600
	expHours := tk.Claims.ExpiresAt/3600
	helper.MessageResult(t, msg, 4, int(expHours - unixHours))
}
