package models

import "github.com/mappcpd/api/db"

func ReportModulesByDate() (map[string]int, error) {

	r := make(map[string]int)

	sql := `SELECT DATE_FORMAT(created_at, '%Y-%m') as 'Date',
  		Count(*) As 'Modules'
		FROM ol_m_module
		GROUP BY Year(created_at), Month(created_at);`
	rows, err := db.MySQL.Session.Query(sql)
	defer rows.Close()

	if err != nil {
		return r, err
	}

	for rows.Next() {
		var d string
		var c int
		rows.Scan(&d, &c)
		r[d] = c
	}

	return r, nil
}

// ReportPointsByRecordDate groups cpd activity (points) by date the record was created.
func ReportPointsByRecordDate() (map[string]float32, error) {

	r := make(map[string]float32)

	sql := `SELECT DATE_FORMAT(created_at, '%Y-%m') as 'Date',
  		SUM(quantity * points_per_unit) AS 'Points'
		FROM ce_m_activity
		GROUP BY Year(created_at), Month(created_at)
		ORDER BY Year(created_at), Month(created_at);`
	rows, err := db.MySQL.Session.Query(sql)
	defer rows.Close()

	if err != nil {
		return r, err
	}

	for rows.Next() {
		var d string
		var p float32
		rows.Scan(&d, &p)
		r[d] = p
	}

	return r, nil
}

// ReportPointsByActivityDate groups cpd activity (points) by the date the activity occurred
func ReportPointsByActivityDate() (map[string]float32, error) {

	r := make(map[string]float32)

	sql := `SELECT DATE_FORMAT(activity_on, '%Y-%m') as 'Date',
  		SUM(quantity * points_per_unit) AS 'Points'
		FROM ce_m_activity
		GROUP BY Year(activity_on), Month(created_at)
		ORDER BY Year(activity_on), Month(created_at);`
	rows, err := db.MySQL.Session.Query(sql)
	defer rows.Close()

	if err != nil {
		return r, err
	}

	for rows.Next() {
		var d string
		var p float32
		rows.Scan(&d, &p)
		r[d] = p
	}

	return r, nil
}
