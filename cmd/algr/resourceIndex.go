package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/mappcpd/web-services/internal/resources"
)

type resourceIndex struct {
	Name      string
	RawData   []resources.Resource
	IndexData []algoliasearch.Object
	Error     error
}

// newResourceIndex returns a pointer to resourceIndex value initialised with the index name
func newResourceIndex(name string) resourceIndex {
	return resourceIndex{
		Name: name,
	}
}

func (ri *resourceIndex) freshIndex() ([]algoliasearch.Object, error) {
	ri.fetchRawData()
	ri.createIndexObjects()
	return ri.IndexData, ri.Error
}

func (ri *resourceIndex) indexName() string {
	return ri.Name
}

func (ri *resourceIndex) fetchRawData() {
	query := bson.M{"active": true, "primary": true}
	ri.RawData, ri.Error = resources.FetchResources(query, 0)
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
