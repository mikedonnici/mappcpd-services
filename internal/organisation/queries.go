package organisation

const querySelectActiveOrganisation = `SELECT 
	id, 
	short_name, 
	name, phone, 
	fax, 
	email, 
	web 
FROM 
	organisation 
WHERE 
	active = 1`

const querySelectParentOrganisations = querySelectActiveOrganisation + ` AND parent_organisation_id IS NULL ORDER BY name`

const querySelectOrganisationByID = querySelectActiveOrganisation + ` AND id = ? ORDER BY name`

const querySelectOrganisationByParentID = querySelectActiveOrganisation + ` AND parent_organisation_id = ? ORDER BY name`

const querySelectOrganisationByTypeID = querySelectActiveOrganisation + ` AND organisation_type_id = ? ORDER BY name`
