package handlers

import (
	"encoding/json"
	"net/http"

	_json "github.com/mappcpd/web-services/cmd/webd/router/handlers/json"
	"github.com/mappcpd/web-services/cmd/webd/router/middleware"
	"github.com/mappcpd/web-services/internal/attachments"
	//"github.com/mappcpd/web-services/internal/platform/datastore"
	"fmt"
)

// PutRequest issues a signed url to allow for an object to be PUT to Amazon S3
func S3PutRequest(w http.ResponseWriter, r *http.Request) {

	p := _json.NewPayload(middleware.UserAuthToken.Token)

	// The body should contain the bucket and the full path (key) including file name and ext.
	body := struct {
		key      string `json:"key"`
		bucket   string `json:"bucket"`
		fileName string `json:"fileName"`
		fileType string `json:"fileType"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		msg := fmt.Sprintf("Could not decode json in request body - %s", err.Error())
		p.Message = _json.Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	fmt.Println(body)
	return

	url, err := attachments.S3PutRequest(body.key, body.bucket)
	if err != nil {
		p.Message = _json.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	p.Message = _json.Message{http.StatusOK, "success", "Signed URL issued."}
	p.Data = url
	p.Send(w)
}
