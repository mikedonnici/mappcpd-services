package member

var queries = map[string]string{
	"insert-member-row":                      insertMemberRow,
	"insert-member-qualification-row":        insertMemberQualificationRow,
	"insert-member-position-row":             insertMemberPositionRow,
	"insert-member-speciality-row":           insertMemberSpecialityRow,
	"insert-member-accreditation-row":        insertMemberAccreditationRow,
	"insert-member-tag-row":                  insertMemberTagRow,
	"insert-member-application-row":          insertMemberApplicationRow,
	"insert-member-contact-row":              insertMemberContactRow,
	"insert-member-status-row":               insertMemberStatusRow,
	"select-member":                          selectMember,
	"select-member-honorific":                selectMemberHonorific,
	"select-member-country":                  selectMemberCountry,
	"select-member-contact-locations":        selectMemberContactLocations,
	"select-membership-title":                selectMembershipTitle,
	"select-membership-title-history":        selectMembershipTitleHistory,
	"select-membership-status":               selectMembershipStatus,
	"select-membership-status-history":       selectMembershipStatusHistory,
	"select-member-qualifications":           selectMemberQualifications,
	"select-member-accreditations":           selectMemberAccreditations,
	"select-member-positions":                selectMemberPositions,
	"select-member-specialities":             selectMemberSpecialities,
	"select-member-tags":                     selectMemberTags,
	"update-member-current-status":           updateMemberCurrentStatus,
	"update-member-deactivate-subscriptions": updateMemberDeactivateSubscriptions,
}

const insertMemberRow = `
INSERT INTO member (
    acl_member_role_id, 
    a_name_prefix_id, 
    country_id, 
    consent_directory, 
    consent_contact, 
    created_at, 
    updated_at, 
    date_of_birth, 
    gender, 
    first_name, 
    middle_names, 
    last_name, 
    suffix, 
    mobile_phone, 
    primary_email,
    password
) VALUES (
    ?, ?, ?, ?, ?, 
    NOW(), NOW(), 
    ?, ?, ?, ?, ?, ?, ?, ?, MD5("a-not-so-random-string")
)`

const insertMemberQualificationRow = `
INSERT INTO mp_m_qualification (
    member_id, 
    mp_qualification_id, 
    organisation_id, 
    created_at, 
    updated_at, 
    year, 
    qualification_suffix,
    comment
) VALUES (?, ?, ?, NOW(), NOW(), ?, ?, ?)
`

const insertMemberPositionRow = `
INSERT INTO mp_m_position (
    member_id, 
    mp_position_id, 
    organisation_id, 
    created_at, 
    updated_at
) VALUES (?, ?, ?, NOW(), NOW())
`

const insertMemberSpecialityRow = `
INSERT INTO mp_m_speciality (
    member_id, 
    mp_speciality_id, 
    created_at, 
    updated_at, 
    preference,
    comment
) VALUES (?, ?, NOW(), NOW(), ?, ?)
`

const insertMemberAccreditationRow = `
INSERT INTO mp_m_accreditation (
    member_id, 
    mp_accreditation_id, 
    created_at, 
    updated_at, 
    start_on,
    end_on,
    comment
) VALUES (?, ?, NOW(), NOW(), ?, ?, ?)
`

const insertMemberTagRow = `
INSERT INTO mp_m_tag (
    member_id, 
    mp_tag_id, 
    created_at, 
    updated_at
) VALUES (?, ?, NOW(), NOW())
`

const insertMemberApplicationRow = `
INSERT INTO ms_m_application(
  member_id, 
  member_id_nominator, 
  member_id_seconder, 
  ms_title_id, 
  updated_at, 
  applied_on, 
  comment) 
VALUES (?, ?, ?, ?, NOW(), NOW(), ?)`

const insertMemberContactRow = `
INSERT INTO mp_m_contact(
  member_id, 
  mp_contact_type_id, 
  country_id,
  updated_at,
  phone,
  fax,
  email,
  web,
  address1,
  address2,
  address3,
  locality,
  state,
  postcode 
  ) 
VALUES (?, ?, ?, NOW(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

const insertMemberStatusRow = `
INSERT INTO ms_m_status(
    member_id, 
    ms_status_id, 
    current, 
    updated_at,
    comment
) 
VALUES (?, ?, ?, NOW(), ?)`

const selectMember = `SELECT 
	active,
    created_at as CreatedAt,
    updated_at as UpdatedAt,
    COALESCE(first_name, '') as FirstName,  
    COALESCE(middle_names, '') as MiddleNames,
    COALESCE(last_name, '') as LastName,
    COALESCE(suffix, '') as PostNom,
    COALESCE(qualifications_other, '') as QualificationsOther, 
    COALESCE(gender, '') as Gender,
    COALESCE(date_of_birth, '') as DateOfBirth,
    COALESCE(date_of_entry, '') as DateOfEntry,
    COALESCE(primary_email, '') as Email,
    COALESCE(secondary_email, '') as Email2,
    COALESCE(mobile_phone, '') as Mobile,
    COALESCE(journal_number, '') as JournalNumber,
    COALESCE(bpay_number, '') as BpayNumber,
    consent_directory as ConsentDirectory,
    consent_contact as ConsentContact
FROM
    member
WHERE
    id = ?`

const selectMemberHonorific = `SELECT
	COALESCE(a.name, '') FROM a_name_prefix a
	RIGHT JOIN member m ON m.a_name_prefix_id = a.id
    WHERE m.id = ?`

const selectMemberCountry = `SELECT
	COALESCE(c.name, '') FROM country c
	RIGHT JOIN member m ON m.country_id = c.id
	WHERE m.id = ?`

const selectMemberContactLocations = `SELECT 
    COALESCE(mpct.name, ''),
    CONCAT(COALESCE(mpmc.address1, ''), '\n', COALESCE(mpmc.address2, ''), '\n', COALESCE(mpmc.address3, '')),
    COALESCE(mpmc.locality, ''),
    COALESCE(mpmc.state, ''),
    COALESCE(mpmc.postcode, ''),
    COALESCE(country.name, ''),
    COALESCE(mpmc.phone, ''),
    COALESCE(mpmc.fax, ''),
    COALESCE(mpmc.email, ''),
    COALESCE(mpmc.web, ''),
    COALESCE(mpct.order, '')
FROM
    mp_m_contact mpmc
        LEFT JOIN
    mp_contact_type mpct ON mpmc.mp_contact_type_id = mpct.id
        LEFT JOIN
    country ON mpmc.country_id = country.id
WHERE
    mpmc.member_id = ?
GROUP BY mpmc.id
ORDER BY mpct.order ASC`

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

const selectMembershipStatus = `SELECT 
    COALESCE(ms.name, '')
FROM
    ms_status ms
        INNER JOIN
    ms_m_status mms ON ms.id = mms.ms_status_id
WHERE
	current = 1 AND mms.member_id = ? 
ORDER BY mms.id DESC
LIMIT 1`

const selectMembershipStatusHistory = `SELECT
	mms.created_at as Date,
    COALESCE(ms.name, ''),
    COALESCE(ms.description, ''),
    COALESCE(mms.comment, '')
FROM
    ms_status ms
        INNER JOIN
    ms_m_status mms ON ms.id = mms.ms_status_id
WHERE
    mms.member_id = ?
ORDER BY mms.id DESC`

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

const selectMemberAccreditations = `SELECT 
    COALESCE(ma.short_name, ''),
    COALESCE(ma.name, ''),
    COALESCE(ma.description, ''),
    COALESCE(mma.start_on, ''),
    COALESCE(mma.end_on, '')
FROM
    mp_m_accreditation mma
        LEFT JOIN
    mp_accreditation ma ON mma.mp_accreditation_id = ma.id
WHERE
    mma.member_id = ?
ORDER BY mma.start_on DESC`

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
    ms.member_id = ?
ORDER BY ms.preference ASC`

const selectMemberTags = `SELECT 
    COALESCE(t.name, '') as Tag
FROM
    mp_m_tag mt
        LEFT JOIN
    mp_tag t ON mt.mp_tag_id = t.id
WHERE
    mt.member_id = ?`

// ensure only ONE member status record is current
const updateMemberCurrentStatus = `
UPDATE ms_m_status SET current = 0 
WHERE id != ? AND member_id = ?`

// de-activate all subscriptions for a member
const updateMemberDeactivateSubscriptions = `UPDATE fn_m_subscription SET active = 0 WHERE member_id = ?`
