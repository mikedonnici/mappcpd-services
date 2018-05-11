package member_test

import (
	"log"
	"testing"

	//"github.com/mappcpd/web-services/internal/member"
	"github.com/mappcpd/web-services/testdata"
)

var db = testdata.NewTestDB()
var helper = testdata.NewHelper()

func TestMain(m *testing.M) {
	err := db.Setup()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.TearDown()

	m.Run()
}

func TestPingDatabase(t *testing.T) {
	err := db.DS.MySQL.Session.Ping()
	if err != nil {
		t.Fatal("Could not ping database")
	}
}
