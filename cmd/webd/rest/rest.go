package rest

import (
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// Store represents the global datastore passed to internal packages by the handlers
var DS datastore.Datastore

// Start fires up the REST server and sets the global datastore
func Start(port string, ds datastore.Datastore) {
	DS = ds
	StartServer(port)
}
