package handlers

import (
	"io"
	"fmt"
	"os"
	"log"

	"io/ioutil"
	"database/sql"
	"net/http"
	"github.com/mappcpd/web-services/cmd/webd/rest/router/handlers/responder"
	"github.com/mappcpd/web-services/cmd/webd/rest/router/middleware"
	"github.com/mappcpd/web-services/internal/member"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/member/activity"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/34South/envr"
	"encoding/base64"
)

func init() {
	envr.New("testEmail", []string{"SENDGRID_API_KEY"}).Clean()
}

// MembersProfile fetches a member record by id
func MembersProfile(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Get user id from token
	id := middleware.UserAuthToken.Claims.ID

	// Get the Member record
	m, err := member.MemberByID(id)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
		p.Data = m

		// TODO: remove this when fetching - should only be on update
		member.SyncMember(m)
	}

	p.Send(w)
}

// MembersActivities fetches activity records for a member
func MembersActivities(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	a, err := activity.MemberActivitiesByMemberID(middleware.UserAuthToken.Claims.ID)

	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
	p.Meta = map[string]int{"count": len(a)}
	p.Data = a
	p.Send(w)
}

// MembersEvaluation created reports for each evaluation period
// by gathering the CPD activities within the dates, adding them up, applying caps etc
func MembersEvaluation(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	// Collect the evaluation periods
	es, err := activity.MemberActivityReports(middleware.UserAuthToken.Claims.ID)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
	p.Meta = map[string]int{"count": len(es)}
	p.Data = es
	p.Send(w)
}

// MembersReports
func MembersReports(w http.ResponseWriter, r *http.Request) {

	p := responder.New(middleware.UserAuthToken.Token)

	ce, err := activity.CurrentMemberActivityReport(middleware.UserAuthToken.Claims.ID)

	// The PDF file is written to the PipeWriter (pw) by PDFReport and can then be read
	// from PipeReader (pr). Then we need to decide what we do with it!
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		activity.PDFReport(ce, pw)
	}()

	xb, err := ioutil.ReadAll(pr)
	if err != nil {
		fmt.Println(err)
	}

	if len(xb) > 0 {
		fmt.Println("ok")
	}

	// here we write it to a file, howerver we can just pass 'w' to PDf function and we could
	// send it straight to the requester!
	// In this case we will trigger a job that emails toe report, and that job will, in turn, return the report as a
	// slice of bytes that we can do whatever we want with.
	//ioutil.WriteFile("buffReport.pdf", xb, 0666)

	from := mail.NewEmail("Example User", "test@example.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Michael Donnici", "michael@mesa.net.au")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"

	a := mail.NewAttachment()
	encoded := base64.StdEncoding.EncodeToString(xb)
	a.SetContent(encoded)
	a.SetType("application/pdf")
	a.SetFilename("report.pdf")
	a.SetDisposition("attachment")
	a.SetContentID("CPD Report")

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	message.AddAttachment(a)


	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = responder.Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = responder.Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = responder.Message{http.StatusOK, "success", "Data retrieved from " + datastore.MySQL.Source}
	p.Data = ce
	p.Send(w)
}
