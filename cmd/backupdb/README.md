# backupdb

Utility that fetches the latest database snapshot from Autobus and copies it to Dropbox.

AutoBus is an Heroku add-on that automatically backs up the MySQL database.

It is accessible via the Heroku web panel, in the Resources of the CSANZ admin app.

Access to Dropbox is configures via [Dropbox Apps](https://www.dropbox.com/developers/apps)

## Configuration

**Env Vars**

```bash
AUTOBUS_API_TOKEN=[from autobus web management]
DROPBOX_ACCESS_TOKEN = [Dropbox app access token]
```

## Usage

```bash
$ backupdb 
```

Note: This utility cleans up the local copy after uploading to Dropbox but does **NOT** clean up 
old backup files on Dropbox.

This utility is included in the `run_services.sh` script to ensure daily snapshot backup.





