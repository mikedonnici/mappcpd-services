package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/34South/envr"
	"github.com/cardiacsociety/web-services/cmd/webd/server"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

const defaultServerPort = "5000"

func init() {
	msg := fmt.Sprint("Initialising environment... ")
	env := envr.New("webdEnv", []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SES_REGION",
		"AWS_SES_ACCESS_KEY_ID",
		"AWS_SES_SECRET_ACCESS_KEY",
		"MAILGUN_DOMAIN",
		"MAILGUN_API_KEY",
		"MAPPCPD_API_URL",
		"MAPPCPD_JWT_TTL_HOURS",
		"MAPPCPD_JWT_SIGNING_KEY",
		"MAPPCPD_MYSQL_DESC",
		"MAPPCPD_MYSQL_URL",
		"MAPPCPD_MONGO_DESC",
		"MAPPCPD_MONGO_DBNAME",
		"MAPPCPD_MONGO_URL",
		"MAPPCPD_MX_SERVICE",
		"MAPPCPD_SHORT_LINK_URL",
		"MAPPCPD_SHORT_LINK_PREFIX",
		"SENDGRID_API_KEY",
	}).Auto()
	if env.Ready {
		msg += "done."
	}

	log.Println(msg)
}

func main() {

	// Set the datastore from env vars
	ds, err := datastore.FromEnv()
	if err != nil {
		log.Fatalln("Could not set datastore -", err)
	}

	// Override default port numbers with optional -p flag (if set) or with env var PORT.
	var serverPort = defaultServerPort
	portFlag := flag.String("p", "", "Override default port")
	flag.Parse()
	if *portFlag != "" {
		serverPort = *portFlag
	}
	// Env var must have highest precedence for Heroku
	if os.Getenv("PORT") != "" {
		serverPort = os.Getenv("PORT")
	}

	// Server Handlers
	h := server.Router(ds)
	log.Printf("Starting web services on port %s", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, h))
}
