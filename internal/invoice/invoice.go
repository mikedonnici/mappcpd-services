package invoice

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/cardiacsociety/web-services/internal/member"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// Invoice represents an invoice :)
type Invoice struct {
	ID             int           `json:"id" bson:"id"`
	CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt" bson:"updatedAt"`
	MemberID       int           `json:"memberId" bson:"memberId"`
	IssueDate      time.Time     `json:"issueDate" bson:"issueDate"`
	DueDate        time.Time     `json:"dueDate" bson:"dueDate"`
	SubscriptionID int           `json:"subscriptionID" bson:"subscriptionID"`
	Subscription   string        `json:"subscription" bson:"subscription"`
	Amount         float64       `json:"Amount" bson:"Amount"`
	Paid           bool          `json:"paid" bson:"paid"`
	Comment        string        `json:"comment" bson:"comment"`
	Member         member.Member `json:"member"`
}

// ByID fetches an invoice by invoice ID
func ByID(ds datastore.Datastore, invoiceID int) (Invoice, error) {
	var i Invoice
	q := fmt.Sprintf(queries["select-invoice-by-id"], invoiceID)
	xi, err := execute(ds, q)
	if err != nil {
		return i, err
	}
	if len(xi) == 0 {
		return i, sql.ErrNoRows
	}
	i = xi[0] // one result

	i.attachMember(ds)

	return i, nil
}

func (i *Invoice) attachMember(ds datastore.Datastore) {
	m, err := member.ByID(ds, i.MemberID)
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		i.Member = *m
	}
}

// ByIDs returns multiple Invoice values identified by invoiceIDs
func ByIDs(ds datastore.Datastore, invoiceIDs []int) ([]Invoice, error) {

	var xi []Invoice

	idList := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(invoiceIDs)), ","), "[]")
	q := queries["select-invoices"] + fmt.Sprintf(" AND i.id IN (%s)", idList)
	xi, err := execute(ds, q)
	if err != nil {
		return nil, err
	}

	for i := range xi {
		xi[i].attachMember(ds)
	}

	return xi, err
}

func execute(ds datastore.Datastore, query string) ([]Invoice, error) {

	var xi []Invoice

	rows, err := ds.MySQL.Session.Query(query)
	if err != nil {
		return xi, fmt.Errorf("Query() err = %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		p, err := scanRow(rows)
		if err != nil {
			return xi, err
		}

		err = rows.Err()
		if err != nil {
			return xi, err
		}

		xi = append(xi, p)
	}

	return xi, nil
}

// scanRow scans the current row via *sql.Rows
func scanRow(row *sql.Rows) (Invoice, error) {

	var i Invoice

	// values that will need some manipulation
	var createdAt, updatedAt string // data dates
	var issueDate, dueDate string   // invoice dates
	var paid int                    // 0,1 represents boolean in database

	err := row.Scan(
		&i.ID,
		&createdAt,
		&updatedAt,
		&i.MemberID,
		&issueDate,
		&dueDate,
		&i.SubscriptionID,
		&i.Subscription,
		&i.Amount,
		&paid,
		&i.Comment,
	)
	if err != nil {
		return i, err
	}

	// convert date strings to time.Time
	i.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return i, err
	}
	i.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt)
	if err != nil {
		return i, err
	}
	i.IssueDate, err = time.Parse("2006-01-02", issueDate)
	if err != nil {
		return i, err
	}
	i.DueDate, err = time.Parse("2006-01-02", dueDate)
	if err != nil {
		return i, err
	}

	// Paid bool is 0,1 in the database
	i.Paid = false
	if paid == 1 {
		i.Paid = true
	}

	return i, nil
}
