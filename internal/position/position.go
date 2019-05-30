// Package position provides a way to query member positions within organisational entities such as
// councils, commitees and so on.
package position

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// Position represents a member's position or association with an organisation or similar entity.
// MemberPositionID is the primary id as this is data that originates from a junction table - ie
// member positions. ID and Name identify the type of position such as chair, member and so on.
type Position struct {
	MemberPositionID int       `json:"memberPositionId" bson:"memberPositionId"`
	CreatedAt        time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt"`
	MemberID         int       `json:"memberId" bson:"memberId"`
	Member           string    `json:"member" bson:"member"`
	Email            string    `json:"email" bson:"email"`
	ID               int       `json:"id" bson:"id"`
	Name             string    `json:"name" bson:"name"`
	OrganisationID   int       `json:"organisationId" bson:"organisationId"`
	OrganisationName string    `json:"organisationName" bson:"organisationName"`
	StartDate        time.Time `json:"startDate" bson:"startDate"`
	EndDate          time.Time `json:"endDate" bson:"endDate"`
	Comment          string    `json:"comment" bson:"comment"`
}

// ByID fetches a Position by member-position ID
func ByID(ds datastore.Datastore, memberPositionID int) (Position, error) {
	var p Position
	q := fmt.Sprintf(queries["select-position-by-id"], memberPositionID)
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

// ByIDs returns multiple Position values identified by memberPositionIDs
func ByIDs(ds datastore.Datastore, memberPositionIDs []int) ([]Position, error) {
	idList := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(memberPositionIDs)), ","), "[]")
	q := queries["select-positions"] + fmt.Sprintf(" AND mp.id IN (%s)", idList)
	return execute(ds, q)
}

func execute(ds datastore.Datastore, query string) ([]Position, error) {

	var xp []Position

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

		xp = append(xp, p)
	}

	return xp, nil
}

// scanRow scans the current row via *sql.Rows
func scanRow(row *sql.Rows) (Position, error) {

	var p Position

	// values that will need some manipulation
	var createdAt, updatedAt string // data dates
	var startDate, endDate string   // position dates

	err := row.Scan(
		&p.MemberPositionID,
		&createdAt,
		&updatedAt,
		&p.MemberID,
		&p.Member,
		&p.Email,
		&p.ID,
		&p.Name,
		&p.OrganisationID,
		&p.OrganisationName,
		&startDate,
		&endDate,
		&p.Comment,
	)
	if err != nil {
		return p, err
	}

	p.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return p, err
	}
	p.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt)
	if err != nil {
		return p, err
	}
	p.StartDate, _ = time.Parse("2006-01-02", startDate) // ignore bung dates for this field
	p.EndDate, _ = time.Parse("2006-01-02", endDate)     // ignore bung dates for this field

	return p, nil
}
