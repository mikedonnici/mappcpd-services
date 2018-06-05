// Package events provides access to Events data
package events

import (
	"database/sql"
	"log"
	"math"
	"time"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/pkg/errors"
)

// Event is a conference, workshop or some other calendar event that is relevant to CPD activity
type Event struct {
	ID          int    `json:"id" bson:"id"`
	DateCreated string `json:"dateCreated" bson:"dateCreated"`
	DateUpdated string `json:"dateUpdated" bson:"dateUpdated"`
	DateStart   string `json:"dateStart" bson:"dateStart"`
	DateEnd     string `json:"dateEnd" bson:"dateEnd"`
	Location    string `json:"location" bson:"location"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	URL         string `json:"url" bson:"url"`
}

// ByID fetches a single Event by ID
func ByID(ds datastore.Datastore, id int) (Event, error) {

	// Create Note value
	e := Event{ID: id}

	// Coalesce any NULL-able fields
	q := `SELECT created_at, updated_at,
		  COALESCE(start_on, ''),
		  COALESCE(end_on, ''),
		  COALESCE(location, ''),
		  COALESCE(name, ''),
		  COALESCE(description, ''),
		  COALESCE(information_url, '')
          FROM ce_event WHERE id = ?
          ORDER BY start_on DESC`

	err := ds.MySQL.Session.QueryRow(q, id).Scan(
		&e.DateCreated,
		&e.DateUpdated,
		&e.DateStart,
		&e.DateEnd,
		&e.Location,
		&e.Name,
		&e.Description,
		&e.URL,
	)

	return e, err
}

// ByDateRange returns Events that have a start date within the specified date range, including the start and end dates
func ByDateRange(ds datastore.Datastore, start, end time.Time) ([]Event, error) {

	var xe []Event

	// MySQL DATE format
	sd := start.Format("2006-01-02")
	ed := end.Format("2006-01-02")

	q := `SELECT id, created_at, updated_at,
		  COALESCE(start_on, ''),
		  COALESCE(end_on, ''),
		  COALESCE(location, ''),
		  COALESCE(name, ''),
		  COALESCE(description, ''),
		  COALESCE(information_url, '')
		  FROM ce_event WHERE
		  start_on >= ? AND end_on <= ?
		  ORDER BY start_on DESC`

	rows, err := ds.MySQL.Session.Query(q, sd, ed)
	switch {
	case err == sql.ErrNoRows:
		return xe, nil
	case err != nil:
		msg := "ByDateRange() sql error"
		return xe, errors.Wrap(err, msg)
	}
	defer rows.Close()

	for rows.Next() {

		e := Event{}

		err := rows.Scan(
			&e.ID,
			&e.DateCreated,
			&e.DateUpdated,
			&e.DateStart,
			&e.DateEnd,
			&e.Location,
			&e.Name,
			&e.Description,
			&e.URL,
		)
		if err != nil {
			msg := "ByDateRange() failed to scan row"
			log.Println(msg, err)
			return xe, errors.Wrap(err, msg)
		}

		xe = append(xe, e)
	}

	return xe, nil
}

// DaysRange returns events with start dates falling within daysBack to daysForward.
func DaysRange(ds datastore.Datastore, daysBack, daysForward int) ([]Event, error) {

	// ensure daysBack is negative, and daysForward is positive
	daysBack = -int(math.Abs(float64(daysBack)))
	daysForward = int(math.Abs(float64(daysForward)))

	from := time.Now().AddDate(0, 0, daysBack)
	to := time.Now().AddDate(0, 0, daysForward)

	return ByDateRange(ds, from, to)
}

// Past is a convenience function that fetches events with a start date that falls between today and n days ago.
// If n < 0 it will return all past events.
func Past(ds datastore.Datastore, days int) ([]Event, error) {

	if days < 0 {
		days = 35000 // 100 years should be enough!
	}

	return DaysRange(ds, days, 0)
}

// Future is a convenience function that fetches events with a start date that falls between today and n days forward.
// If n < 0 it will return all future events.
func Future(ds datastore.Datastore, days int) ([]Event, error) {

	if days < 0 {
		days = 35000 // 100 years should be enough!
	}

	return DaysRange(ds, 0, days)
}
