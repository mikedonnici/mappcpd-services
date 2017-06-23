package datastore

import (
	"os"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLConnection struct {
	url     string
	Source  string
	Session *sql.DB
}

// Connects to MySQL server, returns an error if it fails
func (m *MySQLConnection) Connect() error {

	// Set properties
	m.url = os.Getenv("MYSQL_URL")
	m.Source = os.Getenv("MYSQL_SRC")

	// Establish session
	var err error
	m.Session, err = sql.Open("mysql", m.url)
	if err != nil {
		return err
	}

	return nil
}

// Close terminates the session - don't really need?
func (m *MySQLConnection) Close() {
	m.Session.Close()
}
