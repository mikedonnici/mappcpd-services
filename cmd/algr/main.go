package main

import (
	"flag"
	"log"
	"fmt"
	"os"

	"github.com/34South/envr"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

var collections = flag.String("c", "", "collections to sync - 'all', 'directory', 'members', 'modules' or 'resources'")

var directoryIndexName string
var memberIndexName string
var moduleIndexName string
var resourceIndexName string

func init() {

	envr.New("algrEnv", []string{
		"MAPPCPD_ALGOLIA_DIRECTORY_INDEX",
		"MAPPCPD_ALGOLIA_MEMBERS_INDEX",
		"MAPPCPD_ALGOLIA_MODULES_INDEX",
		"MAPPCPD_ALGOLIA_RESOURCES_INDEX",
	}).Auto()

	datastore.Connect()
}

func main() {

	flag.Parse()

	directoryIndexName = os.Getenv("MAPPCPD_ALGOLIA_DIRECTORY_INDEX")
	memberIndexName = os.Getenv("MAPPCPD_ALGOLIA_MEMBERS_INDEX")
	resourceIndexName = os.Getenv("MAPPCPD_ALGOLIA_RESOURCES_INDEX")
	moduleIndexName = os.Getenv("MAPPCPD_ALGOLIA_MODULES_INDEX")

	switch *collections {
	case "all":
		updateDirectoryIndex()
		updateMemberIndex()
		updateModuleIndex()
		updateResourceIndex()
	case "directory":
		updateDirectoryIndex()
	case "members":
		updateMemberIndex()
	case "modules":
		updateModuleIndex()
	case "resources":
		updateResourceIndex()
	default:
		fmt.Println("Unknown flag, -h for help.")
	}
}

func updateDirectoryIndex() {

	fmt.Println("Updating directory index --------------------------------")

	if directoryIndexName == "" {
		log.Println("Directory index name is an empty string - skipping")
		return
	}

	di := newDirectoryIndex(directoryIndexName)
	if err := updateIndex(&di); err != nil {
		log.Fatalln("Error updating member index -", err)
	}
}

func updateMemberIndex() {

	fmt.Println("Updating member index  --------------------------------")

	if memberIndexName == "" {
		log.Println("Member index name is an empty string - skipping")
		return
	}

	mi := newMemberIndex(memberIndexName)
	if err := updateIndex(&mi); err != nil {
		log.Fatalln("Error updating member index -", err)
	}
}

func updateModuleIndex() {

	fmt.Println("Updating module index  --------------------------------")

	if moduleIndexName == "" {
		log.Println("Module index name is an empty string - skipping")
		return
	}

	mi := newModuleIndex(moduleIndexName)
	if err := updateIndex(&mi); err != nil {
		log.Fatalln("Error updating module index -", err)
	}
}

func updateResourceIndex() {

	fmt.Println("Updating resource index  --------------------------------")

	if resourceIndexName == "" {
		log.Println("Resource index name is an empty string - skipping")
		return
	}

	mi := newResourceIndex(resourceIndexName)
	if err := updateIndex(&mi); err != nil {
		log.Fatalln("Error updating resource index -", err)
	}
}
