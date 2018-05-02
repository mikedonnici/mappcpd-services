package rest

import (
	"github.com/mappcpd/web-services/cmd/webd/rest/router"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// Start fires up the REST server
func Start(port string) {

	// ConnectEnv to the databases
	datastore.Connect()

	// Crank up the router
	router.Start(port)
}
