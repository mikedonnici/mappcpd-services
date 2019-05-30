package speciality

var Queries = map[string]string{
	"select-specialities": selectSpecialities,
}

const selectSpecialities = `
SELECT 
	id,
    mp.name,
    COALESCE(mp.description, '')
FROM
    mp_speciality mp
WHERE
	active = 1`
