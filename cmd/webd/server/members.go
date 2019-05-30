package server

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/cardiacsociety/web-services/internal/cpd"
	"github.com/cardiacsociety/web-services/internal/member"
	"github.com/cardiacsociety/web-services/internal/notification"
	"github.com/cardiacsociety/web-services/internal/platform/email"
)

// MembersProfile fetches a member record by id
func MembersProfile(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Get user id from token
	id := UserAuthToken.Claims.ID

	// Get the Member record
	m, err := member.ByID(DS, id)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = Message{http.StatusNotFound, "failed", err.Error()}
	case err != nil:
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
	default:
		p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
		err := m.SyncUpdated(DS)
		if err != nil {
			p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		}
		p.Data = m
	}

	p.Send(w)
}

// MembersActivities fetches activity records for a member
func MembersActivities(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	a, err := cpd.ByMemberID(DS, UserAuthToken.Claims.ID)

	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Meta = map[string]int{"count": len(a)}
	p.Data = a
	p.Send(w)
}

// MembersEvaluation created reports for each evaluation period
// by gathering the CPD activities within the dates, adding them up, applying caps etc
func MembersEvaluation(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	// Collect the evaluation periods
	es, err := cpd.MemberActivityReports(DS, UserAuthToken.Claims.ID)
	// Response
	switch {
	case err == sql.ErrNoRows:
		p.Message = Message{http.StatusNotFound, "failed", err.Error()}
		p.Send(w)
		return
	case err != nil:
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Meta = map[string]int{"count": len(es)}
	p.Data = es
	p.Send(w)
}

// CurrentActivityReport
func CurrentActivityReport(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)
	reportData, err := cpd.CurrentEvaluationPeriodReport(DS, UserAuthToken.Claims.ID)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	msg := fmt.Sprintf("Data retrieved from ???")
	p.Message = Message{http.StatusOK, "success", msg}
	p.Data = reportData
	p.Send(w)
}

// EmailCurrentActivityReport
func EmailCurrentActivityReport(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)
	reportData, err := cpd.CurrentEvaluationPeriodReport(DS, UserAuthToken.Claims.ID)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// The PDF file is written to the PipeWriter (pw) by PDFReport and can then be read
	// from PipeReader (pr). Then we need to decide what we do with it!
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		cpd.PDFReport(reportData, pw)
	}()
	xb, err := ioutil.ReadAll(pr)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}
	reportAttachment := base64.StdEncoding.EncodeToString(xb)
	//ioutil.WriteFile("buffReport.pdf", xb, 0666)

	e := email.New()
	e.FromName = "MappCPD Report"
	e.FromEmail = "system@mappcpd.com"
	e.ToName = "Dr Mike Donnici"
	e.ToEmail = "michael@mesa.net.au"
	e.Subject = "Your CPD Report"
	e.HTMLContent = "Please find you report attached"
	e.PlainContent = "Please find you report attached"
	e.Attachments = []email.Attachment{
		{"application/pdf", "cpdReport.pdf", reportAttachment},
	}
	err = e.Send()
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	msg := fmt.Sprintf("Report has been created an emailed to %s.", e.ToEmail)
	p.Message = Message{http.StatusOK, "success", msg}
	p.Data = reportData
	p.Send(w)
}

// MemberSendNotification sends an email to the member identified in the token
func MemberSendNotification(w http.ResponseWriter, r *http.Request) {
	p := NewResponder(UserAuthToken.Encoded)

	// member record id in token
	mem, err := member.ByID(DS, UserAuthToken.Claims.ID)
	if err != nil {
		msg := fmt.Sprintf("Could not find member record with id %v", UserAuthToken.Claims.ID)
		p.Message = Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	var body struct {
		SenderName  string   `json:"senderName"`
		SenderEmail string   `json:"senderEmail"`
		Subject     string   `json:"subject"`
		HTML        string   `json:"html"`
		Text        string   `json:"text"`
		Attachments []string `json:"attachments"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		msg := fmt.Sprintf("Could not read request body - %s", err)
		p.Message = Message{http.StatusBadRequest, "failed", msg}
		p.Send(w)
		return
	}

	// map to notification.Email
	fullName := fmt.Sprintf("%s %s", mem.FirstName, mem.LastName)
	em := notification.Email{
		ToName:       fullName,
		ToEmail:      mem.Contact.EmailPrimary,
		FromName:     body.SenderName,
		FromEmail:    body.SenderEmail,
		Subject:      body.Subject,
		HTMLContent:  body.HTML,
		PlainContent: body.Text,
	}
	fmt.Println("Sending to", em.ToName, em.ToEmail)
	err = em.Send()
	if err != nil {
		msg := fmt.Sprintf("Could not sent to '%s' - %s", em.ToEmail, err)
		p.Message = Message{http.StatusInternalServerError, "failed", msg}
		p.Send(w)
		return
	}

	msg := fmt.Sprintf("Sent to '%s'", em.ToEmail)
	p.Message = Message{http.StatusAccepted, "success", msg}
	p.Send(w)
}
