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

// ByID fetches an organisation record from the default datastore
func ByID(id int) (Organisation, error) {
	return orgByID(id, datastore.MySQL)
}

// ByIDStore fetches an organisation record from the specified datastore - used for testing
func ByIDStore(id int, conn datastore.MySQLConnection) (Organisation, error) {
	return orgByID(id, conn)
}

// All fetches all active Organisations
func All() ([]Organisation, error) {
	return orgList(datastore.MySQL)
}

// AllStore fetches all active Organisations from the specified datastore - used for testing
func AllStore(conn datastore.MySQLConnection) ([]Organisation, error) {
	return orgList(conn)
}

func orgByID(id int, conn datastore.MySQLConnection) (Organisation, error) {

	var o Organisation

	query := `SELECT id, short_name, name, phone, fax, email, web 
				FROM organisation WHERE active = 1 AND id = ?`
	err := conn.Session.QueryRow(query, id).Scan(&o.ID, &o.Code, &o.Name, &o.Phone, &o.Fax, &o.Email, &o.URL)
	if err != nil {
		return o, err
	}

	o.Groups, err = childOrgs(o.ID, conn)

	return o, err
}

func orgList(conn datastore.MySQLConnection) ([]Organisation, error) {

	var xo []Organisation

	// Top-level organisations have parent_organisation_id = NULL
	query := `SELECT id, short_name, name, phone, fax, email, web 
				FROM organisation WHERE active = 1 AND parent_organisation_id IS NULL`

	rows, err := conn.Session.Query(query)
	if err != nil {
		return xo, err
	}
	defer rows.Close()

	for rows.Next() {
		o := Organisation{}
		rows.Scan(&o.ID, &o.Code, &o.Name, &o.Phone, &o.Fax, &o.Email, &o.URL)

		var err error
		o.Groups, err = childOrgs(o.ID, conn)
		if err != nil {
			return xo, err
		}

		xo = append(xo, o)
	}

	return xo, nil
}

// childOrgs fetches active sub organisations (groups) for the organisation specified by id.
func childOrgs(id int, conn datastore.MySQLConnection) ([]Organisation, error) {

	var xo []Organisation

	query := `SELECT id, short_name, name, phone, fax, email, web FROM 
				organisation WHERE active = 1 AND parent_organisation_id = ?`
	rows, err := conn.Session.Query(query, id)
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
