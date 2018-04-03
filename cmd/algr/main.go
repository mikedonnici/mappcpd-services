package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
	"encoding/json"

	"github.com/34South/envr"
	"github.com/mappcpd/web-services/internal/members"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/resources"
	"github.com/mappcpd/web-services/internal/modules"
)

var api string
var apiAuth string
var apiMembers string
var membersIndex string
var directoryIndex string
var apiResources string
var resourcesIndex string
var apiModules string
var modulesIndex string

const tempNameSuffix = "_TEMP_COPY"

var maxBatch int

var token string

type AuthRequest struct {
	Status  int
	Message string
	Data    AuthData
}

type AuthData struct {
	Token     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type Index struct {
	Name     string
	TempName string
	Data     []map[string]interface{}
}

// flags
var collections = flag.String("c", "", "collections to sync - 'all', 'members', 'modules' or 'resources'")

// backDate is a string date in format "2017-01-21T13:35:30+10:00" (RFC3339) so we can pass it to the API
var backDate string

func init() {
	envr.New("algrEnv", []string{
		"MAPPCPD_ALGOLIA_APP_ID",
		"MAPPCPD_ALGOLIA_API_KEY",
		"MAPPCPD_ALGOLIA_BATCH_SIZE",
		"MAPPCPD_ALGOLIA_DIRECTORY_INDEX",
		"MAPPCPD_ALGOLIA_MEMBERS_INDEX",
		"MAPPCPD_ALGOLIA_MODULES_INDEX",
		"MAPPCPD_ALGOLIA_RESOURCES_INDEX",
		"MAPPCPD_ALGOLIA_DIRECTORY_EXCLUDE_TITLES",
	}).Auto()

	directoryIndex = os.Getenv("MAPPCPD_ALGOLIA_DIRECTORY_INDEX")
	membersIndex = os.Getenv("MAPPCPD_ALGOLIA_MEMBERS_INDEX")
	resourcesIndex = os.Getenv("MAPPCPD_ALGOLIA_RESOURCES_INDEX")
	modulesIndex = os.Getenv("MAPPCPD_ALGOLIA_MODULES_INDEX")

	var err error
	maxBatch, err = strconv.Atoi(os.Getenv("MAPPCPD_ALGOLIA_BATCH_SIZE"))
	if err != nil {
		log.Fatalln("Could not set batch size:", err)
	}

	datastore.Connect()
}

func main() {

	flag.Parse()

	switch *collections {
	case "all":
		indexMembers()
		indexResources()
		indexModules()
	case "members":
		indexMembers()
	case "resources":
		indexResources()
	case "modules":
		indexModules()
	default:
		fmt.Println("Unknown flag. Try -h for help.")
	}
}

func indexMembers() {

	fmt.Println("Indexing member records... ")
	if strings.ToLower(membersIndex) == "off" {
		fmt.Println("... member index is set to 'OFF' - nothing to do")
		return
	}

	// Two indexes to update (members and directory) - can use the same query for both.
	// The reshape function can be used to filter out records not suitable for the directory.
	// In all cases we only want members with a membership record
	// mongo shell query is: db.Members.find({"memberships.title": {$exists : true}})
	fmt.Println("... fetching member docs updated since", backDate)
	query := bson.M{"memberships.title": bson.M{"$exists": true}}
	members, err := members.FetchMembers(query, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("... creating index")
	mi := Index{
		Name:     membersIndex,
		TempName: membersIndex + tempNameSuffix,
	}
	mi.Data = createMemberIndex(members)
	fmt.Println("... updating Algolia index:", os.Getenv("MAPPCPD_ALGOLIA_MEMBERS_INDEX"))
	mi.atomicUpdate()

	fmt.Println("... creating index")
	di := Index{
		Name:     directoryIndex,
		TempName: directoryIndex + tempNameSuffix,
	}
	di.Data = createDirectoryIndex(members)
	fmt.Println("... updating Algolia index:", os.Getenv("MAPPCPD_ALGOLIA_DIRECTORY_INDEX"))
	di.atomicUpdate()
}

// indexResources manages the resources index. Note the Mongo collection now has a bool field called 'active'
// which mirrors the flag (0/1) field used in MySQL for soft deletes. So any record with active=false should be
// removed from the index.
func indexResources() {

	fmt.Println("Indexing resource records... ")
	if strings.ToLower(resourcesIndex) == "off" {
		fmt.Println("... resource index is set to OFF - nothing to do")
		return
	}

	fmt.Println("... fetching resource docs updated since", backDate)
	query := bson.M{"active": true, "primary": true}
	resources, err := resources.FetchResources(query, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("... creating index")
	i := Index{
		Name:     resourcesIndex,
		TempName: resourcesIndex + tempNameSuffix,
	}
	i.Data = createResourceIndex(resources)

	fmt.Println("... updating Algolia index:", os.Getenv("MAPPCPD_ALGOLIA_RESOURCES_INDEX"))
	i.atomicUpdate()
}

func indexModules() {

	fmt.Println("Indexing module records... ")

	if strings.ToLower(modulesIndex) == "off" {
		fmt.Println("... modules index is set to 'OFF' - nothing to do")
		return
	}

	fmt.Println("... fetching module docs updated since", backDate)
	query := bson.M{"current": true}
	modules, err := modules.FetchModules(query, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("... creating index")
	i := Index{
		Name:     modulesIndex,
		TempName: modulesIndex + tempNameSuffix,
	}
	i.Data = createModuleIndex(modules)
	fmt.Println("... updating Algolia index:", os.Getenv("MAPPCPD_ALGOLIA_MODULES_INDEX"))
	i.atomicUpdate()
}

// createResourceIndex creates a json document suitable for the Algolia resources index
func createResourceIndex(resources []resources.Resource) []map[string]interface{} {

	var resourceIndex []map[string]interface{}

	for _, resource := range resources {

		pubDate := resource.PubDate.Date.Format(time.RFC3339)
		pubStamp := resource.PubDate.Date.Unix()

		r := map[string]interface{}{
			"_id":                  resource.OID,
			"id":                   resource.ID,
			"createdAt":            resource.CreatedAt,
			"updatedAt":            resource.UpdatedAt,
			"publishedAt":          pubDate,
			"publishedAtTimestamp": pubStamp,
			"type":                 resource.Type,
			"name":                 resource.Name,
			"description":          resource.Description,
			"keywords":             resource.Keywords,
			"shortUrl":             resource.ShortURL,
			"resourceUrl":          resource.ResourceURL,
		}

		// Pubmed Attributes
		attributes := resource.Attributes

		v, ok := attributes["sourceId"]
		if ok {
			r["sourceId"] = v
		}

		v, ok = attributes["sourceName"]
		if ok {
			r["sourceName"] = v
		}

		v, ok = attributes["sourceNameAbbrev"]
		if ok {
			r["sourceNameAbbrev"] = v
		}

		v, ok = attributes["sourcePubDate"]
		if ok {
			r["sourcePubDate"] = v
		}

		v, ok = attributes["sourceVolume"]
		if ok {
			r["sourceVolume"] = v
		}

		v, ok = attributes["sourceIssue"]
		if ok {
			r["sourceIssue"] = v
		}

		v, ok = attributes["sourcePages"]
		if ok {
			r["sourcePages"] = v
		}

		resourceIndex = append(resourceIndex, r)
	}

	return resourceIndex
}

// createMemberIndex creates a json document suitable for the Algolia member index
func createMemberIndex(members []members.Member) []map[string]interface{} {

	var memberIndex []map[string]interface{}

	for _, member := range members {

		// concat name fields
		name := fmt.Sprintf("%s %s %s", member.Title, member.FirstName, member.LastName)

		// personal contact details
		email := member.Contact.EmailPrimary
		mobile := member.Contact.Mobile

		// only use location info from the directory contact record, and only the general locality
		var location string
		for _, l := range member.Contact.Locations {
			if l.Description == "Directory" {
				location = fmt.Sprintf("%s %s %s %s", l.City, l.State, l.Postcode, l.Country)
			}
		}

		// Membership title - dig into the memberships array even though there is only one.
		membership := member.Memberships[0].Title

		// Specialities
		var specialities []string
		for _, s := range member.Specialities {
			specialities = append(specialities, s.Name)
		}

		// Affiliations (Positions) with certain groups. Only include positions with no end, or a future end date
		var affiliations []string
		for _, p := range member.Positions {
			endDate, err := time.Parse("2006-01-02", p.End)
			if err != nil || endDate.After(time.Now()) {
				affiliations = append(affiliations, p.OrgName)
			}
		}

		m := map[string]interface{}{
			"_id":          member.OID,
			"id":           member.ID,
			"active":       member.Active,
			"name":         name,
			"email":        email,
			"mobile":       mobile,
			"location":     location,
			"membership":   membership,
			"affiliations": affiliations,
			"specialities": specialities,
		}

		memberIndex = append(memberIndex, m)
	}

	return memberIndex
}

// createDirectoryIndex creates a json document suitable for the Algolia directory index
func createDirectoryIndex(members []members.Member) []map[string]interface{} {

	var directoryIndex []map[string]interface{}

	for _, member := range members {

		// concat name fields
		name := fmt.Sprintf("%s %s %s", member.Title, member.FirstName, member.LastName)

		if excludeMemberFromDirectory(member) {
			fmt.Println(" - excluding", name)
			continue
		}

		// personal contact details
		email := member.Contact.EmailPrimary
		mobile := member.Contact.Mobile

		// only use location info from the directory contact record, and only the general locality
		var location string
		for _, l := range member.Contact.Locations {
			if l.Description == "Directory" {
				location = fmt.Sprintf("%s %s %s %s", l.City, l.State, l.Postcode, l.Country)
			}
		}

		// Membership title - dig into the memberships array even though there is only one.
		membership := member.Memberships[0].Title

		// Specialities
		var specialities []string
		for _, s := range member.Specialities {
			specialities = append(specialities, s.Name)
		}

		// Affiliations (Positions) with certain groups. Only include positions with no end, or a future end date
		var affiliations []string
		for _, p := range member.Positions {
			endDate, err := time.Parse("2006-01-02", p.End)
			if err != nil || endDate.After(time.Now()) {
				affiliations = append(affiliations, p.OrgName)
			}
		}

		m := map[string]interface{}{
			"_id":          member.OID,
			"id":           member.ID,
			"active":       member.Active,
			"name":         name,
			"email":        email,
			"mobile":       mobile,
			"location":     location,
			"membership":   membership,
			"affiliations": affiliations,
			"specialities": specialities,
		}

		directoryIndex = append(directoryIndex, m)
	}

	return directoryIndex
}

// createModuleIndex creates a json document suitable for the Algolia modules index
func createModuleIndex(modules []modules.Module) []map[string]interface{} {

	var moduleIndex []map[string]interface{}

	for _, module := range modules {

		m := map[string]interface{}{
			"_id":         module.OID,
			"id":          module.ID,
			"createdAt":   module.CreatedAt,
			"updateAt":    module.UpdatedAt,
			"publishedAt": module.PublishedAt,
			"name":        module.Name,
			"description": module.Description,
			"started":     module.Started,
			"finished":    module.Finished,
		}

		moduleIndex = append(moduleIndex, m)
	}

	return moduleIndex
}

// excludeMemberFromDirectory returns true if member record should be excluded from the directory
func excludeMemberFromDirectory(member members.Member) bool {

	// no active status value
	if member.Active != true {
		fmt.Print("Member inactive or status unknown")
		return true
	}

	// No directory consent
	if member.Contact.Directory != true {
		fmt.Print("Member has not consented to directory listing")
		return true
	}

	// no membership title
	if member.Title == "" {
		fmt.Print("No membership title")
		return true
	}

	// membership title in exclude list
	xs := strings.Split(os.Getenv("MAPPCPD_ALGOLIA_DIRECTORY_EXCLUDE_TITLES"), ",")
	title := strings.ToLower(member.Title)
	for _, s := range xs {
		excludeTitle := strings.ToLower(strings.TrimSpace(s))
		if title == excludeTitle {
			fmt.Printf("title '%s' matches exclude list value '%s'", title, excludeTitle)
			return true
		}
	}

	// include in directory
	return false
}

// atomicUpdate creates a temporary index, pushes the data to the temporary index, and then moves the temporary (new) index
// to the original index. The last step is atomic and there is no down time or interruption to queries.
// ref: https://www.algolia.com/doc/tutorials/indexing/synchronization/atomic-reindexing/
func (i Index) atomicUpdate() {

	client := algoliasearch.NewClient(os.Getenv("MAPPCPD_ALGOLIA_APP_ID"), os.Getenv("MAPPCPD_ALGOLIA_API_KEY"))

	// temp index
	_, err := client.ScopedCopyIndex(i.Name, i.TempName, []string{"settings", "synonyms"})
	if err != nil {
		fmt.Println("Could not create temporary index for", i.Name)
		os.Exit(1)
	}

	tempIndex := client.InitIndex(i.TempName)

	batchSize := maxBatch
	for j := 0; j < len(i.Data); j++ {

		// Reset this or we accumulate!!!
		objects := []algoliasearch.Object{}

		// If remaining items is less than batch size...
		if len(i.Data)-j < batchSize {
			batchSize = len(i.Data) - j
		}

		fmt.Println("--- next batch of", batchSize)
		for c := 0; c < batchSize; c++ {
			// set algolia objectID to _id
			i.Data[j]["objectID"] = i.Data[j]["_id"]
			objects = append(objects, i.Data[j])
			// Don't increment j on the last loop because the outer loop
			// also increments it so we end up missing a value
			if c < (batchSize - 1) {
				j++
			}
		}

		batch, err := tempIndex.AddObjects(objects)
		if err != nil {
			fmt.Println("Error indexing batch -", err)
			os.Exit(1)
		}
		fmt.Println("Algolia batch taskID:", batch.TaskID, "completed indexing of", len(batch.ObjectIDs), "objects for index:", i.TempName)
	}

	// Move temp index into place then delete it
	_, err = client.MoveIndex(i.TempName, i.Name)
	if err != nil {
		fmt.Println("Error moving temp index -", err)
		os.Exit(1)
	}

	_, err = tempIndex.Delete()
	if err != nil {
		fmt.Println("Error deleting temp index -", err)
		os.Exit(1)
	}
}

// indexBatch adds a batch of items to an index. Note that this seems to time out at somewhere over 1000 objects
// regardless of batch size, timeout settings or anything else.
func indexBatch(xo []algoliasearch.Object, indexName string) {

	client := algoliasearch.NewClient(os.Getenv("MAPPCPD_ALGOLIA_APP_ID"), os.Getenv("MAPPCPD_ALGOLIA_API_KEY"))
	index := client.InitIndex(indexName)
	batch, err := index.AddObjects(xo)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Algolia taskID:", batch.TaskID, "completed indexing of", len(batch.ObjectIDs), "objects for index:", indexName)
}

func timeStampFromDate(date string) int64 {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		fmt.Println("Error parsing date string", err)
	}
	return t.Unix()
}

// outputJSON creates easy-to-read JSON representations of values for testing / debugging
func outputJSON(v interface{}) {
	xb, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("outputJSON() could not marshal the value -", err)
		return
	}
	fmt.Println(string(xb))
}
