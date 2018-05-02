package organisations

import (
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// Organisation defines a business, society or similar legal entity
type Organisation struct {
	OID  string `json:"_id" bson:"_id"`
	ID   int    `json:"id" bson:"id"`
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
	OrganisationContact
}

// OrganisationContact struct holds contact information for an Organisation
type OrganisationContact struct {
	Phone       string                   `json:"phone" bson:"phone"`
	Fax         string                   `json:"fax" bson:"fax"`
	Email       string                   `json:"email" bson:"email"`
	URL         string                   `json:"url" bson:"url"`
}

// OrganisationGroupPosition is a position that can be held within an Organisation Groups
// For example the 'President' (position) of the Board (group)
type OrganisationGroupPosition struct {
	ID   int    `json:"id" bson:"id"`
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}

// OrganisationByID fetches an organisation record from the default datastore
func OrganisationByID(id int) (Organisation, error) {
	return orgByID(id, datastore.MySQL)
}

// OrganisationByIDStore fetches an organisation record from the specified datastore - used for testing
func OrganisationByIDStore(id int, conn datastore.MySQLConnection) (Organisation, error) {
	return orgByID(id, conn)
}

// OrganisationList fetches a list of all the 'active' Organisations
func OrganisationsList() ([]Organisation, error) {
	return orgList(datastore.MySQL)
}

// OrganisationList fetches a list of all the 'active' Organisations from the specified datastore - used for testing
func OrganisationsListStore(conn datastore.MySQLConnection) ([]Organisation, error) {
	return orgList(conn)
}

// ChildOrganisations fetches organisations belonging to parentOrgID
func ChildOrganisations(parentOrgID int) ([]Organisation, error) {
	return childOrgs(parentOrgID, datastore.MySQL)
}

// ChildOrganisationsStore fetches organisations belonging to parentOrgID from the specified datastore - used for testing
func ChildOrganisationsStore(parentOrgID int, conn datastore.MySQLConnection) ([]Organisation, error) {
	return childOrgs(parentOrgID, conn)
}

func orgByID(id int, conn datastore.MySQLConnection) (Organisation, error) {
	var o Organisation
	query := `SELECT id, short_name, name FROM organisation WHERE active = 1 AND id = ?`
	err := conn.Session.QueryRow(query, id).Scan(&o.ID, &o.Code, &o.Name)
	return o, err
}

func orgList(conn datastore.MySQLConnection) ([]Organisation, error) {

	var xo []Organisation

	// Top-level organisations have parent_organisation_id = NULL
	query := `SELECT id, short_name, name FROM organisation WHERE active = 1 AND parent_organisation_id IS NULL`

	rows, err := conn.Session.Query(query)
	if err != nil {
		return xo, err
	}
	defer rows.Close()

	for rows.Next() {
		org := Organisation{}
		rows.Scan(&org.ID, &org.Code, &org.Name)
		xo = append(xo, org)
	}

	return xo, nil
}

// childOrgs fetches active sub organisations (groups) for the organisation specified by id.
func childOrgs(id int, conn datastore.MySQLConnection) ([]Organisation, error) {

	var xo []Organisation

	query := `SELECT id, short_name, name FROM organisation WHERE active = 1 AND parent_organisation_id = ?`
	rows, err := conn.Session.Query(query, id)
	if err != nil {
		return xo, err
	}
	defer rows.Close()

	for rows.Next() {
		g := Organisation{}
		rows.Scan(&g.ID, &g.Code, &g.Name)

		//err := g.SetCurrentGroupPositions()
		//if err != nil {
		//	return xo, err
		//}

		xo = append(xo, g)
	}

	return xo, nil
}