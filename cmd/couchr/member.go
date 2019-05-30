package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cardiacsociety/web-services/internal/activity"
	"github.com/cardiacsociety/web-services/internal/cpd"
	"github.com/cardiacsociety/web-services/internal/generic"
	"github.com/cardiacsociety/web-services/internal/member"
)

const memberIdPrefix = "member"

// MemberDoc stores a couchbase member doc
type MemberDoc struct {
	Type           string                    `json:"type"`
	Created        time.Time                 `json:"created"`
	Updated        time.Time                 `json:"updated"`
	Status         string                    `json:"status,omitempty"`
	Title          string                    `json:"title,omitempty"`
	Gender         string                    `json:"gender,omitempty"`
	PreNom         string                    `json:"preNom,omitempty"`
	FirstName      string                    `json:"firstName,omitempty"`
	MiddleNames    []string                  `json:"middleNames,omitempty"`
	LastName       string                    `json:"lastName,omitempty"`
	PostNom        string                    `json:"postNom,omitempty"`
	Email          string                    `json:"email,omitempty"`
	Email2         string                    `json:"email2,omitempty"`
	Mobile         string                    `json:"mobile,omitempty"`
	Directory      bool                      `json:"directoryConsent"`
	Consent        bool                      `json:"contactConsent"`
	Locations      []member.Location         `json:"locations,omitempty"`
	Qualifications []member.Qualification    `json:"qualifications,omitempty"`
	Accreditations []member.Accreditation    `json:"accreditations,omitempty"`
	Positions      []member.Position         `json:"positions,omitempty"`
	Specialities   []member.Speciality       `json:"specialities"`
	TitleHistory   []member.MembershipTitle  `json:"titleHistory,omitempty"`
	StatusHistory  []member.MembershipStatus `json:"statusHistory,omitempty"`
	CPD            []CPD                     `json:"cpd,omitempty"`
}

type CPD struct {
	Domain      string  `json:"domain"`
	Category    string  `json:"category"`
	Activity    string  `json:"activity"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
	UnitCredit  float64 `json:"creditPerUnit"`
	Credit      float64 `json:"credit"`
}

func syncMembers() {

	xi, err := generic.GetIDs(ds, "member", "")
	if err != nil {
		log.Fatalln("mysql err", err)
	}

	for _, id := range xi {

		md := &MemberDoc{}

		fmt.Print("Syncing member id ", id)
		m, err := member.ByID(ds, id)
		if err != nil {
			log.Fatalln("Could not get member id ", id, "-", err)
		}
		md.mapMemberProfile(*m)

		fmt.Print("... fetching cpd\n")

		xa, err := cpd.ByMemberID(ds, id)
		if err != nil {
			log.Fatalln("Could not get CPD for member id", id, "-", err)
		}
		if len(xa) > 0 {
			md.mapCPD(xa)
		}

		id := fmt.Sprintf("%v::%v", memberIdPrefix, m.ID)
		_, err = cb.Upsert(id, md, 0)
		if err != nil {
			log.Println("Upsert error", err)
		}
	}
}

// mapMemberProfile maps profile data from member.Member to couchbase memberDoc
func (md *MemberDoc) mapMemberProfile(m member.Member) {

	var title string
	var titleHistory []member.MembershipTitle

	var status string
	var statusHistory []member.MembershipStatus

	if len(m.Memberships) > 0 {

		title = m.Memberships[0].Title
		xt := m.Memberships[0].TitleHistory
		for _, t := range xt {
			titleHistory = append(titleHistory, t)
		}

		status = m.Memberships[0].Status
		xs := m.Memberships[0].StatusHistory
		for _, s := range xs {
			statusHistory = append(statusHistory, s)
		}
	}

	var locations []member.Location
	if len(m.Contact.Locations) > 0 {
		for _, l := range m.Contact.Locations {
			locations = append(locations, l)
		}
	}

	md.Type = "member"
	md.Created = m.CreatedAt
	md.Updated = m.UpdatedAt
	md.Gender = m.Gender
	md.PreNom = m.Title
	md.FirstName = m.FirstName
	md.MiddleNames = m.MiddleNames
	md.LastName = m.LastName
	md.PostNom = m.PostNominal
	md.Email = m.Contact.EmailPrimary
	md.Email2 = m.Contact.EmailSecondary
	md.Mobile = m.Contact.Mobile
	md.Directory = m.Contact.Directory
	md.Consent = m.Contact.Consent
	md.Locations = locations
	md.Title = title
	md.TitleHistory = titleHistory
	md.Status = status
	md.StatusHistory = statusHistory
	md.Qualifications = m.Qualifications
	md.Accreditations = m.Accreditations
	md.Specialities = m.Specialities
	md.Positions = m.Positions
}

// mapCPD maps cpd.CPD values to local, simpler version
func (md *MemberDoc) mapCPD(cpd []cpd.CPD) {

	for _, c := range cpd {
		var err error
		ca := CPD{}
		ca.Domain = "CME"

		ca.Category, err = mapActivityToCategory(c.Activity)
		if err != nil {
			log.Println("Error mapping old activity to new category -", err)
		}

		ca.Activity, err = mapTypeToActivity(c)
		if err != nil {
			log.Println("Error mapping activity type to new activity -", err)
		}

		ca.Date = c.Date
		ca.Description = c.Description
		ca.Quantity = c.CreditData.Quantity
		ca.UnitCredit = c.CreditData.UnitCredit
		ca.Unit = c.CreditData.UnitName
		ca.Credit = c.Credit

		md.CPD = append(md.CPD, ca)
	}
}

// maps an old CPD activity (from ce_activity) to a new category
func mapActivityToCategory(a activity.Activity) (string, error) {

	// ID >= 20 is a new activity which has a type - these activities are the new categories
	if a.ID >= 20 {
		return a.Name, nil
	}

	// Old activity needs to be mapped to a new category
	aid := oldToNewActivityID(a.ID)
	a2, err := activity.ByID(ds, aid)
	if err != nil {
		fmt.Println("Can't get activity name for activity with id", aid)
	}
	return a2.Name, nil
}

func mapTypeToActivity(c cpd.CPD) (string, error) {
	// ...will have a type
	if c.Activity.ID >= 20 {
		return c.Type.Name, nil
	}
	// Won't have a type so return original activity name
	return c.Activity.Name, nil
}

// Maps one of the older activity ids to the newer equivalent
func oldToNewActivityID(oldActivityID int) int {

	m := map[int]int{
		1:  22,
		2:  24,
		3:  23,
		4:  22,
		5:  22,
		6:  24,
		7:  23,
		8:  24,
		9:  20,
		10: 20,
		11: 20,
		12: 20,
		13: 20,
		14: 23,
		15: 24,
		16: 22,
		17: 24,
		18: 24,
		19: 20,
	}

	return m[oldActivityID]
}
