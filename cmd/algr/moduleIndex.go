package main

import (
	"time"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/cardiacsociety/web-services/internal/module"
	"gopkg.in/mgo.v2/bson"
)

type moduleIndex struct {
	Name      string
	RawData   []module.Module
	IndexData []algoliasearch.Object
	Error     error
}

// newModuleIndex returns a pointer to moduleIndex value initialised with the index name
func newModuleIndex(name string) moduleIndex {
	return moduleIndex{
		Name: name,
	}
}

func (mi *moduleIndex) indexName() string {
	return mi.Name
}

func (mi *moduleIndex) partialIndex() ([]algoliasearch.Object, error) {
	mi.fetchLimitedData()
	mi.createIndexObjects()
	return mi.IndexData, mi.Error
}

func (mi *moduleIndex) fullIndex() ([]algoliasearch.Object, error) {
	mi.fetchAllData()
	mi.createIndexObjects()
	return mi.IndexData, mi.Error
}

func (mi *moduleIndex) fetchLimitedData() {
	timeBack := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	query := bson.M{"current": true, "updatedAt": bson.M{"$gte": timeBack}}
	mi.RawData, mi.Error = module.FetchModules(DS, query, 0)
}

func (mi *moduleIndex) fetchAllData() {
	query := bson.M{"current": true}
	mi.RawData, mi.Error = module.FetchModules(DS, query, 0)
}

func (mi *moduleIndex) createIndexObjects() {
	for i := range mi.RawData {
		mi.IndexData = append(mi.IndexData, algoliasearch.Object{})
		mi.createObject(i)
	}
}

func (mi *moduleIndex) createObject(i int) {

	module := mi.RawData[i]

	mi.IndexData[i] = make(map[string]interface{})
	mi.IndexData[i]["objectID"] = module.OID
	mi.IndexData[i]["_id"] = module.OID
	mi.IndexData[i]["id"] = module.ID
	mi.IndexData[i]["createdAt"] = module.CreatedAt
	mi.IndexData[i]["updateAt"] = module.UpdatedAt
	mi.IndexData[i]["publishedAt"] = module.PublishedAt
	mi.IndexData[i]["name"] = module.Name
	mi.IndexData[i]["description"] = module.Description
	mi.IndexData[i]["started"] = module.Started
	mi.IndexData[i]["finished"] = module.Finished
}
