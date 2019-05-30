package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/34South/envr"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
)

// updateSched is a date flag used to determine the updateSched for a particular index
type updateSched int

const (
	none updateSched = iota
	daily
	weekly
	monthly
)

var updateScheds = []string{"none", "daily", "weekly", "monthly"}

// updateType specifies the level of index rebuild - atomic, complete or partial
type updateType int

const (
	partial updateType = iota
	full
	atomic
)

var updateTypes = []string{"partial", "full", "atomic"}

var collections = flag.String("c", "", "collections to sync - 'all', 'directory', 'members', 'modules', 'resources', 'qualifications', 'organisations'")

var directoryIndexName string
var memberIndexName string
var moduleIndexName string
var resourceIndexName string
var qualificationIndexName string
var organisationIndexName string

var sched = scheduledUpdateType()

var DS datastore.Datastore

func init() {

	envr.New("algrEnv", []string{
		"MAPPCPD_ALGOLIA_DIRECTORY_INDEX",
		"MAPPCPD_ALGOLIA_MEMBERS_INDEX",
		"MAPPCPD_ALGOLIA_MODULES_INDEX",
		"MAPPCPD_ALGOLIA_RESOURCES_INDEX",
		"MAPPCPD_ALGOLIA_QUALIFICATIONS_INDEX",
		"MAPPCPD_ALGOLIA_ORGANISATIONS_INDEX",
		"MAPPCPD_MONGO_DBNAME",
		"MAPPCPD_MONGO_DESC",
		"MAPPCPD_MONGO_URL",
		"MAPPCPD_MYSQL_DESC",
		"MAPPCPD_MYSQL_URL",
	}).Auto()

	var err error
	DS, err = datastore.FromEnv()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {

	flag.Parse()

	directoryIndexName = os.Getenv("MAPPCPD_ALGOLIA_DIRECTORY_INDEX")
	memberIndexName = os.Getenv("MAPPCPD_ALGOLIA_MEMBERS_INDEX")
	resourceIndexName = os.Getenv("MAPPCPD_ALGOLIA_RESOURCES_INDEX")
	moduleIndexName = os.Getenv("MAPPCPD_ALGOLIA_MODULES_INDEX")
	qualificationIndexName = os.Getenv("MAPPCPD_ALGOLIA_QUALIFICATIONS_INDEX")
	organisationIndexName = os.Getenv("MAPPCPD_ALGOLIA_ORGANISATIONS_INDEX")

	switch *collections {
	case "all":
		updateDirectoryIndex()
		updateMemberIndex()
		updateModuleIndex()
		updateResourceIndex()
		updateQualificationIndex()
		updateOrganisationIndex()
	case "directory":
		updateDirectoryIndex()
	case "members":
		updateMemberIndex()
	case "modules":
		updateModuleIndex()
	case "resources":
		updateResourceIndex()
	case "qualifications":
		updateQualificationIndex()
	case "organisations":
		updateOrganisationIndex()
	default:
		fmt.Println("Unknown flag, -h for help.")
	}
}

func updateDirectoryIndex() {

	if directoryIndexName == "" {
		log.Println("Directory index name is an empty string - skipping")
		return
	}

	var ut updateType
	switch sched {
	case monthly, weekly:
		ut = atomic
	case daily:
		ut = full
	default:
		ut = partial
	}
	updateLogMessage(directoryIndexName, ut)

	di := newDirectoryIndex(directoryIndexName)
	if err := update(&di, ut); err != nil {
		log.Fatalln("Error updating member index -", err)
	}
}

func updateMemberIndex() {

	if memberIndexName == "" {
		log.Println("Member index name is an empty string - skipping")
		return
	}

	var ut updateType
	switch sched {
	case monthly:
		ut = atomic
	case weekly, daily:
		ut = full
	default:
		ut = partial
	}
	updateLogMessage(memberIndexName, ut)

	mi := newMemberIndex(memberIndexName)
	if err := update(&mi, ut); err != nil {
		log.Fatalln("Error updating member index -", err)
	}
}

func updateModuleIndex() {

	if moduleIndexName == "" {
		log.Println("Module index name is an empty string - skipping")
		return
	}

	ut := atomic
	updateLogMessage(moduleIndexName, ut)

	mi := newModuleIndex(moduleIndexName)
	if err := update(&mi, ut); err != nil {
		log.Fatalln("Error updating module index -", err)
	}
}

func updateResourceIndex() {

	if resourceIndexName == "" {
		log.Println("Resource index name is an empty string - skipping")
		return
	}

	var ut updateType
	switch sched {
	case monthly:
		ut = atomic
	default:
		ut = partial
	}
	updateLogMessage(resourceIndexName, ut)

	mi := newResourceIndex(resourceIndexName)
	if err := update(&mi, ut); err != nil {
		log.Fatalln("Error updating resource index -", err)
	}
}

func updateQualificationIndex() {

	if qualificationIndexName == "" {
		log.Println("Qualifications index name is an empty string - skipping")
		return
	}

	// Always atomic because this index has no mongo collection and hence no Object IDs. Any other type
	// of index update will result in duplicate records.
	var ut updateType
	switch sched {
	default:
		ut = atomic
	}
	updateLogMessage(qualificationIndexName, ut)

	i := newQualificationIndex(qualificationIndexName)
	if err := update(&i, ut); err != nil {
		log.Fatalln("Error updating qualifications index -", err)
	}
}

func updateOrganisationIndex() {

	if organisationIndexName == "" {
		log.Println("Organisation index name is an empty string - skipping")
		return
	}

	// Always atomic because this index has no mongo collection and hence no Object IDs. Any other type
	// of index update will result in duplicate records.
	var ut updateType
	switch sched {
	default:
		ut = atomic
	}
	updateLogMessage(organisationIndexName, ut)

	i := newOrganisationIndex(organisationIndexName)
	if err := update(&i, ut); err != nil {
		log.Fatalln("Error updating organisation index -", err)
	}
}

// scheduledUpdateType returns an updateSched constant based on the date
func scheduledUpdateType() updateSched {

	// first of the month
	if time.Now().Day() == 1 {
		return monthly
	}

	// First day of the week
	if time.Now().Weekday().String() == "Sunday" {
		return weekly
	}

	// limited rebuild otherwise
	return daily
}

func updateLogMessage(in string, ut updateType) {
	log.Printf("Updating index: %s, sched: %s, type: %s", in, updateScheds[sched], updateTypes[ut])
}
