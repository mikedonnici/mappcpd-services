package datastore

import (
	"fmt"
	"log"
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

	fmt.Println("Initialising MySQL connection...")

	//var m MySQLConnection
	err := MySQL.Connect() // this does not really open a new connection
	if err != nil {
		log.Fatalln("models: Could not connect to MySQLConnection server")
	}

	// Set the global db pointer
	//MySQL = m.Session
	//MySQL.Session.SetMaxIdleConns(100) // no idea if this is required

	err = MySQL.Session.Ping() // This DOES open a connection if necessary. This makes sure the db is accessible
	if err != nil {
		log.Fatalf("Error opening MySQL connection: %s", err.Error())
	}

	fmt.Println("Connected to ", MySQL.Source)
}

// connectMongoDB establishes a connection to DB2
func connectMongoDB() {

	fmt.Println("Establishing connection to MongoDB...")

	err := MongoDB.Connect()
	if err != nil {
		log.Fatalln("models: Could not connect to MongoDB server")
	}

	// Set global pointer to session for convenience
	err = MongoDB.session.Ping()
	if err != nil {
		log.Fatalf("Error opening MongoDB connection: %s\n", err.Error())
	}

	fmt.Println("Connected to ", MongoDB.Source)
}
