## datastore/

The `datastore` package provides access to the various data sources
for the application.

A `datastore.Datastore` value contains fields that point to database
connections - thus far it may contain a `MySQLConnection` and a `MongoDBConnection`.

The `Datastore` values is passed to the `internal` package functions
and those functions then determine which of the data sources are required
to perform the task.

For testing, or some other reason, a `Datastore` can be connected to one
source, or to all.

Examples:

**Connect to MySQL only**
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

**Connect to MongoDB only**
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

**Connect to all data sources using values from env vars**
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

In most cases within the application we need a `Datastore` that is
connected to the sources specified by the env vars, so this convenience
function can be used in place of the above:

```go
    ds, err := datastore.FromEnv()
    if err != nil {
        log.Fatalln("Could not set datastore -", err)
    }
```


