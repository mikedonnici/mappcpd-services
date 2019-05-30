package graphql

import (
	"github.com/cardiacsociety/web-services/internal/member"
)

// memberData is a local representation of member.Member
type memberData struct {
	ID             int                    `json:"id"`
	Token          string                 `json:"token"`
	Active         bool                   `json:"active"`
	Title          string                 `json:"title"`
	FirstName      string                 `json:"firstName"`
	MiddleNames    []string               `json:"middleNames"`
	LastName       string                 `json:"lastName"`
	PostNominal    string                 `json:"postNominal"`
	DateOfBirth    string                 `json:"dateOfBirth"`
	Email          string                 `json:"email"`
	Mobile         string                 `json:"mobile"`
	Locations      []member.Location      `json:"locations"`
	Qualifications []member.Qualification `json:"qualifications"`
	Positions      []member.Position      `json:"positions"`
}

// mapMemberData fetches a member record by id, and maps field values to the local memberData type
func mapMemberData(id int) (memberData, error) {
	var m memberData
	mp, err := member.ByID(DS, id)
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

// memberActivityAttachmentRequest requests a signed URL for uploading to S3
//func memberActivityAttachmentRequest(memberID int) string {
//
//	var url string
//
//	// Get the file set data
//	fs, err := fileset.ActivityAttachment()
//	if err != nil {
//		msg := "Could not determine the storage information for activity attachments - " + err.Error()
//		return msg
//	}
//	fmt.Println(fs)
//
//	// Use the file set information to create an upload value
//
//	return url
//}
