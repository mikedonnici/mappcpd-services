package data

import (
	"fmt"
	"github.com/mappcpd/web-services/internal/members"
)

// MemberViewer holds some very basic info about the member, such as ID,
// that can be passed down to child queries.
type MemberViewer struct {
	UIID string `json:"_id"`
	ID   int    `json:"id"`
	Active bool `json:"active"`
}

// GetMemberViewer fetches...
func GetMemberViewer(id int) (MemberViewer, error) {
	var m MemberViewer
	mp, err := GetMemberProfile(id)
	if err != nil {
		return m, err
	}
	m.ID = mp.ID
	m.UIID = "not available"
	m.Active = mp.Active
	return m, nil
}

// GetMember fetches a single member record by id
func GetMemberProfile(id int) (members.Member, error) {
	// MemberByID returns a pointer to a members.Member so dereference in return
	m, err := members.MemberByID(id)
	return *m, err
}

func GetMembers() []members.Member {

	q := map[string]interface{}{}
	p := map[string]interface{}{}
	l := 10

	xm, err := members.SearchMembersCollection(q, p, l)
	if err != nil {
		fmt.Println(err)
	}

	return xm
}
