# algr

MappCPD worker that leverages the [API](/cmd/webd/README.md) to maintains [Algolia](https://www.algolia.com/) search indexes.

## Configuration

**Env Vars**

```bash
# Admin auth credentials
MAPPCPD_ADMIN_PASS="demo-user"
MAPPCPD_ADMIN_USER"demo-pass"

# API
MAPPCPD_API_URL="https://mappcpd-api.com"

# Algolia creds - need write access
MAPPCPD_ALGOLIA_APP_ID=MZ......
MAPPCPD_ALGOLIA_API_KEY=7e2................
MAPPCPD_ALGOLIA_MEMBERS_INDEX=mappcpd_demo_MEMBERS
MAPPCPD_ALGOLIA_MODULES_INDEX=mappcpd_demo_MODULES
MAPPCPD_ALGOLIA_RESOURCES_INDEX=mappcpd_demo_RESOURCES
MAPPCPD_ALGOLIA_BATCH_SIZE=1000

# Index names, setting to "OFF" will skip 
MEMBERS_INDEX="mappcpd_demo_MEMBERS"
MODULES_INDEX="mappcpd_demo_MODULES"
RESOURCES_INDEX="mappcpd_demo_RESOURCES" 
```

## Usage

```bash
$ algr -b [int]
```

**Flags** 

`-b` - include records with a modification date up to *this many* days back - default 2

