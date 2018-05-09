package organisation

import (
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// Organisation defines a business, society or similar legal entity
type Organisation struct {
	OID    string         `json:"_id,omitempty" bson:"_id,omitempty"`
	ID     int            `json:"id" bson:"id"`
	Code   string         `json:"code,omitempty" bson:"code,omitempty"`
	Name   string         `json:"name" bson:"name"`
	Phone  string         `json:"phone,omitempty" bson:"phone,omitempty"`
	Fax    string         `json:"fax,omitempty" bson:"fax,omitempty"`
	Email  string         `json:"email,omitempty" bson:"email,omitempty"`
	URL    string         `json:"url,omitempty" bson:"url,omitempty"`
	Groups []Organisation `json:"groups,omitempty" bson:"groups,omitempty"`
}

// ByID fetches an organisation record
func ByID(ds datastore.Datastore, id int) (Organisation, error) {
	return orgByID(ds, id)
}

// All fetches all active Organisations
func All(ds datastore.Datastore) ([]Organisation, error) {
	return orgList(ds)
}

func orgByID(ds datastore.Datastore, id int) (Organisation, error) {

	var o Organisation

	query := `SELECT id, short_name, name, phone, fax, email, web 
				FROM organisation WHERE active = 1 AND id = ?`
	err := ds.MySQL.Session.QueryRow(query, id).Scan(&o.ID, &o.Code, &o.Name, &o.Phone, &o.Fax, &o.Email, &o.URL)
	if err != nil {
		return o, err
	}

	o.Groups, err = childOrgs(ds, o.ID)

	return o, err
}

func orgList(ds datastore.Datastore) ([]Organisation, error) {

	var xo []Organisation

	// Top-level organisations have parent_organisation_id = NULL
	query := `SELECT id, short_name, name, phone, fax, email, web 
				FROM organisation WHERE active = 1 AND parent_organisation_id IS NULL`

	rows, err := ds.MySQL.Session.Query(query)
	if err != nil {
		return xo, err
	}
	defer rows.Close()

	for rows.Next() {
		o := Organisation{}
		rows.Scan(&o.ID, &o.Code, &o.Name, &o.Phone, &o.Fax, &o.Email, &o.URL)

		var err error
		o.Groups, err = childOrgs(ds, o.ID)
		if err != nil {
			return xo, err
		}

		xo = append(xo, o)
	}

	return xo, nil
}

// childOrgs fetches active sub organisations (groups) for the organisation specified by id.
func childOrgs(ds datastore.Datastore, id int) ([]Organisation, error) {

	var xo []Organisation

	query := `SELECT id, short_name, name, phone, fax, email, web FROM 
				organisation WHERE active = 1 AND parent_organisation_id = ?`
	rows, err := ds.MySQL.Session.Query(query, id)
	if err != nil {
		return xo, err
	}
	defer rows.Close()

	for rows.Next() {
		o := Organisation{}
		rows.Scan(&o.ID, &o.Code, &o.Name, &o.Phone, &o.Fax, &o.Email, &o.URL)
		xo = append(xo, o)
	}

	return xo, nil
}
