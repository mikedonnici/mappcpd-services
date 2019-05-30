package organisation

import (
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
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

// All fetches all active Organisations
func All(ds datastore.Datastore) ([]Organisation, error) {
	var xo []Organisation
	rows, err := ds.MySQL.Session.Query(querySelectParentOrganisations)
	if err != nil {
		return xo, err
	}
	defer rows.Close()
	for rows.Next() {
		o := Organisation{}
		rows.Scan(&o.ID, &o.Code, &o.Name, &o.Phone, &o.Fax, &o.Email, &o.URL)
		var err error
		o.Groups, err = ByParentID(ds, o.ID)
		if err != nil {
			return xo, err
		}
		xo = append(xo, o)
	}
	return xo, nil
}

// ByID fetches an organisation record
func ByID(ds datastore.Datastore, id int) (Organisation, error) {
	var o Organisation
	err := ds.MySQL.Session.QueryRow(querySelectOrganisationByID, id).Scan(&o.ID, &o.Code, &o.Name, &o.Phone, &o.Fax, &o.Email, &o.URL)
	if err != nil {
		return o, err
	}
	o.Groups, err = ByParentID(ds, o.ID)

	return o, err
}

// ByParentID fetches active sub organisations (groups) for the organisation specified by id.
func ByParentID(ds datastore.Datastore, id int) ([]Organisation, error) {
	var xo []Organisation
	rows, err := ds.MySQL.Session.Query(querySelectOrganisationByParentID, id)
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

// ByTypeID fetchjes organisations by type id
func ByTypeID(ds datastore.Datastore, organisationTypeID int) ([]Organisation, error) {
	var xo []Organisation
	rows, err := ds.MySQL.Session.Query(querySelectOrganisationByTypeID, organisationTypeID)
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
