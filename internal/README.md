## internal/

Contains packages for accessing and manipulating data.

Most of packages here are analogous to `models` in that they implement
structure and function for entities in the application.

The `platform/` packages provide database and other services for the
`internal` packages, the most important being the
[`datastore`](/internal/platform/datastore) package.

Basic integration tests are being included with most `internal` packages
 and the [`/testdata`](/testdata) folder contains the setup
 sql as well as some helper functions.








 
 