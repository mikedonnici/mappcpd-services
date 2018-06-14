package main

import (
	"time"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/cardiacsociety/web-services/internal/resource"
	"gopkg.in/mgo.v2/bson"
)

type resourceIndex struct {
	Name      string
	RawData   []resource.Resource
	IndexData []algoliasearch.Object
	Error     error
}

// newResourceIndex returns a pointer to resourceIndex value initialised with the index name
func newResourceIndex(name string) resourceIndex {
	return resourceIndex{
		Name: name,
	}
}

func (ri *resourceIndex) indexName() string {
	return ri.Name
}

func (ri *resourceIndex) partialIndex() ([]algoliasearch.Object, error) {
	ri.fetchLimitedData()
	ri.createIndexObjects()
	return ri.IndexData, ri.Error
}

func (ri *resourceIndex) fullIndex() ([]algoliasearch.Object, error) {
	ri.fetchAllData()
	ri.createIndexObjects()
	return ri.IndexData, ri.Error
}

func (ri *resourceIndex) fetchLimitedData() {
	timeBack := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	query := bson.M{"active": true, "primary": true, "updatedAt": bson.M{"$gte": timeBack}}
	ri.RawData, ri.Error = resource.FetchResources(DS, query, 0)
}

func (ri *resourceIndex) fetchAllData() {
	query := bson.M{"active": true, "primary": true}
	ri.RawData, ri.Error = resource.FetchResources(DS, query, 0)
}

func (ri *resourceIndex) createIndexObjects() {
	for i := range ri.RawData {
		ri.IndexData = append(ri.IndexData, algoliasearch.Object{})
		ri.createObject(i)
	}
}

func (ri *resourceIndex) createObject(i int) {

	obj := make(map[string]interface{})
	resource := ri.RawData[i]

	pubDate := resource.PubDate.Date.Format(time.RFC3339)
	pubStamp := resource.PubDate.Date.Unix()

	obj["objectID"] = resource.OID
	obj["_id"] = resource.OID
	obj["id"] = resource.ID
	obj["createdAt"] = resource.CreatedAt
	obj["updateAt"] = resource.UpdatedAt
	obj["publishedAt"] = pubDate
	obj["publishedAtTimestamp"] = pubStamp
	obj["type"] = resource.Type
	obj["name"] = resource.Name
	obj["description"] = resource.Description
	obj["keywords"] = resource.Keywords
	obj["shortUrl"] = resource.ShortURL
	obj["resourceUrl"] = resource.ResourceURL

	// Attributes are map[string]interface{} and may, or may not be present
	xa := []string{
		"sourceId", "sourceName", "sourceNameAbbrev",
		"sourcePubDate", "sourceVolume", "sourceIssue", "sourcePages",
	}
	for _, a := range xa {
		v, ok := resource.Attributes[a]
		if ok {
			obj[a] = v
		}
	}

	ri.IndexData[i] = obj
}
