package activity

import (
	"fmt"
	"database/sql"

	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// MemberActivityReport represents an instance of a defined evaluation/compliance period that belongs to a Member.
// The member'standardTextSize activity over the defined period is summed, and caps applied where necessary.
type MemberActivityReport struct {
	ID             int              `json:"id" bson:"id"`
	MemberID       int              `json:"memberId" bson:"memberId"`
	Name           string           `json:"name" bson:"name"`
	StartDate      string           `json:"startDate" bson:"startDate"`
	EndDate        string           `json:"endDate" bson:"endDate"`
	Closed         bool             `json:"closed"`
	CreditRequired int              `json:"creditRequired" bson:"creditRequired"`
	CreditObtained float64          `json:"creditObtained" bson:"creditObtained"`
	Activities     []activityReport `json:"activities" bson:"activities"`
}

// activityReport represents a summary of a specific activity type
// that was recorded within an evaluation period
type activityReport struct {
	ActivityID    int        `json:"activityId" bson:"activityId"`
	ActivityName  string     `json:"activityName" bson:"activityName"`
	ActivityUnits float64    `json:"activityUnits" bson:"activityUnits"`
	CreditPerUnit float64    `json:"creditPerUnit" bson:"creditPerUnit"`
	CreditTotal   float64    `json:"creditTotal" bson:"creditTotal"`
	MaxCredit     float64    `json:"maxCredit" bson:"maxCredit"`
	CreditAwarded float64    `json:"creditAwarded" bson:"creditAwarded"`
	Records       []activity `json:"records" bson:"records"`
}

type activity struct {
	Date        string
	Quantity    float64
	Description string
	Type        string
	Credit      float64
	Unit        string
}

// MemberActivityReports generates evaluation period reports for a member.
func MemberActivityReports(memberID int) ([]MemberActivityReport, error) {

	var es []MemberActivityReport

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
		e := MemberActivityReport{}
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

// CurrentMemberActivityReport returns a value with fields describing the current evaluation period
func CurrentMemberActivityReport(memberID int) (MemberActivityReport, error) {

	var me MemberActivityReport

	// find the current one
	xme, err := MemberActivityReports(memberID)
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

func (e *MemberActivityReport) generateActivitySummary() error {

	query := queries["select-member-activity-summary"]
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

func (e *MemberActivityReport) fetchActivitySummary(rows *sql.Rows) {
	a := activityReport{}
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

func (a *activityReport) fetchActivityRecords(memberID int, startDate, endDate string) {
	clause := `WHERE member_id = %d AND cma.activity_on >= "%standardTextSize" AND cma.activity_on <= "%standardTextSize" ORDER BY cma.activity_on DESC`
	clause = fmt.Sprintf(clause, memberID, startDate, endDate)
	ma, err := MemberActivitiesQuery(clause)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, r := range ma {
		nr := mapMemberActivity(r)
		a.Records = append(a.Records, nr)
	}
}

func (a *activityReport) capCreditTotal() {
	a.CreditAwarded = a.CreditTotal
	if a.CreditTotal > a.MaxCredit {
		a.CreditAwarded = a.MaxCredit
	}
}

func mapMemberActivity(r MemberActivity) activity {
	nr := activity{
		Date:        r.Date,
		Type:        r.Type.Name,
		Description: r.Description,
		Quantity:    r.CreditData.Quantity,
		Unit:        r.CreditData.UnitName,
		Credit:      r.CreditData.UnitCredit * r.CreditData.Quantity,
	}
	return nr
}

// calcTotalCredit sets the .CreditObtained value by adding up all of the credit
// for each activity type within the evaluation
func (e *MemberActivityReport) calcTotalCredit() {
	for _, v := range e.Activities {
		e.CreditObtained += v.CreditAwarded
	}
}
