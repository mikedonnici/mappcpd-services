package main

import (
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/cardiacsociety/web-services/internal/qualification"
)

type qualificationIndex struct {
	Name      string
	RawData   []qualification.Qualification
	IndexData []algoliasearch.Object
	Error     error
}

// newQualificationIndex returns a pointer to qualificationIndex value initialised with the index name
func newQualificationIndex(name string) qualificationIndex {
	return qualificationIndex{
		Name: name,
	}
}

func (qi *qualificationIndex) indexName() string {
	return qi.Name
}

// Number of qualifications is relatively small so this does same as fullIndex
func (qi *qualificationIndex) partialIndex() ([]algoliasearch.Object, error) {
	return qi.fullIndex()
}

func (qi *qualificationIndex) fullIndex() ([]algoliasearch.Object, error) {
	qi.fetchAllData()
	qi.createIndexObjects()
	return qi.IndexData, qi.Error
}

func (qi *qualificationIndex) fetchAllData() {
	qi.RawData, qi.Error = qualification.All(DS)
}

func (qi *qualificationIndex) createIndexObjects() {
	for i := range qi.RawData {
		qi.IndexData = append(qi.IndexData, algoliasearch.Object{})
		qi.createObject(i)
	}
}

func (qi *qualificationIndex) createObject(i int) {
	qual := qi.RawData[i]
	qi.IndexData[i] = make(map[string]interface{})
	qi.IndexData[i]["id"] = qual.ID
	qi.IndexData[i]["code"] = qual.Code
	qi.IndexData[i]["name"] = qual.Name
	qi.IndexData[i]["description"] = qual.Description
}
