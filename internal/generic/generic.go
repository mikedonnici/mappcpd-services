package generic

import (
	"fmt"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// GetIDs returns a list of primary keys (id) from any table. Takes the table name and an sql
// clause.
func GetIDs(ds datastore.Datastore, table string, clause string) ([]int, error) {
	return GetIntCol(ds, table, "id", clause)
}

// GetIntCol will return integer values from a table, from the specified column.
// It is used for fetching ids or fk ids based on the sql clause.
func GetIntCol(ds datastore.Datastore, table, column, clause string, ) ([]int, error) {

	var ids []int

	sql := fmt.Sprintf("SELECT %s FROM %s %s", column, table, clause)
	rows, err := ds.MySQL.Session.Query(sql)
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
func GetRows(ds datastore.Datastore, sql string) ([]map[string]string, error) {

	rows, e := ds.MySQL.Session.Query(sql)
	if e != nil {
		return nil, e
	}

	r := toSlice(rows)

	return r, nil
}
