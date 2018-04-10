package members

import (
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"fmt"
	"database/sql"
)

// MemberEvaluation represents a reporting period for CPD activity, for a particular member - a defined period
// over which the credit for CPD activity is summed, and caps applied.
type MemberEvaluation struct {
	ID             int                  `json:"id" bson:"id"`
	MemberID       int                  `json:"memberId" bson:"memberId"`
	Name           string               `json:"name" bson:"name"`
	StartDate      string               `json:"startDate" bson:"startDate"`
	EndDate        string               `json:"endDate" bson:"endDate"`
	Closed         bool                 `json:"closed"`
	CreditRequired int                  `json:"creditRequired" bson:"creditRequired"`
	CreditObtained float64              `json:"creditObtained" bson:"creditObtained"`
	Activities     []EvaluationActivity `json:"activities" bson:"activities"`
}

// EvaluationActivity represents a summary of a specific activity type
// that was recorded within an evaluation period
type EvaluationActivity struct {
	ActivityID    int                      `json:"activityId" bson:"activityId"`
	ActivityName  string                   `json:"activityName" bson:"activityName"`
	ActivityUnits float64                  `json:"activityUnits" bson:"activityUnits"`
	CreditPerUnit float64                  `json:"creditPerUnit" bson:"creditPerUnit"`
	CreditTotal   float64                  `json:"creditTotal" bson:"creditTotal"`
	MaxCredit     float64                  `json:"maxCredit" bson:"maxCredit"`
	CreditAwarded float64                  `json:"creditAwarded" bson:"creditAwarded"`
	Records       []map[string]interface{} `json:"records" bson:"records"`
}

// EvaluationsByMemberID generates evaluation period reports for a member.
func EvaluationsByMemberID(memberID int) ([]MemberEvaluation, error) {

	es := []MemberEvaluation{}

	query := `SELECT cme.id, cme.member_id, ce.name,
	cme.cpd_points_required, cme.start_on, cme.end_on, cme.closed
	FROM ce_m_evaluation cme
	LEFT JOIN ce_evaluation ce ON cme.ce_evaluation_id = ce.id
	WHERE member_id = ?`

	rows, err := datastore.MySQL.Session.Query(query, memberID)
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

		err := e.generateActivitySummary()
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

func (e *MemberEvaluation) generateActivitySummary() error {

	query := `SELECT
					ca.id as ActivityID,
    				ca.name as ActivityName,
    				SUM(cma.quantity) as TotalUnits,
    				cma.points_per_unit as UnitCredit,
    				SUM(cma.quantity * cma.points_per_unit) as CreditObtained,
    				cma.annual_points_cap as CappedCredit
				FROM
    				ce_m_activity cma
        		LEFT JOIN
    				ce_activity ca ON cma.ce_activity_id = ca.id
				WHERE
    				cma.active = 1
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
		e.fetchActivitySummary(rows)
	}

	e.calcTotalCredit()

	return nil
}

func (e *MemberEvaluation) fetchActivitySummary(rows *sql.Rows) {
	a := EvaluationActivity{}
	rows.Scan(
		&a.ActivityID,
		&a.ActivityName,
		&a.ActivityUnits,
		&a.CreditPerUnit,
		&a.CreditTotal,
		&a.MaxCredit,
	)
	a.capCreditTotal()
	a.fetchActivityRecords(e.MemberID, e.StartDate, e.EndDate)
	e.Activities = append(e.Activities, a)
}

func (a *EvaluationActivity) fetchActivityRecords(memberID int, startDate, endDate string) {
	clause := `WHERE member_id = %d AND cma.activity_on >= "%s" AND cma.activity_on <= "%s" ORDER BY cma.activity_on DESC`
	clause = fmt.Sprintf(clause, memberID, startDate, endDate)
	ma, err := MemberActivitiesQuery(clause)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, r := range ma {
		nr := reducedMemberActivity(r)
		a.Records = append(a.Records, nr)
	}
}

func (a *EvaluationActivity) capCreditTotal() {
	a.CreditAwarded = a.CreditTotal
	if a.CreditTotal > a.MaxCredit {
		a.CreditAwarded = a.MaxCredit
	}
}

func reducedMemberActivity(r MemberActivity) map[string]interface{} {
	nr := map[string]interface{}{
		"activityDate":        r.Date,
		"activityType":        r.Type.Name,
		"activityDescription": r.Description,
		"activityQuantity":    r.CreditData.Quantity,
		"unit":                r.CreditData.UnitName,
		"activityCredit":      r.CreditData.UnitCredit * r.CreditData.Quantity,
	}
	return nr
}

// calcTotalCredit sets the .CreditObtained value by adding up all of the credit
// for each activity type within the evaluation
func (e *MemberEvaluation) calcTotalCredit() {
	for _, v := range e.Activities {
		e.CreditObtained += v.CreditAwarded
	}
}
