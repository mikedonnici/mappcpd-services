// package speciality provides access to speciality records
package speciality

import "github.com/cardiacsociety/web-services/internal/platform/datastore"

// Speciality represents an area of professional interest
type Speciality struct {
	ID          int    `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

// All returns all the active Specialities
func All(ds datastore.Datastore) ([]Speciality, error) {
	var xs []Speciality
	q := Queries["select-specialities"]
	rows, err := ds.MySQL.Session.Query(q)
	if err != nil {
		return xs, err
	}
	defer rows.Close()

	for rows.Next() {
		s := Speciality{}
		err := rows.Scan(&s.ID, &s.Name, &s.Description)
		if err != nil {
			return xs, err
		}
		xs = append(xs, s)
	}

	return xs, nil
}
