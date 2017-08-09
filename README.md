# MappCPD Web Services

[![Build Status](https://travis-ci.org/mappcpd/web-services.svg?branch=master)](https://travis-ci.org/mappcpd/web-services)

Set of web services supporting [MappCPD](https://mappcpd.com).

* [cmd/](/cmd/README.md) - executable packages
  * [pubmedr/](/cmd/pubmedr/README.md) - worker to fetch pubmed articles
  * [mongr/](/cmd/mongr/README.md) - worker to sync data from MySQL to MongoDB
  * [algr/](/cmd/algr/README.md) - worker to sync Algolia indexes
  * [webd/](/cmd/webd/README.md) - web server for API
* [internal/](/internal/README.md) - internal packages
* [vendor/](/vendor/README.md) - vendor packages

## Configuration

The following configuration includes env vars required for all of the services, set using a `.env` file or by some other method.

```bash
# API Base URL (don't forget PORT number if required)
MAPPCPD_API_URL="https://mappcpd-api.com"


# Data stores ------------------------------------------------------------------

# MySQL connection string
MAPPCPD_MYSQL_URL="dbuser:dbpass@tcp(db.hostname.com:3306)/dbname"

# MySQL db descriptive string - used to ensure data is coming from the 
# right place in dev, staging, production! 
MAPPCPD_MYSQL_DESC="Production server"

# MongoDB connection string
MAPPCPD_MONGO_URL="mongodb://mongodb.hostname.com/mongodbname"

# MongoDB database name
MAPPCPD_MONGO_DBNAME="mongodbname"

# MongoDB descriptive string, for same reason as above ;)
MAPPCPD_MONGO_DESC="Production server"


# AWS S3 -----------------------------------------------------------------------

# The API currently plays a small role in facilitating direct uploads 
# from the client to S3, thus bypassing the server. To do this the API 
# issues a signed url to the client requiring the following credntials:
AWS_REGION="ap-southeast-1"
AWS_ACCESS_KEY_ID="ABC....4546QRST"
AWS_SECRET_ACCESS_KEY="fghjkl...5678asdfg"


# Worker Services: pubmedr, mongr, algr ----------------------------------------

# Admin creds to access the API
MAPPCPD_ADMIN_USER="admin-user"
MAPPCPD_ADMIN_PASS="admin-pass"

# Pubmedr ----
# URL or relative file path location of the pubmedr query config (JSON) 
MAPPCPD_PUBMED_BATCH_FILE="https://some-cloud.io/pubmed.json"
# Number of articles to return in each batch
MAPPCPD_PUBMED_RETMAX=200

# Algr ----
# Algolia API - must have write access 
MAPPCPD_ALGOLIA_APP_ID="MZQPVRPXFY"
MAPPCPD_ALGOLIA_API_KEY="7e25770a12191493af086d4a03ee4acb"
# Names of the relevant indexes 
MAPPCPD_ALGOLIA_MEMBERS_INDEX="mappcpd_members"
MAPPCPD_ALGOLIA_MODULES_INDEX="mappcpd_modules"
MAPPCPD_ALGOLIA_RESOURCES_INDEX="mappcpd_resources"
# Days back to check for updated records
MAPPCPD_ALGOLIA_BACK_DAYS=1
# Batch size for index update - 1000 seems to work fine
MAPPCPD_ALGOLIA_BATCH_SIZE=1000


# linkr - short link redirector and stats --------------------------------------

# The short link redirector is run as a separate service however 
# the API does set and retrieve information relating to it service 

# URL for the short link (linkr) service 
MAPPCPD_SHORT_LINK_URL="https://mapp.to"

# This is a bit of a hack and will be removed at some stage, but is required to 
# prepend the record id in a short link. For example, resource with is 1234 is
# referenced by the short link service as "/r1234". The prefix was put in place
# to distinguish short links for different collections, that may have 
# overlapping id numbers. For now, just stick an "r" here.
MAPPCPD_SHORT_LINK_PREFIX="r"

``` 

  





## Service Architecture

![resources](https://docs.google.com/drawings/d/1zJ4pQCb94syzpCvoqRBXwbMUvs8LhpFlFE2Gax6LTfM/pub?w=691&h=431)

Note that [linkr](https://github.com/34South/linkr) short link redirection service is a separate service.

See [MappCPD Architecture](https://github.com/mappcpd/architecture/wiki) for more info.



## References

Project structure based on Bill Kennedy's [package oriented design](https://www.goinggo.net/2017/02/package-oriented-design.html).

