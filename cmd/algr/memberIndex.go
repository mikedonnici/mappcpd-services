package main

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/mappcpd/web-services/internal/member"
	"gopkg.in/mgo.v2/bson"
)

type memberIndex struct {
	Name      string
	RawData   []member.Member
	IndexData []algoliasearch.Object
	Error     error
}

// newMemberIndex returns a pointer to memberIndex value initialised with the index name
func newMemberIndex(name string) memberIndex {
	return memberIndex{
		Name: name,
	}
}

func (mi *memberIndex) freshIndex() ([]algoliasearch.Object, error) {
	mi.fetchRawData()
	mi.createIndexObjects()
	return mi.IndexData, mi.Error
}

func (mi *memberIndex) indexName() string {
	return mi.Name
}

func (mi *memberIndex) fetchRawData() {
	query := bson.M{"memberships.title": bson.M{"$exists": true}}
	mi.RawData, mi.Error = member.FetchMembers(query, 0)
}

func (mi *memberIndex) createIndexObjects() {
	for i := range mi.RawData {
		mi.IndexData = append(mi.IndexData, algoliasearch.Object{})
		mi.createObject(i)
	}
}

func (mi *memberIndex) createObject(i int) {

	member := mi.RawData[i]

	mi.IndexData[i] = make(map[string]interface{})
	mi.IndexData[i]["objectID"] = member.OID
	mi.IndexData[i]["_id"] = member.OID
	mi.IndexData[i]["id"] = member.ID
	mi.IndexData[i]["active"] = member.Active
	mi.IndexData[i]["name"] = fmt.Sprintf("%s %s %s", member.Title, member.FirstName, member.LastName)
	mi.IndexData[i]["email"] = member.Contact.EmailPrimary
	mi.IndexData[i]["mobile"] = member.Contact.Mobile
	mi.IndexData[i]["membership"] = member.Memberships[0].Title

	mi.setLocationByType(i, "Directory")
	mi.setSpecialities(i)
	mi.setAffiliations(i)
}

func (mi *memberIndex) setLocationByType(i int, locationType string) {
	var s string
	for _, l := range mi.RawData[i].Contact.Locations {
		if l.Description == locationType {
			s = fmt.Sprintf("%s %s %s %s", l.City, l.State, l.Postcode, l.Country)
		}
	}
	mi.IndexData[i]["location"] = s
}

func (mi *memberIndex) setSpecialities(i int) {
	var xs []string
	for _, s := range mi.RawData[i].Specialities {
		xs = append(xs, s.Name)
	}
	mi.IndexData[i]["specialities"] = xs
}

func (mi *memberIndex) setAffiliations(i int) {
	var xs []string
	for _, s := range mi.RawData[i].Positions {
		xs = append(xs, s.OrgName)
	}
	mi.IndexData[i]["affiliations"] = xs
}
