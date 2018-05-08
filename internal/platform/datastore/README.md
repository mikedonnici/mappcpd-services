## datastore/

The `datastore` package provides access to database sessions.

A `datastore` value contains fields that point to database connections -
thus far it may contain a `MySQLConnection` and a `MongoDBConnection`

A `datastore` can be passed to internal package functions so that the
package can access the data it needs. However, as the `datastore`
contains pointers to all of the databases resources the choice of how
the data is obtained is left with the internal package itself.




