package datastore

import (
	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/34South/envr"
)

func init() {
	envr.New("datastoreEnv", []string{
		"MAPPCPD_MYSQL_URL",
		"MAPPCPD_MYSQL_DESC",
	}).Auto()
}

type MySQLConnection struct {
	url     string
	Source  string
	Session *sql.DB
}

// Connects to MySQL server, returns an error if it fails
func (m *MySQLConnection) Connect() error {

	// Set properties
	m.url = os.Getenv("MAPPCPD_MYSQL_URL")
	m.Source = os.Getenv("MAPPCPD_MYSQL_DESC")

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
