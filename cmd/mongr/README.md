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

`-c` - collections to sync - 'none' does an auth check, 'all' syncs or 'only' followed by '[collectioname]' or '[collectioname] id1,id2,id3' to specify list of IDs - note no spaces in this string.

**Examples**

```bash
# Check auth (no flags)
$ mongr

# Sync everything updated within last 24 hours
$ mongr -c all

# Sync all members updated in last 7 days
$ mongr -b 7 -c only members   

# Sync single resource record, id 65473
$mongr -c only resources 65473

# Sync two member records with id 1234 and 5678
$mongr -c only members 1234,5678 
```

