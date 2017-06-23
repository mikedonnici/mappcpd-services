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
		"MYSQL_URL",
		"MYSQL_SRC",
		"MONGO_URL",
		"MONGO_DB",
		"MONGO_SRC",
		"BASE_URL",
		"SHORT_LINK_BASE_URL",
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
