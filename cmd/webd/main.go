package main

import (
	"fmt"

	"github.com/34South/envr"

	r_ "github.com/mappcpd/web-services/cmd/webd/router"
	ds_ "github.com/mappcpd/web-services/internal/platform/datastore"
)

func main() {

	msg := fmt.Sprint("Initialising environment...")
	env := envr.New("myEnv", []string{
		"MAPPCPD_API_URL",
		"MAPPCPD_MYSQL_URL",
		"MAPPCPD_MYSQL_DESC",
		"MAPPCPD_MONGO_URL",
		"MAPPCPD_MONGO_DESC",
		"MAPPCPD_MONGO_DBNAME",
		"MAPPCPD_SHORT_LINK_URL",
		"MAPPCPD_SHORT_LINK_PREFIX",
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
	}).Auto()
	if env.Ready {
		msg += "ready!"
	}
	fmt.Println(msg)

	// Connect to the databases
	ds_.Connect()

	// Crank up the router
	r_.Start()
}
