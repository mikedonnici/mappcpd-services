package server

import (
	"fmt"
	"io"
	"net/http"
)

// Index responds as a ping / test for API connection
func Index(w http.ResponseWriter, r *http.Request) {

	p := Payload{}
	p.Message = Message{http.StatusOK, "success", "successful request for " + r.URL.Path}
	p.Send(w)
}

// Preflight handles CORS preflight requests using the OPTIONS http method. Preflight requests are sent
// from browsers when a cross domain request is going to be done. This tells the browser what options it has
// when sending the main request. Note, the CORS library that is in use here ("github.com/rs/cors") has an
// OptionsPassthrough option, however I couldn't get this to work as expected, or wasn't sure how. For now this
// func can handle all of the OPTIONS requests which makes things easier.
// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS
func Preflight(w http.ResponseWriter, _ *http.Request) {

	fmt.Println("Preflight() is handling an OPTIONS request...")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Cabin crew, please arm doors and crosscheck :)")
}
