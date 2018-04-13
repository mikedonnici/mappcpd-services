package models

import "github.com/mappcpd/web-services/internal/platform/datastore"

// Organisation defines a business, society or similar legal entity
type Organisation struct {
	OID         string `json:"_id" bson:"_id"`
	ID          int    `json:"id" bson:"id"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	OrganisationContact
}

// OrganisationContact struct holds contact information for an Organisation
type OrganisationContact struct {
	Phone       string                   `json:"phone" bson:"phone"`
	Fax         string                   `json:"fax" bson:"fax"`
	Email       string                   `json:"email" bson:"email"`
	URL         string                   `json:"url" bson:"url"`
	Locations   []OrganisationLocation   `json:"locations" bson:"locations"`
	Memberships []OrganisationMembership `json:"memberships" bson:"memberships"`
	Groups      []OrganisationGroup      `json:"groups" bson:"groups"`
}

// OrganisationLocation defines alternative contact data or locations
type OrganisationLocation struct {
	Preference  int    `json:"order" bson:"order"`
	Description string `json:"type" bson:"type"`
	Address     string `json:"address" bson:"address"`
	City        string `json:"city" bson:"city"`
	State       string `json:"state" bson:"state"`
	Postcode    string `json:"postcode" json:"postcode"`
	Country     string `json:"country" bson:"country"`
	Phone       string `json:"phone" bson:"phone"`
	Fax         string `json:"fax" bson:"fax"`
	Email       string `json:"email" bson:"email"`
	URL         string `json:"url" bson:"url"`
}

// OrganisationMembership types describe one or more membership categories available
// within an organisation. For Organisations that don't have members this is simply omitted
type OrganisationMembership struct {
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	TitlePNL    string `json:"titlePNL" bson:"titlePNL"`
	TitleShort  string `json:"titleShort" bson:"titleShort"`
	TitleFull   string `json:"titleFull" bson:"titleFull"`
}

// OrganisationGroup describes boards, councils, committees, working groups etc
// that fall under an Organisation
type OrganisationGroup struct {
	ID          int                         `json:"id" bson:"id"`
	Code        string                      `json:"code" bson:"code"`
	Name        string                      `json:"name" bson:"name"`
	Description string                      `json:"description" bson:"description"`
	Positions   []OrganisationGroupPosition `json:"positions" bson:"positions"`
}

// OrganisationGroupPosition is a position that can be held within an Organisation Groups
// For example the 'President' (position) of the Board (group)
type OrganisationGroupPosition struct {
	ID          int    `json:"id" bson:"id"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

// OrganisationByID fetches a single organisation record
func OrganisationByID(id int) (Organisation, error) {

	var o Organisation

	query := `SELECT id, short_name, name, 'no group description in model'
		  FROM organisation
		  WHERE 1
		  AND id = ?
		  AND active = 1`

	err := datastore.MySQL.Session.QueryRow(query, id).Scan(
		&o.ID,
		&o.Code,
		&o.Name,
		&o.Description,
	)
	if err != nil {
		return o, err
	}

	return o, nil
}

// OrganisationList fetches a list of all the 'active' Organisations
func OrganisationsList() ([]Organisation, error) {

	var xo []Organisation

	// Top-level organisations have parent_organisation_id = NULL
	query := `SELECT
		  id,
		  short_name,
		  name,
		  'no organisation description in model'
		  FROM organisation
		  WHERE 1
		  AND parent_organisation_id IS NULL
		  AND active = 1`

	rows, err := datastore.MySQL.Session.Query(query)
	if err != nil {
		return xo, err
	}
	defer rows.Close()

	for rows.Next() {
		org := Organisation{}
		rows.Scan(&org.ID, &org.Code, &org.Name, &org.Description)
		xo = append(xo, org)
	}

	return xo, nil
}

// OrganisationGroupsList fetches a list of all the 'active' Organisations. Currently this reads
// from the Organisation table which is a mash of standalone organisations as well as sub
// groups. Added parent_org_id column to filter top-level orgs from sub groups
// this function received the id of the top-level Organisation for which we wish to
// fetch groups.
func OrganisationGroupsList(id int) ([]OrganisationGroup, error) {

	var gs []OrganisationGroup

	query := `SELECT id, short_name, name, 'no group description in model'
		  FROM organisation
		  WHERE 1
		  AND parent_organisation_id = ?
		  AND active = 1`

	rows, err := datastore.MySQL.Session.Query(query, id)
	if err != nil {
		return gs, err
	}
	defer rows.Close()

	for rows.Next() {
		g := OrganisationGroup{}
		rows.Scan(&g.ID, &g.Code, &g.Name, &g.Description)

		err := g.SetCurrentGroupPositions()
		if err != nil {
			return gs, err
		}

		gs = append(gs, g)
	}

	return gs, nil
}

// SetCurrentGroupPositions sets the Positions field in an OrganisationGroup
func (og *OrganisationGroup) SetCurrentGroupPositions() error {

	var ps []OrganisationGroupPosition

	//ot.name as 'Group Type',
	// CONCAT(m.first_name, ' ', m.last_name) as 'Member',
	//mp.start_on as 'Start',
	//mp.end_on as 'End'

	query := `SELECT
		o.id as 'Group ID',
		o.short_name as 'Group Code',
		o.name as 'Group Name',
		p.name as 'Position'
		FROM
		mp_m_position mp
		LEFT JOIN
		mp_position p ON mp.mp_position_id = p.id
		LEFT JOIN
		organisation o ON mp.organisation_id = o.id
		LEFT JOIN
		organisation_type ot on o.organisation_type_id = ot.id
		LEFT JOIN
		member m ON mp.member_id = m.id
		WHERE 1
		AND p.name != 'First Council Affiliation'
		AND p.name != 'Second Council Affiliation'
		AND p.name != 'Third Council Affiliation'
		AND ((mp.end_on IS NULL) OR (mp.end_on = '0000-00-00') OR (mp.end_on > NOW()))
		AND o.id = ?;`

	rows, err := datastore.MySQL.Session.Query(query, og.ID)
	if err != nil {
		return err
	}

	for rows.Next() {
		p := OrganisationGroupPosition{}
		rows.Scan(&p.ID, &p.Code, &p.Name, &p.Description)
		ps = append(ps, p)
	}

	return nil
}
