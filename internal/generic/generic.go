package internal

import (
	"fmt"

	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// GetIDs returns a list of ids from any table. Takes the table name (t) and a filter (f)
// Note the table MUST have a field called `id`
func GetIDs(t string, f string) ([]int, error) {

	var ids []int

	sql := fmt.Sprintf("SELECT id FROM %s %s", t, f)
	rows, err := datastore.MySQL.Session.Query(sql)
	if err != nil {
		return ids, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids, nil
}

// GetRows runs any query and returns a map slice where each slice is a row
func GetRows(sql string) ([]map[string]string, error) {

	rows, e := datastore.MySQL.Session.Query(sql)
	if e != nil {
		return nil, e
	}

	r := toSlice(rows)

	return r, nil
}

// Execute
func Execute(sql string) error {
	_, err := datastore.MySQL.Session.Exec(sql)
	return err
}
