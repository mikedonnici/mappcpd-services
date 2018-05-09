package datastore

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type MySQLConnection struct {
	DSN     string // Data Desc Name - connection string
	Desc    string
	Session *sql.DB
}

// NewMySQLConnection returns a pointer to an initialised MySQLConnection
//func NewMySQLConnection(dsn, desc string) *MySQLConnection {
//	return &MySQLConnection{
//		DSN:         dsn,
//		Desc: desc,
//	}
//}

// ConnectSource establishes the Session using the specified connection string - handy for testing.
func (m *MySQLConnection) Connect() error {
	err := m.checkFields()
	if err != nil {
		return err
	}
	m.Session, err = sql.Open("mysql", m.DSN)
	return err
}

// Close terminates the Session - don't really need?
func (m *MySQLConnection) Close() {
	m.Session.Close()
}

func (m *MySQLConnection) checkFields() error {
	if m.DSN == "" {
		return errors.New("MySQLConnection.DSN (data source name / connection string) is not set")
	}
	if m.Desc == "" {
		return errors.New("MySQLConnection.Desc is not set")
	}
	return nil
}
