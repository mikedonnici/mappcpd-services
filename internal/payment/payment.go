// Package payment provides access to membership payment data
package payment

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// Payment describes receipt of an amount of money
type Payment struct {
	ID          int              `json:"id" bson:"id"`
	CreatedAt   time.Time        `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt" bson:"updatedAt"`
	MemberID    int              `json:"memberId" bson:"memberId"`
	Member      string           `json:"member" bson:"member"`
	Date        time.Time        `json:"date" bson:"date"`
	Type        string           `json:"type" bson:"type"`
	Amount      float64          `json:"Amount" bson:"Amount"`
	Comment     string           `json:"comment" bson:"comment"`
	DataField1  string           `json:"dataField1" bson:"dataField1"`
	DataField2  string           `json:"dataField2" bson:"dataField2"`
	DataField3  string           `json:"dataField3" bson:"dataField3"`
	DataField4  string           `json:"dataField4" bson:"dataField4"`
	Allocations []InvoicePayment `json:"allocations" bson:"allocations"`
}

// InvoicePayment represents the allocation of part of all of the payment amount, to an invoice.
type InvoicePayment struct {
	InvoiceID int       `json:"invoiceId" bson:"paymentId"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Amount    float64   `json:"amount" bson:"amount"`
}

// ByID returns the Payment identified by paymentID, or an error if not found.
func ByID(ds datastore.Datastore, paymentID int) (Payment, error) {
	var p Payment
	q := fmt.Sprintf(queries["select-payment-by-id"], paymentID)
	xp, err := execute(ds, q)
	if err != nil {
		return p, err
	}
	if len(xp) == 0 {
		return p, sql.ErrNoRows
	}
	p = xp[0] // one result
	return p, nil
}

// ByIDs returns multiple Payment values identified by paymentIDs
func ByIDs(ds datastore.Datastore, paymentIDs []int) ([]Payment, error) {
	idList := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(paymentIDs)), ","), "[]")
	q := queries["select-payments"] + fmt.Sprintf(" AND p.id IN (%s)", idList)
	return execute(ds, q)
}

func execute(ds datastore.Datastore, query string) ([]Payment, error) {
	var xp []Payment

	rows, err := ds.MySQL.Session.Query(query)
	if err != nil {
		return xp, fmt.Errorf("Query() err = %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		p, err := scanRow(rows)
		if err != nil {
			return xp, err
		}

		err = rows.Err()
		if err != nil {
			return xp, err
		}

		invoicePayments, err := paymentAllocations(ds, p.ID)
		if err != nil {
			return xp, fmt.Errorf("paymentAllocations() err = %s", err)
		}
		p.Allocations = invoicePayments

		xp = append(xp, p)
	}

	return xp, nil
}

// scanRow scans the current row via *sql.Rows, to avoid duplicating the .Scan() for Query and QueryRow.
func scanRow(row *sql.Rows) (Payment, error) {

	var p Payment
	var createdAt, updatedAt, paymentDate string

	err := row.Scan(
		&p.ID,
		&createdAt,
		&updatedAt,
		&p.MemberID,
		&p.Member,
		&paymentDate,
		&p.Type,
		&p.Amount,
		&p.Comment,
		&p.DataField1,
		&p.DataField2,
		&p.DataField3,
		&p.DataField4,
	)
	if err != nil {
		return p, err
	}

	// convert date strings to time.Time
	p.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return p, err
	}
	p.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt)
	if err != nil {
		return p, err
	}
	p.Date, err = time.Parse("2006-01-02", paymentDate)
	if err != nil {
		return p, err
	}

	return p, nil
}

// paymentAllocations fetches the InvoicePayment values for the payment identified by paymentID. Part or all of the payment amount can be
// allocated to one or more invoices.
func paymentAllocations(ds datastore.Datastore, paymentID int) ([]InvoicePayment, error) {

	var result []InvoicePayment

	q := fmt.Sprintf(queries["select-payment-allocations"], paymentID)
	rows, err := ds.MySQL.Session.Query(q)
	if err != nil {
		return result, fmt.Errorf("Query() err = %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ip InvoicePayment
		var createdAt string
		err := rows.Scan(
			&ip.InvoiceID,
			&createdAt,
			&ip.Amount,
		)
		if err != nil {
			return result, fmt.Errorf("rows.Scan() err = %s", err)
		}

		ip.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return result, err
		}

		result = append(result, ip)
	}
	return result, nil
}
