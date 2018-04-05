# algr

MappCPD utility that completely rebuilds the [Algolia](https://www.algolia.com/) search indexes.

## Configuration

**Env Vars**

```bash
MAPPCPD_ALGOLIA_APP_ID = [Algolia app id]
MAPPCPD_ALGOLIA_API_KEY = [Algolia api key with write access]
MAPPCPD_ALGOLIA_DIRECTORY_INDEX = [name of member contact directory index]
MAPPCPD_ALGOLIA_MEMBERS_INDEX = [name of member admin index]
MAPPCPD_ALGOLIA_MODULES_INDEX = [name of modules index]
MAPPCPD_ALGOLIA_RESOURCES_INDEX = [name of resources index]
MAPPCPD_ALGOLIA_DIRECTORY_EXCLUDE_TITLES = [comma-sep list of titles to exclude from directory index]
```



## Usage

```bash
$ algr -c ['all', 'directory', 'members', 'modules', 'resources']
```

**Flags** 

`-c` - collection to be updated
