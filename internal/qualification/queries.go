package qualification

var Queries = map[string]string{
	"select-qualifications": selectQualifications,
}

const selectQualifications = `SELECT 
	id,
    COALESCE(mq.short_name, ''),
    COALESCE(mq.name, ''),
    COALESCE(mq.description, '')
FROM
    mp_qualification mq
WHERE
	active = 1`
