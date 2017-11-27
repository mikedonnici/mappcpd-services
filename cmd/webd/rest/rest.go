package rest

import (
	"fmt"
	"github.com/mappcpd/web-services/cmd/webd/rest/router"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

func Start() {
	fmt.Println("Starting REST server...")

	// Connect to the databases
	datastore.Connect()

	// Crank up the router
	router.Start()
}
