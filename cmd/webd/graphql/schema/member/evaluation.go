package member

import (
	"github.com/mappcpd/web-services/internal/member/activity"
)

// evaluationData representations the member evaluation data
type evaluationData struct {
	ID             int     `json:"id"`
	ReportName     string  `json:"name"`
	StartDate      string  `json:"startDate"`
	EndDate        string  `json:"endDate"`
	CreditRequired float64 `json:"creditRequired"`
	CreditObtained float64 `json:"creditObtained"`
	Closed         bool    `json:"closed"`
}

// evaluations fetches all evaluations member and maps to local evaluationData values.
func evaluations(memberID int) ([]evaluationData, error) {

	var xed []evaluationData

	xar, err := activity.MemberActivityReports(memberID)
	for _, ar := range xar {
		e := mapEvaluationData(ar)
		xed = append(xed, e)
	}

	return xed, err
}

// currentEvaluation fetches the current evaluation period report for a member
func currentEvaluation(memberID int) (evaluationData, error) {

	var ed evaluationData

	ce, err := activity.CurrentEvaluationPeriodReport(memberID)
	if err != nil {
		return ed, err
	}
	ed = mapEvaluationData(ce)

	return ed, nil
}

// mapEvaluationData maps am activity.MemberActivityReport to a local evaluationData value
func mapEvaluationData(ar activity.MemberActivityReport) evaluationData {

	var ed evaluationData

	ed.ID = ar.ID
	ed.ReportName = ar.ReportName
	ed.StartDate = ar.StartDate
	ed.EndDate = ar.EndDate
	ed.CreditRequired = float64(ar.CreditRequired)
	ed.CreditObtained = float64(ar.CreditObtained)
	ed.Closed = ar.Closed

	return ed
}
