# mongr

A MappCPD worker that leverages the [API](/cmd/webd/README.md) to create denormalized records (JSON) from the primary MySQL database, and synchronise them to MongoDB.
                      
This serves two main purposes:

1. Provide faster search operations for the API
2. Provide index documents for the [algr](/cmd/algr/README.md) service

## Configuration

**Env vars**

```bash
# Admin auth credentials
MAPPCPD_ADMIN_PASS="demo-user"
MAPPCPD_ADMIN_USER"demo-pass"

# API
MAPPCPD_API_URL="https://mappcpd-api.com"
```

## Usage

```bash
$ mongr -b [int] -c [collection]
```

**Flags** 

`-b` - include records that were updated up to this many days back, default '1'

`-c` - collections to sync - 'none' does an auth check, 'all' syncs or 'only' followed by 'name:all' or 'name:id1,id2,id3' to be more specific.

**Examples**

```bash
# Check auth (no flags)
$ mongr

# Sync everything updated within last 24 hours
$ mongr -c all

# Sync members updated in last 7 days
$ mongr -c members:all -b 7  

# Sync a single resource record, id 65473
$mongr -c resources:65473 
```

