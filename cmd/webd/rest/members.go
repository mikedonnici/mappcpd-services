package rest

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mappcpd/web-services/internal/cpd"
	"github.com/mappcpd/web-services/internal/member"
	"github.com/mappcpd/web-services/internal/platform/email"
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
		p.Data = m

		// TODO: remove this when fetching - should only be on update
		member.SyncMember(DS, m)
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
