package application

var queries = map[string]string{
	"select-applications":             selectActiveApplications,
	"select-application-by-id":        selectApplicationByID,
	"select-applications-by-memberid": selectApplicationsByMemberID,
}

const selectApplications = `SELECT 
  ma.id                                   AS ID,
  ma.created_at                           AS CreatedAt,
  ma.updated_at                           AS UpdatedAt,
  ma.member_id                            AS MemberID,
  COALESCE(CONCAT(m.first_name, ' ', m.last_name), '') AS Member,
  IFNULL(ma.member_id_nominator, 0)       AS NominatorID,
  COALESCE(CONCAT(n.first_name, ' ', n.last_name), '') AS Nominator,
  IFNULL(ma.member_id_seconder, 0)        AS SeconderID,
  COALESCE(CONCAT(s.first_name, ' ', s.last_name), '') AS Seconder,
  ma.applied_on                           AS ApplicationDate,
  COALESCE(t.name, '')                    AS AppliedFor,
  ma.result                               AS Status,
  COALESCE(ma.comment,'')                 AS Comment
FROM
  ms_m_application ma
    LEFT JOIN
  member m ON ma.member_id = m.id
    LEFT JOIN
  member n ON ma.member_id_nominator = n.id
    LEFT JOIN
  member s ON ma.member_id_seconder = s.id
    LEFT JOIN
  ms_title t ON ma.ms_title_id = t.id 
WHERE 1 `

const selectActiveApplications = selectApplications + ` AND ma.active = 1 `

const selectApplicationByID = selectActiveApplications + ` AND ma.id = %v `

const selectApplicationsByMemberID = selectActiveApplications + ` AND ma.member_id = %v `
