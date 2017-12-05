package data

import (
	"github.com/mappcpd/web-services/internal/members"
)

// Member struct - a simpler representation than members.Member
// This is used as a 'viewer' field and the id is passed down to child queries.
type Member struct {
	ID             int                    `json:"id"`
	Active         bool                     `json:"active"`
	Title          string                   `json:"title"`
	FirstName      string                   `json:"firstName"`
	MiddleNames    string                   `json:"middleNames"`
	LastName       string                   `json:"lastName"`
	PostNominal    string                   `json:"postNominal"`
	DateOfBirth    string                   `json:"dateOfbirth"`
	Email          string                   `json:"email"`
	Mobile         string                   `json:"mobile"`
	Locations      []members.MemberLocation `json:"locations"`
	Qualifications []members.Qualification  `json:"qualifications"`
	Positions      []members.Position       `json:"positions"`
}

// GetMember fetches the basic member record
func GetMember(id int) (Member, error) {
	var m Member
	mp, err := GetMemberProfile(id)
	if err != nil {
		return m, err
	}
	m.ID = mp.ID
	m.Active = mp.Active
	m.Title = mp.Title
	m.FirstName = mp.FirstName
	m.MiddleNames = mp.MiddleNames
	m.LastName = mp.LastName
	m.DateOfBirth = mp.DateOfBirth
	m.Email = mp.Contact.EmailPrimary
	m.Mobile = mp.Contact.Mobile
	m.PostNominal = mp.PostNominal
	m.Locations = mp.Contact.Locations
	m.Qualifications = mp.Qualifications
	m.Positions = mp.Positions

	return m, nil
}

// GetMemberProfile fetches a single member record by id
func GetMemberProfile(memberID int) (members.Member, error) {
	// MemberByID returns a pointer to a members.Member so dereference in return
	m, err := members.MemberByID(memberID)
	return *m, err
}
