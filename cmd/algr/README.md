# algr

MappCPD worker that leverages the [API](/cmd/webd/README.md) to maintains [Algolia](https://www.algolia.com/) search indexes.

## Configuration

**Env Vars**

```bash
# Admin auth credentials 
ADMIN_PASS="demo-admin"
ADMIN_USER="demo-pass"

# Algolia creds, with write access
ALG_API_KEY="abc......."
ALG_APP_ID="ABCDEFG..."

# Index names, setting to "OFF" will skip 
MEMBERS_INDEX="mappcpd_demo_MEMBERS"
MODULES_INDEX="mappcpd_demo_MODULES"
RESOURCES_INDEX="mappcpd_demo_RESOURCES" 

# Include records that have been modified up to this many days ago
# Set high for first run, then can be run daily with a value of '1'
# Note this refers to `updateAt` in the MongoDB doc, not `updated_at` in MySQL. 
BACK_DAYS=1

# How many at a time... 100 seems ok
BATCH_SIZE=100

# API
MAPPCPD_API_URL="https://mappcpd-api.com"
```
