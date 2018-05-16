package member_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/mappcpd/web-services/internal/member"
	"github.com/mappcpd/web-services/testdata"
	"github.com/matryer/is"
	"gopkg.in/mgo.v2/bson"
)

var data = testdata.NewDataStore()
var helper = testdata.NewHelper()

func TestMain(m *testing.M) {
	err := data.SetupMySQL()
	if err != nil {
		log.Fatalln(err)
	}
	defer data.TearDownMySQL()

	err = data.SetupMongoDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer data.TearDownMongoDB()

	m.Run()
}

func TestPingDatabase(t *testing.T) {
	is := is.New(t)
	err := data.Store.MySQL.Session.Ping()
	is.NoErr(err) // Could not ping test database
}

func TestByID(t *testing.T) {
	is := is.New(t)
	m, err := member.ByID(data.Store, 1)
	is.NoErr(err)                                              // Error fetching member by id
	is.True(m.Active)                                          // Active should be true
	is.Equal(m.LastName, "Donnici")                            // Last name incorrect
	is.Equal(m.Memberships[0].Title, "Associate")              // Incorrect membership title
	is.Equal(m.Contact.EmailPrimary, "michael@mesa.net.au")    // Email incorrect
	is.Equal(m.Contact.Mobile, "0402123123")                   // Mobile incorrect
	is.Equal(m.Contact.Locations[0].City, "Jervis Bay")        // Location city incorrect
	is.Equal(m.Qualifications[0].Name, "PhD")                  // Qualification incorrect
	is.Equal(m.Specialities[1].Name, "Cardiac Cath Lab Nurse") // Speciality incorrect
	//printJSON(*m)
}

func TestSearchDocDB(t *testing.T) {
	is := is.New(t)
	q := bson.M{"id": 7821}
	m, err := member.SearchDocDB(data.Store, q)
	is.NoErr(err)                     // Error querying MongoDB
	is.Equal(m[0].LastName, "Rousos") // Last name incorrect
}

func TestSaveDocDB(t *testing.T) {
	is := is.New(t)
	mem := member.Member{
		ID:          1,
		Active:      true,
		Title:       "Mr",
		FirstName:   "Michael",
		MiddleNames: "Peter",
		LastName:    "Donnici",
		Gender:      "M",
		DateOfBirth: "1970-11-03",
	}
	err := mem.SaveDocDB(data.Store)
	is.NoErr(err) // Error saving to MongoDB

	q := bson.M{"lastName": "Donnici"}
	xm, err := member.SearchDocDB(data.Store, q)
	m := xm[0]
	is.NoErr(err)     // Error querying MongoDB
	is.Equal(m.ID, 1) // ID should be 1
}

func printJSON(m member.Member) {
	xb, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println("-------------------------------------------------------------------")
	fmt.Print(string(xb))
	fmt.Println("-------------------------------------------------------------------")
}
