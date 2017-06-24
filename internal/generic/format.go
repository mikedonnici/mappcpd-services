package generic

import (
	"database/sql"
	"time"
)

// Formats raw data into a slice of string maps
func toSlice(rows *sql.Rows) []map[string]string {

	// Records (slice of map) to return
	records := []map[string]string{}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		r := make(map[string]string)
		for i, v := range values {
			// Here we can check if the value is nil (NULL value)
			if v == nil {
				value = "NULL"
			} else {
				value = string(v)
			}

			// Create name / value pair in the record
			r[columns[i]] = value
		}

		// Add the record map (row) to our records map slice
		//fmt.Println(r)
		records = append(records, r)

	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return records
}

// MySQLDateToTime converts a MySQL date string (YYY-MM-DD) to a time.Time value
func MySQLDateToTime(s string) (time.Time, error) {

	return time.Parse("2006-01-02", s)
}
