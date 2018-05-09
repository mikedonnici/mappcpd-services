package auth

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// AuthMember checks login & pass against db. Check for md5() or encrypted string.
// Latter is a workaround to allow the old member app to get a token for file uploads.
func AuthMember(ds datastore.Datastore, u, p string) (int, string, error) {

	query := `SELECT id, concat(first_name, ' ', last_name) as name
		  FROM member WHERE primary_email = "%s" AND (password = MD5("%s") OR password = "%s")`
	query = fmt.Sprintf(query, u, p, p)

	var id int
	var name string
	var errMsg error
	err := ds.MySQL.Session.QueryRow(query).Scan(&id, &name)

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

// AdminAuth authenticates an admin user against the db. It received username and password
// strings and returns the id and name of the authenticated admin
func AdminAuth(ds datastore.Datastore, u, p string) (int, string, error) {

	query := `SELECT id, name, active, locked FROM ad_user WHERE
	          username = "%s" AND (password = MD5("%s") OR password = "%s")`
	query = fmt.Sprintf(query, u, p, p)

	var id int
	var name string
	var active int
	var locked int
	var errMsg error
	err := ds.MySQL.Session.QueryRow(query).Scan(&id, &name, &active, &locked)

	// zero rows is a failed login
	if err == sql.ErrNoRows {
		errMsg = errors.New("Login details incorrect")
		return id, name, errMsg
	}

	// Account is locked
	if locked == 1 {
		errMsg = errors.New("Admin account locked - contact systems administrator")
		return id, name, errMsg
	}

	// Account inactive
	if active == 0 {
		errMsg = errors.New("Admin account inactive - contact systems administrator")
		return id, name, errMsg
	}

	// DB or some other error
	if err != nil {
		return id, name, err
	}

	// Authenticated
	return id, name, nil
}

// AdminAuthScope gets the authorizations scopes or 'roles' for an admin user by user id.
func AdminAuthScope(id int) ([]string, error) {

	// TODO - Actually look up scopes - how?
	ss := []string{"admin", "a", "b", "c"}

	return ss, nil
}
