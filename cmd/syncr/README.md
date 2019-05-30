# syncr

A data syncronisation utility that ensures the primary data records are sync'd to the document database.

Replaces the previous utility, `mongr` with the following improvements:

- Does not require API access
- Checks for updates in member _and_ related tables

## Configuration

### Env vars

This utility requires the following env vars to be set:

```bash

# MongoDB
MAPPCPD_MONGO_DBNAME="dbname"
MAPPCPD_MONGO_DESC="Mongo source description"
MAPPCPD_MONGO_URL="mongodb://mongodb.hostname.com/mongodbname"


# MySQL
MAPPCPD_MYSQL_DESC="MySQl source description"
MAPPCPD_MYSQL_URL="dbuser:dbpass@tcp(db.hostname.com:3306)/dbname"
```

## Usage

### Flags

`-b` _backdays_ - include records with `updated_at` value <= this many days ago.

`-c` _collection(s)_ to include - `members`, `modules`, `resources` or `all`

### Examples

```bash
# sync all data updated within the last 14 days
syncr -b 14 -c all

# sync member data updated within the last 24 hours
syncr -b 1 -c member
```