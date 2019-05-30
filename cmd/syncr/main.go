package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/34South/envr"
	"github.com/cardiacsociety/web-services/internal/generic"
	"github.com/cardiacsociety/web-services/internal/member"
	"github.com/cardiacsociety/web-services/internal/module"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/resource"
)

// database table names
const (
	memberTable              = "member"
	memberFKColName          = "member_id" // member id field in related tables
	memberStatusTable        = "ms_m_status"
	memberTitleTable         = "ms_m_title"
	memberQualificationTable = "mp_m_qualification"
	memberPositionTable      = "mp_m_position"
	memberAccreditationTable = "mp_m_accreditation"
	memberSpecialityTable    = "mp_m_speciality"
	moduleTable              = "ol_module"
	resourceTable            = "ol_resource"
)

var memberRelatedTables = []string{
	memberStatusTable,
	memberTitleTable,
	memberAccreditationTable,
	memberQualificationTable,
	memberPositionTable,
	memberSpecialityTable,
}

// Backdays to check for updates
var backdays int

// Collection to sync
var collection string

// sql clause
var clause string

// Datastore
var store datastore.Datastore

func init() {

	envr.New("syncrEnv", []string{
		"MAPPCPD_MONGO_DBNAME",
		"MAPPCPD_MONGO_DESC",
		"MAPPCPD_MONGO_URL",
		"MAPPCPD_MYSQL_DESC",
		"MAPPCPD_MYSQL_URL",
	}).Auto()

	flag.IntVar(&backdays, "b", 0, "Specify backdays as an integer > 0")
	flag.StringVar(&collection, "c", "", "Specify what to sync - 'members', 'modules', 'resources' or 'all'")

	var err error
	store, err = datastore.FromEnv()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {

	err := flagCheck()
	if err != nil {
		log.Fatalf("flagCheck() err = %s", err)
	}
	log.Printf("Running syncr with backdays: %d on collection: %s", backdays, collection)

	clause = sqlClause()

	err = sync()
	if err != nil {
		log.Fatalf("sync() err = %s", err)
	}
}

func flagCheck() error {
	flag.Parse()
	if backdays < 1 {
		return errors.New("Backdays (-b) required, -h for help")
	}
	if collection == "" {
		return errors.New("Sync target (-c) required, -h for help")
	}
	return nil
}

// sqlClause returns an sql clause for selection of records with an updated_at
// date >= current date - backdays.
func sqlClause() string {
	// MySQL timestamp
	t := time.Now().AddDate(0, 0, -backdays).Format("2006-01-02 15:04:05")

	return fmt.Sprintf("WHERE updated_at >= '%s'", t)
}

func sync() error {

	switch collection {
	case "member", "members":
		return syncMembers()
	case "module", "modules":
		return syncModules()
	case "resource", "resources":
		return syncResources()
	case "all":
		return syncAll()
	}
	return nil
}

func syncAll() error {
	err := syncMembers()
	if err != nil {
		return err
	}
	err = syncModules()
	if err != nil {
		return err
	}
	err = syncResources()
	if err != nil {
		return err
	}
	return nil
}

func syncMembers() error {

	var count int

	// IDs of members that need to be synced
	ids, err := updateMemberIDs()
	if err != nil {
		return fmt.Errorf("syncMembers() err = %s", err)
	}

	for _, id := range ids {
		m, err := member.ByID(store, id)
		if err != nil {
			// continue if the member record is not found as it is possible for
			// an member_id value to be present in a related table but not in
			// the member table itself. So log message and continue on.
			//log.Printf("Member ID %v was not found - skipping", id)
			continue
		}
		err = m.Sync(store)
		if err != nil {
			return fmt.Errorf("syncMembers() err = %s", err)
		}
		count++
	}

	log.Printf("Sync'd %d member", count)
	return nil
}

// updateMemberIDs fetches a list of ids for members records that need to be
// updated. It searches the main member table, as well as related tables
// containing data that is included in a complete member record.
func updateMemberIDs() ([]int, error) {

	ids := []int{}

	// check member table for updates
	xi, err := generic.GetIDs(store, memberTable, clause)
	if err != nil {
		return ids, err
	}
	ids = append(ids, xi...)

	// check member-related tables for updates
	for _, t := range memberRelatedTables {
		xi, err = generic.GetIntCol(store, t, memberFKColName, clause)
		if err != nil {
			return ids, err
		}
		ids = append(ids, xi...)
	}

	return unique(ids), nil
}

func syncModules() error {

	var count int

	ids, err := generic.GetIDs(store, moduleTable, clause)
	if err != nil {
		return fmt.Errorf("syncModules() - GetIDs() err = %s", err)
	}

	for _, id := range ids {
		mod, err := module.ByID(store, id)
		if err != nil {
			// most likely error is sql: no rows error for inactive record
			// return fmt.Errorf("syncModules() - ByID() err = %s", err)
			continue
		}
		err = mod.Sync(store)
		if err != nil {
			return fmt.Errorf("syncModules() - Sync() err = %s", err)
		}
		count++
	}

	log.Printf("Sync'd %d modules", count)
	return nil
}

func syncResources() error {

	var count int

	ids, err := generic.GetIDs(store, resourceTable, clause)
	if err != nil {
		return fmt.Errorf("syncResources() err = %s", err)
	}

	for _, id := range ids {
		res, err := resource.ByID(store, id)
		if err != nil {
			// most likely error is sql: no rows error for inactive record
			//return fmt.Errorf("syncResources() err = %s", err)
			continue
		}
		err = res.Sync(store)
		if err != nil {
			return fmt.Errorf("syncResources() err = %s", err)
		}
		count++
	}

	log.Printf("Sync'd %d resources", count)
	return nil
}

// unique removes duplicates from the []int
func unique(xi []int) []int {

	// map with key and value set to the integer
	mi := make(map[int]int)
	for _, v := range xi {
		// identical values will just be overwritten
		mi[v] = v
	}

	res := []int{}
	for _, v := range mi {
		res = append(res, v)
	}

	return res
}
