package data

import (
	"fmt"

	"github.com/mappcpd/web-services/internal/members"
)

// GetMembers fetches a list of members
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
