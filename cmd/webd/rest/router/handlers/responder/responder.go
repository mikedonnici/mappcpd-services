package responder

import (
	"encoding/json"
	"net/http"

	"github.com/mappcpd/web-services/internal/auth"
	"github.com/mappcpd/web-services/internal/platform/jwt"
)

// Payload represents a standard JSON format for ALL responses
// Message - the "header" part of the response (see below)
// Token - wherever possible, return a fresh token
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
	Count      int                    `json:"count" bson:"count"`
	Query      map[string]interface{} `json:"query" bson:"query"`
	Projection map[string]interface{} `json:"projection" bson:"projection"`
}

// MongoMeta stores the meta information from a query to MongoDB
type MongoMeta struct {
	Count int         `json:"count" bson:"count"`
	Query interface{} `json:"query" bson:"query"`
}

// New returns a pointer to a new Payload value. It received a token string which it will
// check for validity and if ok will set a refresh token as part of the payload.
func New(ts string) *Payload {

	p := Payload{}

	// if universal UserAuthToken value is present, use this to set fresh token
	t, err := jwt.Check(ts)
	if err != nil {
		// No panic here, we'll just not do a fresh token
		return &p
	}

	// otherwise, set p.Token to a FRESH token for either member or admin,
	// based on the scope of the current token
	//if t.CheckScope("member") {
	//	fmt.Println("Member token")
	//}
	//if t.CheckScope("admin") {
	//	fmt.Println("Admin token")
	//}

	// Fresh token - re-check the Scope from db rather than copying it from the current
	// token - in case permissions have been changed
	scopes, err := auth.AuthScope(t.Claims.ID)
	if err != nil {
		return &p
	}

	nt, err := jwt.CreateJWT(t.Claims.ID, t.Claims.Name, scopes)
	if err != nil {
		return &p
	}

	p.Token = nt.Token

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
