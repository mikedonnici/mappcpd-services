package datastore

import (
	"os"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLConnection struct {
	connectionString string
	Description      string
	Session          *sql.DB
}

// ConnectEnv establishes the Session using connection details stored in env vars
func (m *MySQLConnection) ConnectEnv() error {
	var err error
	m.connectionString = os.Getenv("MAPPCPD_MYSQL_URL")
	m.Description = os.Getenv("MAPPCPD_MYSQL_DESC")
	m.Session, err = sql.Open("mysql", m.connectionString)
	return err
}

// ConnectSource establishes the Session using the specified connection string - handy for testing.
func (m *MySQLConnection) ConnectSource(connectionString string) error {
	var err error
	m.connectionString = connectionString
	m.Description = "User specified"
	m.Session, err = sql.Open("mysql", m.connectionString)
	return err
}

// Close terminates the session - don't really need?
func (m *MySQLConnection) Close() {
	m.Session.Close()
}
