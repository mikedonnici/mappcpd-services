// Package application provides access to membership applications data, that is, applications to become a member
package application

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// Application describes an application for membership
type Application struct {
	ID          int       `json:"id" bson:"id"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
	MemberID    int       `json:"memberId" bson:"memberId"`
	Member      string    `json:"member" bson:"member"`
	NominatorID int       `json:"nominatorId" bson:"nominatorId"`
	Nominator   string    `json:"nominator" bson:"nominator"`
	SeconderID  int       `json:"seconderId" bson:"seconderId"`
	Seconder    string    `json:"seconder" bson:"seconder"`
	Date        time.Time `json:"date" bson:"date"`
	ForTitle    string    `json:"forTitle" bson:"forTitle"`
	ForTitleID  int       `json:"forTitleId" bson:"forTitleId"`
	Status      int       `json:"status" bson:"status"`
	Comment     string    `json:"comment" bson:"comment"`
}

// ByID fetches an application record by id. This returns an error if no result is found.
func ByID(ds datastore.Datastore, applicationID int) (Application, error) {
	var a Application
	q := fmt.Sprintf(queries["select-application-by-id"], applicationID)
	r, err := execute(ds, q)
	if err != nil {
		return a, err
	}
	if len(r) == 0 {
		return a, sql.ErrNoRows
	}
	a = r[0] // one result
	return a, nil
}

// ByIDs fetches a set of applications by IDs.
func ByIDs(ds datastore.Datastore, applicationIDs []int) ([]Application, error) {
	idList := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(applicationIDs)), ","), "[]")
	clause := fmt.Sprintf(" AND ma.id IN (%s)", idList)
	return Query(ds, clause)
}

// ByMemberID fetches application records by member id. This does not return an error if no results are found, only an empty slice.
func ByMemberID(ds datastore.Datastore, memberID int) ([]Application, error) {
	q := fmt.Sprintf(queries["select-applications-by-memberid"], memberID)
	return execute(ds, q)
}

// Query runs a select query with the given clause
func Query(ds datastore.Datastore, clause string) ([]Application, error) {
	q := fmt.Sprintf(queries["select-applications"]+" %s", clause)
	return execute(ds, q)
}

func execute(ds datastore.Datastore, query string) ([]Application, error) {
	var xa []Application

	rows, err := ds.MySQL.Session.Query(query)
	if err != nil {
		return xa, fmt.Errorf("Query() err = %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		a, err := scanRow(rows)
		if err != nil {
			return xa, err
		}

		err = rows.Err()
		if err != nil {
			return xa, err
		}

		xa = append(xa, a)
	}

	return xa, nil
}

// scanRow scans the current row via *sql.Rows. This avoids duplicating the .Scan() for Query and QueryRow as the latter returns
// *sql.Row and the former *sql.Rows, and there is no way to get the current row, ie *sql.Row from *sql.Rows. Annoying!
func scanRow(row *sql.Rows) (Application, error) {

	var a Application
	var createdAt, updatedAt, applicationDate string

	err := row.Scan(
		&a.ID,
		&createdAt,
		&updatedAt,
		&a.MemberID,
		&a.Member,
		&a.NominatorID,
		&a.Nominator,
		&a.SeconderID,
		&a.Seconder,
		&applicationDate,
		&a.ForTitle,
		&a.Status,
		&a.Comment,
	)
	if err != nil {
		return a, err
	}

	// convert date strings to time.Time
	a.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return a, err
	}
	a.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt)
	if err != nil {
		return a, err
	}
	// application date for some old records is "0000-00-00" which cannot be parsed properly and ends up as
	// "30-08-1754". However, don't want to return an error for these few records so will just ignore for now.
	a.Date, err = time.Parse("2006-01-02", applicationDate)
	if err != nil {
		// return a, err
	}

	return a, nil
}
