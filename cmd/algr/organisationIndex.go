package main

import (
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/cardiacsociety/web-services/internal/organisation"
)

type organisationIndex struct {
	Name      string
	RawData   []organisation.Organisation
	IndexData []algoliasearch.Object
	Error     error
}

// newOrganisationIndex returns a pointer to an organisationIndex value initialised with the index name
func newOrganisationIndex(name string) organisationIndex {
	return organisationIndex{
		Name: name,
	}
}

func (oi *organisationIndex) indexName() string {
	return oi.Name
}

// Number of organisations is relatively small so this does same as fullIndex
func (oi *organisationIndex) partialIndex() ([]algoliasearch.Object, error) {
	return oi.fullIndex()
}

func (oi *organisationIndex) fullIndex() ([]algoliasearch.Object, error) {
	oi.fetchAllData()
	oi.createIndexObjects()
	return oi.IndexData, oi.Error
}

func (oi *organisationIndex) fetchAllData() {
	oi.RawData, oi.Error = organisation.All(DS)
}

func (oi *organisationIndex) createIndexObjects() {

	for i := range oi.RawData {
		oi.IndexData = append(oi.IndexData, algoliasearch.Object{})
		oi.createObject(i)
	}
}

func (oi *organisationIndex) createObject(i int) {
	org := oi.RawData[i]
	oi.IndexData[i] = make(map[string]interface{})
	oi.IndexData[i]["id"] = org.ID
	oi.IndexData[i]["code"] = org.Code
	oi.IndexData[i]["name"] = org.Name
}
