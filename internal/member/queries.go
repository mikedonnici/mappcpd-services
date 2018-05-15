package member

var Queries = map[string]string{
	"select-member": selectMember,
	"select-membership-title": selectMembershipTitle,
	"select-membership-title-history": selectMembershipTitleHistory,
	"select-member-qualifications": selectMemberQualifications,
	"select-member-positions": selectMemberPositions,
	"select-member-specialities": selectMemberSpecialities,
}

const selectMember = `SELECT 
	active,
    created_at as CreatedAt,
    updated_at as UpdatedAt,
    COALESCE(first_name, '') as FirstName,  
    COALESCE(middle_names, '') as MiddleNames,
    COALESCE(last_name, '') as LastName,
    CONCAT(COALESCE(suffix, ''), ' ', COALESCE(qualifications_other, '')) as PostNom,
    COALESCE(gender, '') as Gender,
    COALESCE(date_of_birth, '') as DOB,
    COALESCE(primary_email, '') as Email,
    COALESCE(secondary_email, '') as Email2,
    COALESCE(mobile_phone, '') as Mobile,
    consent_directory as ConsentDirectory,
    consent_contact as ConsentContact
FROM
    member
WHERE
    id = ?`

const selectMembershipTitle = `SELECT 
    COALESCE(mt.name, '')
FROM
    ms_title mt
        INNER JOIN
    ms_m_title mmt ON mt.id = mmt.ms_title_id
WHERE
	current = 1 AND mmt.member_id = ? 
ORDER BY mmt.id DESC
LIMIT 1`

const selectMembershipTitleHistory = `SELECT 
    COALESCE(mmt.granted_on, ''),
    'no code',
    COALESCE(mt.name, ''),
    COALESCE(mt.description, ''),
    COALESCE(mmt.comment, '')
FROM
    ms_title mt
        INNER JOIN
    ms_m_title mmt ON mt.id = mmt.ms_title_id
WHERE
    mmt.member_id = ?
ORDER BY mmt.id DESC`

const selectMemberQualifications = `SELECT 
    COALESCE(mq.short_name, ''),
    COALESCE(mq.name, ''),
    COALESCE(mq.description, ''),
    COALESCE(mmq.year, '')
FROM
    mp_m_qualification mmq
        LEFT JOIN
    mp_qualification mq ON mmq.mp_qualification_id = mq.id
WHERE
    mmq.member_id = ?
ORDER BY year DESC`

const selectMemberPositions = `SELECT 
    COALESCE(organisation.short_name, ''),
    COALESCE(organisation.name, ''),
    COALESCE(mp.short_name, ''),
    COALESCE(mp.name, ''),
    COALESCE(mp.description, ''),
    COALESCE(mmp.start_on, ''),
    COALESCE(mmp.end_on, '')
FROM
    mp_m_position mmp
        LEFT JOIN
    mp_position mp ON mmp.mp_position_id = mp.id
        LEFT JOIN
    organisation ON mmp.organisation_id = organisation.id
WHERE
    mmp.member_id = ?`

const selectMemberSpecialities = `SELECT 
    COALESCE(s.name, ''),
    COALESCE(s.description, ''),
    COALESCE(ms.start_on, '')
FROM
    mp_m_speciality ms
        LEFT JOIN
    mp_speciality s ON ms.mp_speciality_id = s.id
WHERE
    ms.member_id = ?`