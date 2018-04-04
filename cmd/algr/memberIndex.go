package main

import (
	"fmt"

	"github.com/mappcpd/web-services/internal/members"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"gopkg.in/mgo.v2/bson"
)

type MemberIndex struct {
	Name      string
	RawData   []members.Member
	IndexData []algoliasearch.Object
	Error     error
}

// NewMemberIndex returns a pointer to MemberIndex value initialised with the index name
func NewMemberIndex(name string) MemberIndex {
	return MemberIndex{
		Name: name,
	}
}

func (mi *MemberIndex) FreshIndex() ([]algoliasearch.Object, error) {
	mi.fetchRawData()
	mi.createIndexObjects()
	return mi.IndexData, mi.Error
}

func (mi *MemberIndex) IndexName() string {
	return mi.Name
}

func (mi *MemberIndex) fetchRawData() {
	query := bson.M{"memberships.title": bson.M{"$exists": true}}
	mi.RawData, mi.Error = members.FetchMembers(query, 0)
}

func (mi *MemberIndex) createIndexObjects() {
	for i := range mi.RawData {
		mi.IndexData = append(mi.IndexData, algoliasearch.Object{})
		mi.createObject(i)
	}
}

func (mi *MemberIndex) createObject(i int) {

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

func (mi *MemberIndex) setLocationByType(i int, locationType string) {
	var s string
	for _, l := range mi.RawData[i].Contact.Locations {
		if l.Description == locationType {
			s = fmt.Sprintf("%s %s %s %s", l.City, l.State, l.Postcode, l.Country)
		}
	}
	mi.IndexData[i]["location"] = s
}

func (mi *MemberIndex) setSpecialities(i int) {
	var xs []string
	for _, s := range mi.RawData[i].Specialities {
		xs = append(xs, s.Name)
	}
	mi.IndexData[i]["specialities"] = xs
}

func (mi *MemberIndex) setAffiliations(i int) {
	var xs []string
	for _, s := range mi.RawData[i].Positions {
		xs = append(xs, s.OrgName)
	}
	mi.IndexData[i]["affiliations"] = xs
}
