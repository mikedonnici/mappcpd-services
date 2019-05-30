package tag

var Queries = map[string]string{
	"select-tags": selectTags,
}

const selectTags = `
SELECT 
	id,
    name,
    COALESCE(description, '')
FROM
    mp_tag
WHERE
	active = 1`
