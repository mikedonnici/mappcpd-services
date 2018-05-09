## datastore/

The `datastore` package provides access to database sessions.

A `datastore` value contains fields that point to database connections -
thus far it may contain a `MySQLConnection` and a `MongoDBConnection`

A `datastore` can be passed to internal package functions so that the
package can access the data it needs. However, as the `datastore`
contains pointers to all of the databases resources the choice of how
the data is obtained is left with the internal package itself.

For testing and general flexibility the datastore can be connected to
individial databases, as well as to both.

**Connect to MySQL**
```go
	ds := datastore.New()
	ds.MySQL = datastore.MySQLConnection{
		DSN:  "dev:password@tcp(localhost:3306)/mappcpd_demo",
		Desc: "Local MySQL database",
	}
	err := ds.ConnectMySQL()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connected to MySQL")
```

**Connect to MongoDB**
```go
	ds := datastore.New()
	ds.MongoDB = datastore.MongoDBConnection{
		DSN:    "mongodb://localhost/mapp_demo",
		DBName: "mapp_demo",
		Desc:   "Local MongoDB database",
	}
	err := ds.ConnectMongoDB()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connected to MongoDB")
```

**Connect to Datastore using env vars**
```go
    envr.New("testEnv", []string{
		"MAPPCPD_MYSQL_DESC",
		"MAPPCPD_MYSQL_URL",
		"MAPPCPD_MONGO_DESC",
		"MAPPCPD_MONGO_DBNAME",
		"MAPPCPD_MONGO_URL",
	}).Auto()

	ds := datastore.New()
	ds.MySQL = datastore.MySQLConnection{
		DSN:  os.Getenv("MAPPCPD_MYSQL_URL"),
		Desc: os.Getenv("MAPPCPD_MYSQL_DESC"),
	}
	ds.MongoDB = datastore.MongoDBConnection{
		DSN:    os.Getenv("MAPPCPD_MONGO_URL"),
		DBName: os.Getenv("MAPPCPD_MONGO_DBNAME"),
		Desc:   os.Getenv("MAPPCPD_MONGO_DESC"),
	}

	err := ds.ConnectAll()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Connected to Datastore")
```



