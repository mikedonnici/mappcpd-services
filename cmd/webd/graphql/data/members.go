package data

import (
	"fmt"
	"github.com/mappcpd/web-services/internal/members"
)

func GetMembers() []members.Member {

	q := map[string]interface{}{}
	p := map[string]interface{}{}
	l := 10

	xm, err := members.SearchMembersCollection(q, p, l)
	if err != nil {
		fmt.Println(err)
	}

	//m1 := members.Member{
	//	ID: 1,
	//	LastName: "Donnici",
	//}
	//m2 := members.Member{
	//	ID: 2,
	//	LastName: "Smith",
	//}
	//xm := []members.Member{m1,m2}

	return xm
}
