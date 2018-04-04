package main

import (
	"flag"
	"fmt"
	"os"
	"log"

	"encoding/json"
	"github.com/pkg/errors"

	"github.com/34South/envr"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

const maxUpdateBatchCount = 1000

var memberIndexName string
var directoryIndexName string
var resourceIndexName string
var moduleIndexName string

var algoliaClient algoliasearch.Client

type Indexer interface {
	// FreshIndex creates a fresh set of objects to be indexed
	FreshIndex() ([]algoliasearch.Object, error)

	// IndexName returns the name of the algolia index
	IndexName() string
}

// flags
var collections = flag.String("c", "", "collections to sync - 'all', 'members', 'modules' or 'resources'")

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

	directoryIndexName = os.Getenv("MAPPCPD_ALGOLIA_DIRECTORY_INDEX")
	memberIndexName = os.Getenv("MAPPCPD_ALGOLIA_MEMBERS_INDEX")
	resourceIndexName = os.Getenv("MAPPCPD_ALGOLIA_RESOURCES_INDEX")
	moduleIndexName = os.Getenv("MAPPCPD_ALGOLIA_MODULES_INDEX")

	datastore.Connect()
	algoliaClient = algoliasearch.NewClient(os.Getenv("MAPPCPD_ALGOLIA_APP_ID"), os.Getenv("MAPPCPD_ALGOLIA_API_KEY"))
}

func main() {

	flag.Parse()

	switch *collections {
	case "all":
		updateMemberIndex()
		//indexResources()
		//indexModules()
	case "members":
		updateMemberIndex()
	case "resources":
		//indexResources()
	case "modules":
		//indexModules()
	default:
		fmt.Println("Unknown flag. Try -h for help.")
	}
}

func updateMemberIndex() {
	mi := NewMemberIndex(memberIndexName)
	updateIndex(&mi)
}

func updateIndex(i Indexer) {

	name := i.IndexName()
	objects, err := i.FreshIndex()
	if err != nil {
		log.Fatalln("updateIndex error:", err)
	}

	err = atomicUpdate(name, objects)
	if err != nil {
		log.Fatalln("updateIndex error:", err)
	}
}

func atomicUpdate(indexName string, objects []algoliasearch.Object) error {

	tempIndexName := indexName+"_TEMP_COPY"
	tempIndex, err := copyOfIndex(indexName, tempIndexName)
	if err != nil {
		return errors.New("Error making copy of index - " + err.Error())
	}

	err = populateIndex(tempIndex, objects)
	if err != nil {
		return errors.New("Error populating index -" + err.Error())
	}

	_, err = algoliaClient.MoveIndex(tempIndexName, indexName)
	if err != nil {
		return errors.New("Error moving temp index to target - " + err.Error())
	}

	_, err = tempIndex.Delete()
	if err != nil {
		return errors.New("Error deleting temp index - " + err.Error())
	}

	return nil
}

func copyOfIndex(sourceIndexName, destIndexName string) (algoliasearch.Index, error) {

	_, err := algoliaClient.ScopedCopyIndex(sourceIndexName, destIndexName, []string{"settings", "synonyms"})
	if err != nil {
		return nil, errors.New("Could not create temporary index for " + sourceIndexName + "-" + err.Error())
	}

	return algoliaClient.InitIndex(destIndexName), nil
}

func populateIndex(index algoliasearch.Index, objects []algoliasearch.Object) error {

	batches := determineSlices(maxUpdateBatchCount, len(objects))
	for _, b := range batches {
		batch, err := index.AddObjects(objects[b["start"]:b["end"]])
		if err != nil {
			return err
		}
		fmt.Println("Algolia batch TaskID", batch.TaskID, "- count", len(batch.ObjectIDs))
	}

	return nil
}

// determineSlices returns a set of start and end index values for iterating through a slice of totalElements in length
// and ensuring the maximum sub slice size is maxElements. It returns an array of tuples with start and end values.
// For example, determineSlices(3, 10) will return
// [{"start": 0, "end": 3}, {"start": 3, "end": 6}, {"start": 6, "end": 9}, {"start": 9, "end": 10}]
func determineSlices(maxElements, totalElements int) []map[string]int {

	var xm []map[string]int

	remainder := totalElements % maxElements
	sets := totalElements / maxElements
	if remainder > 0 {
		sets++
	}

	for c := 0; c <= totalElements; c += maxElements {
		start := c
		end := start + maxElements
		if end > totalElements {
			end = start + remainder
		}
		s := map[string]int{
			"start": start,
			"end":   end,
		}
		xm = append(xm, s)
	}

	return xm
}

// printJSON creates easy-to-read JSON representations of values for testing / debugging
func printJSON(v interface{}) {
	xb, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("outputJSON() could not marshal the value -", err)
		return
	}
	fmt.Println(string(xb))
}
