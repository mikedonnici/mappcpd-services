package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/34South/envr"
	"github.com/mappcpd/web-services/cmd/webd/graphql"
	"github.com/mappcpd/web-services/cmd/webd/rest"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

const defaultRestServerPort = "5000"
const defaultGraphQLServerPort = "5001"

func init() {
	msg := fmt.Sprint("Initialising environment...")
	env := envr.New("webdEnv", []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"MAPPCPD_API_URL",
		"MAPPCPD_JWT_TTL_HOURS",
		"MAPPCPD_JWT_SIGNING_KEY",
		"MAPPCPD_SHORT_LINK_URL",
		"MAPPCPD_SHORT_LINK_PREFIX",
		"MAPPCPD_MYSQL_DESC",
		"MAPPCPD_MYSQL_URL",
		"MAPPCPD_MONGO_DESC",
		"MAPPCPD_MONGO_DBNAME",
		"MAPPCPD_MONGO_URL",
		"WEBD_TYPE",
	}).Auto()
	if env.Ready {
		msg += "ready!"
	}
	fmt.Println(msg)
}

func main() {

	// Set the datastore from env vars
	ds, err := datastore.FromEnv()
	if err != nil {
		log.Fatalln("Could not set datastore -", err)
	}

	// Options for starting the server are varied.
	// Can only start a single web process on Heroku so use env var WEBD_TYPE to specify which type.
	// However, can override this with flags, fo local dev or if deployed on a server/container
	serverType := os.Getenv("WEBD_TYPE")
	var serverPort string

	// Optional flags to force REST / GraphQL server and port number
	serverFlag := flag.String("s", "", "Specify server type to start - 'rest' or 'graphql'")
	portFlag := flag.String("p", "", "Specify port number")
	flag.Parse()
	if strings.ToLower(*serverFlag) == "rest" {
		fmt.Println("Starting REST server...")
		serverType = "rest"
		serverPort = defaultRestServerPort
	}
	if strings.ToLower(*serverFlag) == "graphql" {
		fmt.Println("Starting GraphQL server...")
		serverType = "graphql"
		serverPort = defaultGraphQLServerPort
	}

	// Override default port numbers with optional -p flag (if set) or with env var PORT.
	// Env var must have highest precedence for Heroku
	if *portFlag != "" {
		serverPort = *portFlag
	}
	if os.Getenv("PORT") != "" {
		serverPort = os.Getenv("PORT")
	}

	if serverType == "graphql" {
		graphql.Start(serverPort, ds)
	}

	if serverType == "rest" {
		rest.Start(serverPort, ds)
	}

	// ??
	msg := "Problem starting server.\n" +
		"The required env var WEBD_TYPE is set to '%s' and should be either 'rest' or 'graphql'.\n" +
		"Alternatively, the flags -s [type] and -p [port] can be used to specify server type and port.\n" +
		"Try webd -h for help.\n" +
		"Also, make sure nothing is already listening on port %s - try this:\n" +
		"$ netstat -tulpn | grep %s\n"
	fmt.Printf(msg, os.Getenv("WEBD_TYPE"), serverPort, serverPort)
}
