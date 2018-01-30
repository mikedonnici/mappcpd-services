package data

import (
	"github.com/mappcpd/web-services/internal/members"
)

// MemberEvaluation representations the member evaluation data
type MemberEvaluation struct {
	//ID          int       `json:"id"`
	Name           string  `json:"name"`
	StartDate      string  `json:"startDate"`
	EndDate        string  `json:"endDate"`
	CreditRequired float64 `json:"creditRequired"`
	CreditObtained float64 `json:"creditObtained"`
	Closed         bool    `json:"closed"`
}

// GetMemberEvaluations fetches evaluation data for a member.
func GetMemberEvaluations(memberID int) ([]MemberEvaluation, error) {

	var xme []MemberEvaluation

	// This returns a nested struct which is simplified below.
	xma, err := members.EvaluationsByMemberID(memberID)

	for _, v := range xma {
		e := MemberEvaluation{
			Name:           v.Name,
			StartDate:      v.StartDate,
			EndDate:        v.EndDate,
			CreditRequired: float64(v.CreditRequired),
			CreditObtained: float64(v.CreditObtained),
			Closed:         v.Closed,
		}
		xme = append(xme, e)
	}

	return xme, err
}

// GetCurrentEvaluation fetches the current evaluation period data for a member.
func GetCurrentEvaluation(memberID int) (MemberEvaluation, error) {

	var me MemberEvaluation

	// This returns a nested struct which is simplified below.
	ce, err := members.CurrentEvaluation(memberID)
	if err != nil {
		return me, err
	}

	me.Name = ce.Name
	me.StartDate = ce.StartDate
	me.EndDate = ce.EndDate
	me.CreditRequired = float64(ce.CreditRequired)
	me.CreditObtained = float64(ce.CreditObtained)
	me.Closed = ce.Closed

	return me, nil
}
