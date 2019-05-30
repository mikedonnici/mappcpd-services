package qualification

import "github.com/cardiacsociety/web-services/internal/platform/datastore"

// Qualification is a formal qualification such as a degree, Masters, PHD etc
type Qualification struct {
	ID          int    `json:"id" bson:"id"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

// All returns all the active Qualifications
func All(ds datastore.Datastore) ([]Qualification, error) {
	var xq []Qualification
	q := Queries["select-qualifications"]
	rows, err := ds.MySQL.Session.Query(q)
	if err != nil {
		return xq, err
	}
	defer rows.Close()

	for rows.Next() {
		q := Qualification{}
		err := rows.Scan(&q.ID, &q.Code, &q.Name, &q.Description)
		if err != nil {
			return xq, err
		}
		xq = append(xq, q)
	}

	return xq, nil
}
