package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/34South/envr"
)

// campaignConfig maps to a JSON read file for setting up the email campaign
type campaignConfig struct {

	// steps to run

	// Authenticate with the MappCPD API
	Authenticate bool `json:"authenticate"`

	// Add active members to the master to the SendGrid master list of recipients
	UpdateMasterList bool `json:"updateMasterList"`

	// Update the segment list for this campaign
	UpdateSegmentList bool `json:"updateSegmentList"`

	// Create campaign at SendGrid, from the specified HTML template
	CreateCampaign bool `json:"createCampaign"`

	// Test campaign only must be false for SendCampaign to work
	TestCampaign bool `json:"testCampaign"`

	// Send campaign, todo must be true and have a scheduled date
	SendCampaign bool `json:"sendCampaign"`

	TestEmail string `json:"testEmail"`

	// Add a date to the campaign title and the email subject
	AppendDate       bool   `json:appendDate`
	AppendDateFormat string `json:"appendDateFormat"`

	// Campaign specific
	CampaignTitle      string `json:"campaignTitle"`
	EmailSubject       string `json:"emailSubject"`
	SenderID           int    `json:"senderId"`
	ListIDs            []int  `json:"listIds"`
	SuppressionGroupID int    `json:"suppressionGroupId"`
	// todo make this a url
	HTMLTemplate string `json:"htmlTemplate"`
	PlainContent string `json:"plainContent"`

	// CampaignID gets set when the campaign is created at SendGrid
	CampaignID int
}

// recipient is formatted for posting to SendGrid
type recipient struct {
	Title     string `json:"title"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// sendgridCampaign maps to the POST body for creating a new campaign via the SendGrid API
type sendgridCampaign struct {
	Title              string `json:"title"`
	Subject            string `json:"subject"`
	SenderID           int    `json:"sender_id"`
	ListIDs            []int  `json:"list_ids"`
	SuppressionGroupID int    `json:"suppression_group_id"`
	HTMLContent        string `json:"html_content"`
	PlainContent       string `json:"plain_content"`
}

// Configuration
var configFile string
var cfg = campaignConfig{}

var httpClient = &http.Client{Timeout: 30 * time.Second}
var api string
var apiAuth string
var apiActiveMembers string
var apiResources string
var token string
var recipients []recipient

func init() {

	envr.New("mongrEnv", []string{
		"MAPPCPD_ADMIN_USER",
		"MAPPCPD_ADMIN_PASS",
		"MAPPCPD_API_URL",
		"SENDGRID_API_KEY",
		"SENDGRID_LIST_ID",
		"SENDGRID_SENDER_ID",
		"SENDGRID_SUPPRESSION_GROUP_ID",
	}).Auto()

	api = os.Getenv("MAPPCPD_API_URL")
	apiAuth = api + "/v1/auth/admin"
	apiActiveMembers = api + "/v1/a/members"
	apiResources = api + "/v1/a/resources"

	flag.StringVar(&configFile, "cfg", "", "Path to JSON read file")
}

func main() {

	// First make sure we have a valid configuration
	flag.Parse()
	if configFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := cfg.read(configFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cfg.Authenticate == true {
		fmt.Print("Authenticate with MappPCD API... ")
		if err := auth(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("done")
	}

	if cfg.UpdateMasterList == true {
		fmt.Println("Update master list at SendGrid:")
		fmt.Print("Fetch active members from the MappCPD Members collection... ")
		if err := getActiveMembers(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("done")

		fmt.Print("Updating recipients master list at SendGrid... ")
		// Overwrite master list TEMPORARY!!
		recipients = []recipient{
			{Title: "Mr", FirstName: "Mike", LastName: "Donnici", Email: "michael@mesa.net.au"},
			{Title: "Prof.", FirstName: "Richmond", LastName: "Jeremy", Email: "michael.donnici@csanz.edu.au"},
		}
		if err := syncRecipients(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("done")

		t := 30
		fmt.Printf("Give SendGrid %v seconds to process the task", t)
		pause(t)
		fmt.Println()
	}

	if cfg.UpdateSegmentList == true {
		fmt.Println("Update SendGrid segment list from SendGrid master list:")
		fmt.Print("Fetch all recipient IDs from SendGrid master list... ")
		ids, err := getRecipientIDs()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("done")
		// DEBUG /////////////////////////////////////////
		// fmt.Println("found ", len(ids), "recipient ids")
		//////////////////////////////////////////////////
		if len(ids) == 0 {
			fmt.Println("No recipients in the SendGrid master list? Exiting.")
			os.Exit(1)
		}

		fmt.Print("Update segment list with recipient ids...")
		if err := updatedRecipientList(ids); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("done")

		t := 30
		fmt.Printf("Give SendGrid %v seconds to process the task", t)
		pause(t)
		fmt.Println()
	}

	if cfg.CreateCampaign == true {

		fmt.Println("Create campaign:")
		fmt.Print("Generate HTML template...")
		t, err := createTemplate()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("done")

		fmt.Print("Create SendGrid campaign with this template...")
		sc := sendgridCampaign{}
		// Copy the values from the config value
		sc.SenderID = cfg.SenderID
		sc.ListIDs = cfg.ListIDs
		sc.SuppressionGroupID = cfg.SuppressionGroupID
		sc.Title = cfg.CampaignTitle
		sc.Subject = cfg.EmailSubject
		// append time?
		if cfg.AppendDate == true {
			sc.Title += " " + time.Now().Format(cfg.AppendDateFormat)
			sc.Subject += " " + time.Now().Format(cfg.AppendDateFormat)
		}
		sc.HTMLContent = t // from above
		sc.PlainContent = cfg.PlainContent
		// Create the campaign and store the SendGrid campaign ID the campaignConfig value
		// for subsequent access
		cfg.CampaignID, err = createCampaign(sc)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("done - CampaignID is", cfg.CampaignID)
	}

	if cfg.TestCampaign == true {
		fmt.Printf("Send test for campaign ID %v to %v...", cfg.CampaignID, cfg.TestEmail)
		if err := sendTestCampaign(cfg.CampaignID); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("done")
	}
}

// read reads the config file and sets up the campaignConfig properties
func (cc *campaignConfig) read(configFile string) error {

	// read read file
	xb, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	// decode
	if err := json.Unmarshal(xb, cc); err != nil {
		return err
	}

	return nil
}

// auth fetches an admin token from the MappCPD API
func auth() error {

	type AuthResponse struct {
		Status  int
		Message string
		Data    struct {
			Token     string
			IssuedAt  time.Time
			ExpiresAt time.Time
		}
	}

	auth := AuthResponse{}
	b := `{"login": "` + os.Getenv("MAPPCPD_ADMIN_USER") + `","password": "` + os.Getenv("MAPPCPD_ADMIN_PASS") + `"}`
	res, err := httpClient.Post(apiAuth, "application/json", strings.NewReader(b))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&auth)
	token = auth.Data.Token
	if len(token) == 0 {
		return errors.New("Token has no length")
	}

	return nil
}

// getActiveMembers fetches active members from the MappCPD Members collection
func getActiveMembers() error {

	// map the individual member elemtns in the returned dta array
	type member struct {
		Contact struct {
			EmailPrimary string `json:"emailPrimary"`
		}
		Title     string `json:"title"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	// Map the overall response... only want the data field
	type memberResponse struct {
		Data []member `json:"data"`
	}
	var mr = memberResponse{}

	// selector
	b :=
		`{
			"query": {
				"active": true
			},
			"projection": {
				"_id": false,
				"title": true,
				"firstName": true,
				"lastName": true,
				"contact.emailPrimary": true
			},
			"limit": 10
		}`

	req, err := http.NewRequest("POST", apiActiveMembers, strings.NewReader(b))
	if err != nil {
		return errors.New("Problem with NewRequest() - " + err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := httpClient.Do(req)
	if err != nil {
		return errors.New("Problem making request - " + err.Error())
	}
	defer res.Body.Close()
	//printResponseBody(res.Body)

	err = json.NewDecoder(res.Body).Decode(&mr)
	if err != nil {
		return errors.New("Problem decoding response - " + err.Error())
	}

	for _, v := range mr.Data {
		r := recipient{}
		r.Title = v.Title
		r.FirstName = v.FirstName
		r.LastName = v.LastName
		r.Email = v.Contact.EmailPrimary
		recipients = append(recipients, r)
	}

	return nil
}

// syncRecipients adds the active members as recipients to SendGrid
func syncRecipients() error {

	// Marshal recipients to JSON string
	xb, err := json.MarshalIndent(recipients, "", " ")
	if err != nil {
		return errors.New("Problem marshaling recipients - " + err.Error())
	}
	// DEBUG ////////////////
	// fmt.Println(string(xb))
	////////////////////////

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/contactdb/recipients", strings.NewReader(string(xb)))
	if err != nil {
		return errors.New("Problem with NewRequest() - " + err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("SENDGRID_API_KEY"))

	// Don't need the response?
	res, err := httpClient.Do(req)
	if err != nil {
		return errors.New("Problem making request - " + err.Error())
	}
	defer res.Body.Close()

	// DEBUG - turn off otherwise can't read again
	//printResponseBody(res.Body)
	//////////////////////////////////////////////

	return nil
}

// getRecipientIDs fetches all of the recipient IDs from SendGrid
func getRecipientIDs() ([]string, error) {

	// return value
	var ids []string

	// mapp all the recipients returned from SendGrid, only need the id
	var xr struct {
		Recipients []struct {
			ID string `json:"id"`
		}
	}
	url := "https://api.sendgrid.com/v3/contactdb/recipients"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ids, errors.New("Problem with NewRequest() - " + err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("SENDGRID_API_KEY"))

	res, err := httpClient.Do(req)
	if err != nil {
		return ids, errors.New("Problem making request - " + err.Error())
	}
	defer res.Body.Close()
	// printResponseBody(res.Body)

	err = json.NewDecoder(res.Body).Decode(&xr)
	if err != nil {
		return ids, errors.New("Problem decoding response - " + err.Error())
	}
	// DEBUG ///////
	// fmt.Println(xr)
	////////////////

	for _, v := range xr.Recipients {
		ids = append(ids, v.ID)
	}

	return ids, nil
}

// updateRecipientList ensures all the recipients (ids) are on the list
func updatedRecipientList(ids []string) error {

	url := "https://api.sendgrid.com/v3/contactdb/lists/" + os.Getenv("SENDGRID_LIST_ID") + "/recipients"
	xb, err := json.Marshal(ids)
	//fmt.Println(string(xb))
	if err != nil {
		return errors.New("Problem marshaling ids - " + err.Error())
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(xb)))
	if err != nil {
		return errors.New("Problem with NewRequest() - " + err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("SENDGRID_API_KEY"))

	res, err := httpClient.Do(req)
	if err != nil {
		return errors.New("Problem making request - " + err.Error())
	}
	defer res.Body.Close()
	// DEBUG ////////////////////
	// fmt.Println(res.Status)
	////////////////////////////

	return nil
}

// createTemplate creates the HTML email template for the campaign
func createTemplate() (string, error) {

	// return string
	var h string

	tpl := template.New("layout")
	tpl, err := tpl.ParseFiles("./cmd/mailr/template.html")
	if err != nil {
		return h, errors.New("Error parsing HTML template -" + err.Error())
	}

	// selector
	sevenDaysAgo := time.Now().AddDate(0,0,-7).Format(time.RFC3339)
	b := `{
  			"find": {
				"createdAt": {
					"$gte": "` + sevenDaysAgo + `"
				}
			},
  			"select": {"_id": 0, "name": 1, "type": 1, "shortUrl": 1},
  			"limit": 30,
			"sort" : "-id"
		  }`
	//fmt.Println(b)

	req, err := http.NewRequest("POST", apiResources, strings.NewReader(b))
	if err != nil {
		return h, errors.New("Problem with NewRequest() - " + err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := httpClient.Do(req)
	if err != nil {
		return h, errors.New("Problem making request - " + err.Error())
	}
	defer res.Body.Close()
	//printResponseBody(res.Body)

	// decode response body...
	var data = struct {
		Data []struct {
			ResourceType string `json:"type"`
			Name         string `json:"name"`
			Link         string `json:"shortUrl"`
		} `json:"data"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return h, err
	}

	var t bytes.Buffer
	if err := tpl.ExecuteTemplate(&t, "template.html", data.Data); err != nil {
		return h, err
	}

	return t.String(), nil
}

// createCampaign creates an email campaign at SendGrid
func createCampaign(campaign sendgridCampaign) (int, error) {

	// return campaign id
	var id int

	url := "https://api.sendgrid.com/v3/campaigns"

	b, err := json.Marshal(campaign)
	if err != nil {
		return id, errors.New("Problem marshaling campaign - " + err.Error())
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(b)))
	if err != nil {
		return id, errors.New("Problem with NewRequest() - " + err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("SENDGRID_API_KEY"))

	res, err := httpClient.Do(req)
	if err != nil {
		return id, errors.New("Problem making request - " + err.Error())
	}
	defer res.Body.Close()
	// printResponseBody(res.Body)

	// map response
	var data = struct {
		ID int `json:"id"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return id, err
	}

	return data.ID, nil
}

func sendTestCampaign(campaignID int) error {

	url := "https://api.sendgrid.com/v3/campaigns/" + strconv.Itoa(campaignID) + "/schedules/test"
	b := `{"to": "` + cfg.TestEmail + `"}`

	req, err := http.NewRequest("POST", url, strings.NewReader(b))
	if err != nil {
		return errors.New("Problem with NewRequest() - " + err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("SENDGRID_API_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Problem making request - " + err.Error())
	}
	defer res.Body.Close()
	// printResponseBody(res.Body)

	return nil
}

// printResponseBody is a utility function to print the response from API requests
func printResponseBody(body io.ReadCloser) {
	xb, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("RESPONSE BODY -----------------------------------------")
	fmt.Println(string(xb))
	fmt.Println("END RESPONSE BODY -----------------------------------------")
}

// pause waits for a number of seconds. This is required as there is a delay between sending updates
// to SendGrid and the data actually being available.
func pause(s int) {
	for i := 0; i < s; i++ {
		duration := time.Second
		time.Sleep(duration)
		fmt.Print(".")
	}
}
