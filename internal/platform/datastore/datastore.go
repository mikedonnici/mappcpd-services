package datastore

import (
	"fmt"
	"log"

	"github.com/34South/envr"
)

// MySQL global provides access to the session value, *sql.DB, via the .Session field.
// Originally it pointed directly at *sql.DB however this was done to make it consistent with the
// MongoDB value, and to allow helper methods to be provided later on.
// So queries are accessed like this: MySQL.Session.Query()
var MySQL MySQLConnection

// MongoDB global is used for convenient access to methods for accessing MongoDB collections.
// This works slightly differently from MySQL as there are methods that return pointers
// to each collection in MongoDB, and these pointers (*Collection values) provide methods to manipulate data.
var MongoDB MongoDBConnection

func Connect() {

	connectMySQL()
	connectMongoDB()
}

// connectMySQL establishes a MySQL connection
func connectMySQL() {

	envr.New("mysqlEnv", []string{
		"MAPPCPD_MYSQL_URL",
		"MAPPCPD_MYSQL_DESC",
	}).Auto()

	err := MySQL.ConnectEnv() // this does not really open a new connection
	if err != nil {
		log.Fatalln(err)
	}

	err = MySQL.Session.Ping() // This DOES open a connection if necessary. This makes sure the db is accessible
	if err != nil {
		log.Fatalf("Error opening MySQL connection: %s\n", err.Error())
	}

	fmt.Println("datastore connected to MySQL:", MySQL.Description)
}

// connectMongoDB establishes a connection to DB2
func connectMongoDB() {

	envr.New("mongoEnv", []string{
		"MAPPCPD_MONGO_URL",
		"MAPPCPD_MONGO_DBNAME",
		"MAPPCPD_MONGO_DESC",
	}).Auto()

	err := MongoDB.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	// Set global pointer to session for convenience
	err = MongoDB.session.Ping()
	if err != nil {
		log.Fatalf("Error opening MongoDB connection: %s\n", err.Error())
	}

	fmt.Println("datastore connected to MongoDB:", MongoDB.Source)
}
