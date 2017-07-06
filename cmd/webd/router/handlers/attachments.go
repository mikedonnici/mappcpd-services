package handlers

import (
	"fmt"

	"net/http"

	_json "github.com/mappcpd/web-services/cmd/webd/router/handlers/json"
	"github.com/mappcpd/web-services/cmd/webd/router/middleware"
	"github.com/mappcpd/web-services/internal/attachments"
)

// PutRequest issues a signed url to allow for an object to be PUT to Amazon S3
func S3PutRequest(w http.ResponseWriter, r *http.Request) {

	p := _json.NewPayload(middleware.UserAuthToken.Token)

	// Return the URL Query string params for the caller's convenience, and signedURL
	upload := struct {
		Key           string `json:"key"`
		Bucket        string `json:"bucket"`
		FileName      string `json:"fileName"`
		FileType      string `json:"fileType"`
		SignedRequest string `json:"signedRequest"`
	}{
		Key:      r.FormValue("key"),
		Bucket:   r.FormValue("bucket"),
		FileName: r.FormValue("filename"),
		FileType: r.FormValue("filetype"),
	}

	// Todo - check for missing bits?
	//if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
	//	msg := fmt.Sprintf("Could not decode json in request body - %s", err.Error())
	//	p.Message = _json.Message{http.StatusBadRequest, "failed", msg}
	//	p.Send(w)
	//	return
	//}

	url, err := attachments.S3PutRequest(upload.Key, upload.Bucket)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	fmt.Println(url)
	upload.SignedRequest = url

	p.Message = _json.Message{http.StatusOK, "success", "Signed request in data.signedRequest."}
	p.Data = upload
	p.Send(w)
}
