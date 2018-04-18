package main

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/34South/envr"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/pkg/errors"
)

const maxUpdateBatchCount = 1000

var algoliaClient algoliasearch.Client

func init() {

	envr.New("algrEnv", []string{
		"MAPPCPD_ALGOLIA_APP_ID",
		"MAPPCPD_ALGOLIA_API_KEY",
	}).Auto()

	algoliaClient = algoliasearch.NewClient(os.Getenv("MAPPCPD_ALGOLIA_APP_ID"), os.Getenv("MAPPCPD_ALGOLIA_API_KEY"))
}

type indexer interface {
	partialIndex() ([]algoliasearch.Object, error)
	fullIndex() ([]algoliasearch.Object, error)
	indexName() string
}

// update updates an index according to updateSched
func update(i indexer, ut updateType) error {

	switch ut {
	case partial:
		return partialUpdate(i)
	case full:
		return fullUpdate(i)
	case atomic:
		return atomicUpdate(i)
	}

	return errors.New("Error - update type could not be determined")
}

func partialUpdate(i indexer) error {

	objects, err := i.partialIndex()
	if err != nil {
		return err
	}

	return updateIndex(i.indexName(), objects)
}

func fullUpdate(i indexer) error {

	objects, err := i.fullIndex()
	if err != nil {
		return err
	}

	return updateIndex(i.indexName(), objects)
}

func atomicUpdate(i indexer) error {

	objects, err := i.fullIndex()
	if err != nil {
		return err
	}

	return rebuildIndex(i.indexName(), objects)
}

// updateIndex handles both partial and full updates using the objects passed in
func updateIndex(indexName string, objects []algoliasearch.Object) error {

	index := algoliaClient.InitIndex(indexName)
	err := populateIndex(index, objects)
	if err != nil {
		return errors.New("Error updating index -" + err.Error())
	}

	return nil
}

// rebuildIndex updates an index without interruption to any queries that may be in progress.
// It makes an empty copy of the original index with the same settings, and then populates the
// temporary index with fresh data. Once that is done the temporary index is moved to replace the original.
func rebuildIndex(indexName string, objects []algoliasearch.Object) error {

	tempIndexName := indexName + "_TEMP_COPY"
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

	batches := sliceBoundaries(maxUpdateBatchCount, len(objects))
	for _, b := range batches {
		batch, err := index.AddObjects(objects[b["start"]:b["end"]])
		if err != nil {
			return err
		}
		fmt.Println("Algolia batch TaskID", batch.TaskID, "- count", len(batch.ObjectIDs))
	}

	return nil
}

// sliceBoundaries returns a set of start and end index values which can be used to iterate over a slice in batches,
// where totalElements is the len of the target slice, and maxElements is the maximum size of the batch.
// For example, sliceBoundaries(3, 10) will return:
// [{"start": 0, "end": 3}, {"start": 3, "end": 6}, {"start": 6, "end": 9}, {"start": 9, "end": 10}]
// - 3 slices containing 3 elements and a final slice containing one element.
func sliceBoundaries(maxElements, totalElements int) []map[string]int {

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
		s := map[string]int{"start": start, "end": end}
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
