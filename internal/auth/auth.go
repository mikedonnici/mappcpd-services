package auth

import (
	"errors"
	"fmt"

	"database/sql"

	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// AuthMember checks login & pass against db_
func AuthMember(u, p string) (int, string, error) {

	query := `SELECT id, concat(first_name, ' ', last_name) as name
		  FROM member WHERE primary_email = "%s" AND password = MD5("%s")`
	query = fmt.Sprintf(query, u, p)

	var id int
	var name string
	var errMsg error
	err := datastore.MySQL.Session.QueryRow(query).Scan(&id, &name)

	// zero rows is a failed login
	if err == sql.ErrNoRows {
		errMsg = errors.New("Login details incorrect")
		return id, name, errMsg
	}

	// DB or some other error
	if err != nil {
		return id, name, err
	}

	// Logged in!
	return id, name, nil
}

// AuthScope gets the authorizations scopes or 'roles' for a user by user (member) id.
func AuthScope(id int) ([]string, error) {

	// TODO - Actually look up scopes - how?
	ss := []string{"member"}

	return ss, nil
}
