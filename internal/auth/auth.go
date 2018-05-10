package auth

import (
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
	err := ds.MySQL.Session.QueryRow(query).Scan(&id, &name)
	// Note: err == sql.ErrorNoRows for a failed login
	return id, name, err
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
	err := ds.MySQL.Session.QueryRow(query).Scan(&id, &name, &active, &locked)

	return id, name, err
}

// AdminAuthScope gets the authorizations scopes or 'roles' for an admin user by user id.
func AdminAuthScope(id int) ([]string, error) {

	// TODO - Actually look up scopes - how?
	ss := []string{"admin", "a", "b", "c"}

	return ss, nil
}
