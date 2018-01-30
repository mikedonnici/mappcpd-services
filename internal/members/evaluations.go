package members

import (
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// Evaluation is a struct representing an 'ep' or evaluation period. This is
// a defined but arbitrary period of time over which the CPD activities are reported
// on. The Evaluation definition is represented here, and MemberEvaluation used to
// represent an instance of an Evaluation belonging to a Member.
type Evaluation struct {
}

// MemberEvaluation represents an Evaluation belonging to a Member
type MemberEvaluation struct {
	ID             int                  `json:"id" bson:"id"`
	MemberID       int                  `json:"memberId" bson:"memberId"`
	Name           string               `json:"name" bson:"name"`
	StartDate      string               `json:"startDate" bson:"startDate"`
	EndDate        string               `json:"endDate" bson:"endDate"`
	Closed         bool                 `json:"closed"`
	CreditRequired int                  `json:"creditRequired" bson:"creditRequired"`
	CreditObtained int                  `json:"creditObtained" bson:"creditObtained"`
	Activities     []EvaluationActivity `json:"activities" bson:"activities"`
}

// EvaluationActivity represents a summary of a specific activity type
// that was recorded within an evaluation period
type EvaluationActivity struct {
	Activity string  `json:"activity" bson:"activity"`
	Total    float64 `json:"total" bson:"total"`
	Cap      float64 `json:"cap" bson:"cap"`
	Credit   float64 `json:"credit" bson:"credit"`
}

// EvaluationsByMemberID fetches all evaluation records for a member
// Received a member id, return a []MemberEvaluation
func EvaluationsByMemberID(id int) ([]MemberEvaluation, error) {

	es := []MemberEvaluation{}

	query := `SELECT cme.id, cme.member_id, ce.name,
	cme.cpd_points_required, cme.start_on, cme.end_on, cme.closed
	FROM ce_m_evaluation cme
	LEFT JOIN ce_evaluation ce ON cme.ce_evaluation_id = ce.id
	WHERE member_id = ?`

	rows, err := datastore.MySQL.Session.Query(query, id)
	if err != nil {
		return es, err
	}
	defer rows.Close()

	for rows.Next() {
		e := MemberEvaluation{}
		rows.Scan(
			&e.ID,
			&e.MemberID,
			&e.Name,
			&e.CreditRequired,
			&e.StartDate,
			&e.EndDate,
			&e.Closed,
		)

		// Evaluate activities for this evaluation period
		err := e.evaluate()
		if err != nil {
			return es, err
		}

		es = append(es, e)
	}

	return es, nil
}

// CurrentEvaluation returns a value with fields describing the current evaluation period
func CurrentEvaluation(memberID int) (MemberEvaluation, error) {
	var me MemberEvaluation

	// find the current one
	xme, err := EvaluationsByMemberID(memberID)
	if err != nil {
		return me, err
	}

	for _, v := range xme {
		if v.Closed == false {
			me = v
		}
	}

	return me, nil
}

// evaluate adds the activities to the MemberEvaluation value
// including total activity by types, caps and credit allowed.
func (e *MemberEvaluation) evaluate() error {

	// Gather activities by types, between the start and end dates...
	query := `SELECT ca.name,
			 SUM(cma.quantity),
			 cma.points_per_unit,
    		  	 SUM(cma.quantity * cma.points_per_unit)
	          FROM ce_m_activity cma
        	  LEFT JOIN ce_activity ca ON cma.ce_activity_id = ca.id
		  WHERE cma.active = 1
      		  AND cma.activity_on >= ?
                  AND cma.activity_on <= ?
                  AND cma.member_id = ?
                  GROUP BY cma.ce_activity_id`

	rows, err := datastore.MySQL.Session.Query(query, e.StartDate, e.EndDate, e.MemberID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {

		a := EvaluationActivity{}
		rows.Scan(
			&a.Activity,
			&a.Total,
			&a.Cap,
			&a.Credit,
		)
		e.Activities = append(e.Activities, a)
	}

	// Work out total credit for this evaluation (period)
	e.creditObtained()

	return nil
}

// creditObtained sets the .CreditObtained value by adding up all of the credit
// for each activity type within the evaluation
func (e *MemberEvaluation) creditObtained() error {

	var c float64
	for _, v := range e.Activities {
		c += v.Credit
	}
	// store it is an int
	e.CreditObtained = int(c)

	return nil
}
