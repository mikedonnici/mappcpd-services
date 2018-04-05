package main

import (
	"fmt"
	"strings"
	"os"

	"gopkg.in/mgo.v2/bson"
	"github.com/mappcpd/web-services/internal/members"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/34South/envr"
)

func init() {

	envr.New("algrEnv", []string{
		"MAPPCPD_ALGOLIA_DIRECTORY_EXCLUDE_TITLES",
	}).Auto()
}

type directoryIndex struct {
	Name      string
	RawData   []members.Member
	IndexData []algoliasearch.Object
	Error     error
}

func newDirectoryIndex(name string) directoryIndex {
	return directoryIndex{
		Name: name,
	}
}

func (di *directoryIndex) freshIndex() ([]algoliasearch.Object, error) {
	di.fetchRawData()
	di.removeExcludedMembers()
	di.createIndexObjects()
	return di.IndexData, di.Error
}

func (di *directoryIndex) indexName() string {
	return di.Name
}

func (di *directoryIndex) fetchRawData() {
	query := bson.M{"memberships.title": bson.M{"$exists": true}}
	di.RawData, di.Error = members.FetchMembers(query, 0)
}

func (di *directoryIndex) removeExcludedMembers() {
	var xm []members.Member
	for _, m := range di.RawData {
		if shouldInclude(m) {
			xm = append(xm, m)
		}
	}
	di.RawData = xm
}

func (di *directoryIndex) createIndexObjects() {
	for i := range di.RawData {
		di.IndexData = append(di.IndexData, algoliasearch.Object{})
		di.createObject(i)
	}
}

func (di *directoryIndex) createObject(i int) {

	member := di.RawData[i]

	di.IndexData[i] = make(map[string]interface{})
	di.IndexData[i]["objectID"] = member.OID
	di.IndexData[i]["_id"] = member.OID
	di.IndexData[i]["id"] = member.ID
	di.IndexData[i]["active"] = member.Active
	di.IndexData[i]["name"] = fmt.Sprintf("%s %s %s", member.Title, member.FirstName, member.LastName)
	di.IndexData[i]["email"] = member.Contact.EmailPrimary
	di.IndexData[i]["mobile"] = member.Contact.Mobile
	di.IndexData[i]["membership"] = member.Memberships[0].Title

	di.setLocationByType(i, "Directory")
	di.setSpecialities(i)
	di.setAffiliations(i)
}

func (di *directoryIndex) setLocationByType(i int, locationType string) {
	var s string
	for _, l := range di.RawData[i].Contact.Locations {
		if l.Description == locationType {
			s = fmt.Sprintf("%s %s %s %s", l.City, l.State, l.Postcode, l.Country)
		}
	}
	di.IndexData[i]["location"] = s
}

func (di *directoryIndex) setSpecialities(i int) {
	var xs []string
	for _, s := range di.RawData[i].Specialities {
		xs = append(xs, s.Name)
	}
	di.IndexData[i]["specialities"] = xs
}

func (di *directoryIndex) setAffiliations(i int) {
	var xs []string
	for _, s := range di.RawData[i].Positions {
		xs = append(xs, s.OrgName)
	}
	di.IndexData[i]["affiliations"] = xs
}

func shouldInclude(m members.Member) bool {
	return isActive(m)
}

func isActive(m members.Member) bool {
	if m.Active != true {
		return false
	}
	return hasDirectoryConsent(m)
}

func hasDirectoryConsent(m members.Member) bool {
	if m.Contact.Directory != true {
		return false
	}
	return hasMembershipTitle(m)
}

func hasMembershipTitle(m members.Member) bool {
	if m.Title == "" {
		return false
	}
	return hasMembershipTitleNotExcluded(m)
}

func hasMembershipTitleNotExcluded(m members.Member) bool {
	xs := strings.Split(os.Getenv("MAPPCPD_ALGOLIA_DIRECTORY_EXCLUDE_TITLES"), ",")
	title := strings.ToLower(m.Title)
	for _, s := range xs {
		excludeTitle := strings.ToLower(strings.TrimSpace(s))
		if title == excludeTitle {
			return false
		}
	}

	return true
}
