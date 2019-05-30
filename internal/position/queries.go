package position

var queries = map[string]string{
	"select-positions":      selectActivePositions,
	"select-position-by-id": selectPositionByID,
}

const selectPositions = `
SELECT 
    mp.id AS MemberPositionID,
    mp.created_at AS CreatedAt,
    mp.updated_at AS UpdatedAt,
    mp.member_id  AS MemberID,
    COALESCE(CONCAT(m.first_name, ' ', m.last_name), '') AS Member,
    COALESCE(m.primary_email, '') AS Email,
    mp.mp_position_id         AS PositionID,
    p.name                    AS PositionName,
    mp.organisation_id        AS OrganisationID,
    COALESCE(o.name, '')      AS OrganisationName,
    COALESCE(mp.start_on, '') AS StartDate,
    COALESCE(mp.end_on, '')   AS EndDate,
    COALESCE(mp.comment, '')  AS Comment
FROM
    mp_m_position mp
        LEFT JOIN
    member m ON mp.member_id = m.id
        LEFT JOIN
    mp_position p ON mp.mp_position_id = p.id
		LEFT JOIN 
	organisation o ON mp.organisation_id = o.id
WHERE 1 `

const selectActivePositions = selectPositions + ` AND mp.active = 1 `

const selectPositionByID = selectActivePositions + ` AND mp.id = %v `
