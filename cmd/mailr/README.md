# mailr

A utility for running sending email broadcasts based on a stored template. 

>As this is a non-system communication we *maybe* should use the Campaign Monitor API so that the user can opt out, and for stats. Alternatively, could use the SendGrid API which provides similar features. Need to check the usage limits etc. Either way, try to design this so it is not heavily BOUND to either.

## How it works

1. Each hour (or less?) mailr checks for scheduled Communications, ie with nextAt in te past. These are stored in a Comminucations collection and look like this:

```json
{
  "_id": "df234tfq....",
  "createdAt": "DATETIME",
  "updatedAt": "DATETIME",
  "nextAt": "DATETIME",
  "frequency": "days?",
  "name": "Weekly resources update",
  "description": "Notifies members of new resources added to library",
  "layoutTemplate": "<html>... {{ content }} ...</html>",
  "contentTemplate": "<div> {{ item }} </div>",
  "contentSelector": "API call that will fetch the content array",
  "recipientSelector": "API call that will fetch the list of users we need?"
 
}
```
2. The contentSelector API call is made to fetch content items - if there are none this is recorded and the job exits.

3. The recipientSelector API call is made to ensure there are recipients - if there are none this is recorded and the job exits.

4. Foreach recipient, the html content is created and then POST'd to the third-party email platform of choice. 

---


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

