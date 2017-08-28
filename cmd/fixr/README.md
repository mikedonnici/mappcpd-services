# fixr

A general utility for checking and fixing various data issues. At this stage it only deals with setting and syncing the short link in the primary datastore and  the corresponding Link doc. This ensures the short link redirection service will work properly.

## Configuration

This utility accesses the datastores directly so does not require API access.

**Env vars**

```bash
# MySQL connection string
MAPPCPD_MYSQL_URL="dbuser:dbpass@tcp(db.hostname.com:3306)/dbname"

# MongoDB connection string
MAPPCPD_MONGO_URL="mongodb://mongodb.hostname.com/mongodbname"

# MongoDB database name
MAPPCPD_MONGO_DBNAME="mongodbname"

# URL for the short link (linkr) service 
MAPPCPD_SHORT_LINK_URL="https://mapp.to"

# This is a bit of a hack and will be removed at some stage, but is required to 
# prepend the record id in a short link. For example, resource with is 1234 is
# referenced by the short link service as "/r1234". The prefix was put in place
# to distinguish short links for different collections, that may have 
# overlapping id numbers. For now, just stick an "r" here.
MAPPCPD_SHORT_LINK_PREFIX="r"
```

## Usage

Use the `-b` flag to specify *backdays* - ie, how far back to include records based on `updated_at`. Default to 1.

```bash
# run fixr on records updated within the last 1 day (default)
$ fixr

# run fixr on records updated within the last 3 days
$ fixr -b 3 
```