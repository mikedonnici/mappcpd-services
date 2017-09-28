# mailr

A utility for running sending email broadcasts based on a stored template. 

>Note this was going to be abstracted from the third-party platform however, in the interest of expediency, have used [SendGrid](https://sendgrid.com) api. This allows for each broadcast to be monitored as a 'campaign', and for users to opt out if they wish.

## How it works

A config file is read in which sets up the campaign. Depending on the options the command will do the following:
1. Authenticate with the MappCPD API to generate a token
1. Fetch all active members from MapPCPD
1. Update the SendGrid recipient list with active users
1. Fetch all recipients from SendGrid recipient list
1. Add all recipients to the SendGrid *segment* list
1. Parse the specified HTML template, and fetch content (list of resources)
1. Create a campaign at SendGrid with HTML template
1. Send a test email
1. Send the campaign at specified time   

**Todo**

* handle removing inactive members!!
* separate layout template from functions to generate content that is embedded into the template  
 

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
  "listIds": [
    1985845
  ],
  "suppressionGroupId": 4933,
  "htmlTemplate": "./cmd/mailr/template.html",
  "plainContent": "Weblink: [weblink]\r\n\r\nUnsubscribe: [unsubscribe]"
}
```

**Explanation of config options**

* authenticate
* updateMasterList
* updateSegmentList
* createCampaign
* testCampaign
* sendCampaign
* testEmail
* appendDate
* appendDateFormat
* campaignTitle
* emailSubject
* senderId
* listIds
* suppressionGroupId
* htmlTemplate
* plainContent

## Usage

Use the `-cfg` flag to specify a remote or local config file:

```bash
$ mailr -cfg https://cdn.somewhere.com/mailr/options.json 
```
