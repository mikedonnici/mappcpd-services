package rest

import (
	"net/http"

	reports "github.com/mikedonnici/mappcpd-services/internal/reports"
)

// ReportsTest handles a request to test the reports route
func ReportsTest(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)
	p.Message = Message{http.StatusOK, "success", "Request to reports test handler successful!"}
	p.Send(w)
}

// ReportsModulesByDate fetches data on modules by year-month
func ReportsModulesByDate(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	report, err := reports.ReportModulesByDate(DS)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Data = report
	m := make(map[string]interface{})
	m["count"] = len(report)
	m["description"] = "Report shows number of modules started by year-month"
	p.Meta = m
	p.Send(w)
}

// ReportsPointsByRecordDate fetches data on cpd activity (points) recorded by year-month
// according to WHEN they were recoded - so it is a measure of system activity. Actual activity
// dates are reported by ReportsPointsByActivityDate
func ReportsPointsByRecordDate(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	report, err := reports.ReportPointsByRecordDate(DS)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Data = report
	m := make(map[string]interface{})
	m["count"] = len(report)
	m["description"] = "Report groups CPD points by date of record creation - indicates system activity"
	p.Meta = m
	p.Send(w)
}

// ReportsPointsByActivityDate fetches data showing the cpd activity (points)
// according to the date of the activity itself - that is CPD Activity as opposed to system activity (above)
func ReportsPointsByActivityDate(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	report, err := reports.ReportPointsByActivityDate(DS)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = Message{http.StatusOK, "success", "Data retrieved from ???"}
	p.Data = report
	m := make(map[string]interface{})
	m["count"] = len(report)
	m["description"] = "Report groups CPD points by date of CPD activity"
	p.Meta = m
	p.Send(w)
}
