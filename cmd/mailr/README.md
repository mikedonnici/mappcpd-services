# mailr

A utility for running sending email broadcasts based on a stored template. 

>Note this was going to be abstracted from the third-party platform however, in the interest of expediency, have used [SendGrid](https://sendgrid.com) api. This allows for each broadcast to be monitored as a 'campaign', and for users to opt out if they wish.

## How it works

A config file is read in which sets up the campaign. Depending on the options the command will do the following:
1. Authenticate with the MappCPD API to generate a token
1. Fetch all active members from MappCPD
1. Update the SendGrid recipient list with active users
1. Fetch all recipients from SendGrid recipient list
1. Add all recipients to the SendGrid *segment* list
1. Parse the specified HTML template, and fetch content (list of resources)
1. Create a campaign at SendGrid with HTML template
1. Send a test email
1. Send the campaign *immediately*   

## Configuration

This utility accesses the MappCPD API, so requires access to:

**Env vars**

```bash
# Admin auth credentials
MAPPCPD_ADMIN_PASS="admin-user"
MAPPCPD_ADMIN_USER"admin-pass"

# API
MAPPCPD_API_URL="https://mappcpd-api.com"

# SendGrid API key
SENDGRID_API_KEY=435243nmb245b2kj4h51kj
```

JSON config file: 

```json
{
  "authenticate": true,
  "updateMasterList": false,
  "updateSegmentList": false,
  "createCampaign": true,
  "testCampaign": false,
  "sendCampaign": false,
  "testEmail": "michael.donnici@csanz.edu.au",
  "appendDate": true,
  "appendDateFormat": "2 Jan 2006",
  "campaignTitle": "CSANZ HeartOne Update",
  "emailSubject": "CSANZ HeartOne Update",
  "senderId": 167946,
  "segmentList": 1985845,
  "suppressionGroupId": 4933,
  "htmlTemplate": "cmd/mailr/template.html",
  "plainContent": "Weblink: [weblink]\r\n\r\nUnsubscribe: [unsubscribe]",
  "backDays": 7,
  "maxContentItems": 20
}
```

The config file allows each step in the process to be switched on and off for testing.  

**authenticate**
Authenticate with the MappCPD API and get a token.
 
**updateMasterList**
Fetch all the *active* members using the MappCPD API, and update the recipient master list at SendGrid. 

**updateSegmentList**
Add all of recipients in the master list to the segment list specified in the config, in this example id 1985845.
 
**createCampaign**
Creates a campaign at SendGrid based on the specified HTML template.

**testCampaign**
Send a test to the specified email

**sendCampaign**
Send the campaign, *testCampaign* must be set to *false*
 
**testEmail**
Send the test to this email address.

**appendDate**
Append the date to the campaign title and email subject.

**appendDateFormat**
Format the appended date

**campaignTitle**
Set the campaign title at SendGrid.

**emailSubject**
Sepcify the email subject

**senderId**
Specify the SendGrid sender ID

**segmentListId**
Specify the segment list id to which the campaign will be sent. Note it is possible 
 to send to more than one segment list so this could be an array. However, in our case there is only one list so it is a single list id for simplicity.

**suppressionGroupId**
SendGrid suppression group id that will store the ids of the recipients who unsubscribe. 

**htmlTemplate**
The remote/local location of the HTML template.

**plainContent**
The (token) plain text version. 

**backDays**
When checking time-based items for inclusion, allow them to be a maximum of this many days old. 

**maxContentItems**
Maximum number of (repeating) items to include in the campaign email.


## Usage

Use the `-cfg` flag to specify a remote or local config file:

```bash
# local config
$ mailr -cfg cmd/mailr/config.json

# remote config
$ mailr -cfg https://cdn.somewhere.com/mailr/config.json 
```

## Running on Heroku
If the config file is pushed to the Heroku repo the command can be run as per the `local config` example above, that is: `mailr -cfg cmd/mailr/config.json`.

This is because the pre-compiled files are sitting on the heroku instance, along with the `bin` directory that contains the compiled executables.

To confirm this run `bash` on the heroku app:

```bash
# local machine
$ heroku login
$ heroku --app mappcpd-deployment-web-services run bash

# heroku instance
$ pwd
/app

$ ls
Procfile  README.md  bin  cmd  internal  test  vendor

$ ls cmd/mailr/
README.md  config.json  main.go  template.html
```

Having said that, it is obviously a pain to commit changes to test things out, so it is more convenient to store both the config and the html template in a remote location, and call `mailr` thus:

```bash
$ heroku run --app mappcpd-csanz-web-services mailr -cfg https://raw.githubusercontent.com/cardiacsociety/h1-cfg/master/mailr/config.json
``` 

`config.json` specifies the html template url:

```json
{
  "htmlTemplate": "https://raw.githubusercontent.com/cardiacsociety/h1-cfg/master/mailr/template.html"
}
```


## Todo 

* handle removal of inactive (in MappCPD) members
