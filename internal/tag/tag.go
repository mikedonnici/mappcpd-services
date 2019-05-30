// package tag provides access to tag records
package tag

import "github.com/cardiacsociety/web-services/internal/platform/datastore"

// Tag represents an area of professional interest
type Tag struct {
	ID          int    `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

// All returns all the active Tags
func All(ds datastore.Datastore) ([]Tag, error) {
	var xt []Tag
	q := Queries["select-tags"]
	rows, err := ds.MySQL.Session.Query(q)
	if err != nil {
		return xt, err
	}
	defer rows.Close()

	for rows.Next() {
		t := Tag{}
		err := rows.Scan(&t.ID, &t.Name, &t.Description)
		if err != nil {
			return xt, err
		}
		xt = append(xt, t)
	}
	return xt, nil
}
