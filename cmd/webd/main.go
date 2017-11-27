package main

import (
	"fmt"
	"os"

	"github.com/34South/envr"
	"github.com/mappcpd/web-services/cmd/webd/graphql"
	"github.com/mappcpd/web-services/cmd/webd/rest"
)

func init() {
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
		"GRAPHQL_SERVER",
	}).Auto()
	if env.Ready {
		msg += "ready!"
	}
	fmt.Println(msg)
}

func main() {
	// starts the GraphQL server if env var GRAPHQL_SERVER = true,
	// otherwise will start the REST server
	if os.Getenv("GRAPHQL_SERVER") == "true" {
		graphql.Start()
	} else {
		rest.Start()
	}
}
