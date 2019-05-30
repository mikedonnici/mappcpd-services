package server

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/cardiacsociety/web-services/internal/platform/jwt"
	"github.com/pkg/errors"
)

// Payload represents a standard JSON format for ALL responses
// Message - the "header" part of the response (see below)
// Encoded - wherever possible, return a fresh token
// Meta - information about the data payload such as count etc
// Data - the actual data being returned, single object or an array of objects
type Payload struct {
	Message
	Token string      `json:"token"`
	Meta  interface{} `json:"meta"`
	Data  interface{} `json:"data"`
}

// Message holds the basic response information - like a header.
// Status - is the http status code
// Result - "success" or "failure" in general terms
// Message - a comment about the operation, or an error message - anything useful
type Message struct {
	Status  int    `json:"status" bson:"status"`
	Result  string `json:"result" bson:"result"`
	Message string `json:"message" bson:"message"`
}

// DocMeta stores the meta information from a query to MongoDB
// todo replace DocMeta with MongoMeta
type DocMeta struct {
	Count int                    `json:"count" bson:"count"`
	Query map[string]interface{} `json:"query" bson:"query"`
}

// MongoMeta stores the meta information from a query to MongoDB
type MongoMeta struct {
	Count int         `json:"count" bson:"count"`
	Query interface{} `json:"query" bson:"query"`
}

// New returns a pointer to a new Payload value. It received a token string which it will
// check for validity and if ok will set a refresh token as part of the payload.
func NewResponder(ts string) *Payload {

	p := Payload{}

	// if universal UserAuthToken value is present, use this to set fresh token
	t, err := jwt.Decode(ts, os.Getenv("MAPPCPD_JWT_SIGNING_KEY"))
	if err != nil {
		// No panic here, we'll just not do a fresh token
		return &p
	}

	ft, err := freshToken(t.Claims.ID, t.Claims.Name, t.Claims.Role)
	if err != nil {
		return &p
	}
	p.Data = ft

	return &p
}

// Send will; send the payload back to the requester
func (p Payload) Send(w http.ResponseWriter) error {

	// set all the preflight responses here as well, and as we have OPTIONS passThrough in
	// CORS this single responder should handle all cases...
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(p.Message.Status) // the http status code is part of the payload message

	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		return err
	}

	return nil
}

// freshToken issues a new token and adds custom claims id (member id) and name (member name) and well as custom scope
func freshToken(id int, name string, role string) (jwt.Token, error) {

	var t jwt.Token

	iss := os.Getenv("MAPPCPD_API_URL")
	key := os.Getenv("MAPPCPD_JWT_SIGNING_KEY")
	ttl, err := strconv.Atoi(os.Getenv("MAPPCPD_JWT_TTL_HOURS"))
	if err != nil {
		return t, errors.Wrap(err, "Could not convert hours string to int")
	}

	c := map[string]interface{}{
		"id":   id,
		"name": name,
		"role": role,
	}

	return jwt.New(iss, key, ttl).CustomClaims(c).Encode()
}
