package resource_test

import (
	"log"
	"testing"
	"time"

	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/internal/resource"
	"github.com/cardiacsociety/web-services/testdata"
	"github.com/matryer/is"
	"gopkg.in/mgo.v2/bson"
)

var ds datastore.Datastore

func TestResource(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("resource", func(t *testing.T) {
		t.Run("testPingDatabases", testPingDatabases)
		t.Run("testByID", testByID)
		t.Run("testDocResourcesAll", testDocResourcesAll)
		t.Run("testDocResourcesLimit", testDocResourcesLimit)
		t.Run("testDocResourcesOne", testDocResourcesOne)
		t.Run("testQueryResourcesCollection", testQueryResourcesCollection)
		t.Run("testFetchResources", testFetchResources)
		t.Run("testSyncResource", testSyncResource)
		t.Run("testSaveNewResource", testSaveNewResource)
		t.Run("testSaveExistingResource", testSaveExistingResource)
	})
}

func setup() (datastore.Datastore, func()) {
	var db = testdata.NewDataStore()
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalf("SetupMySQL() err = %s", err)
	}
	err = db.SetupMongoDB()
	if err != nil {
		log.Fatalf("SetupMongoDB() err = %s", err)
	}
	return db.Store, func() {
		err := db.TearDownMySQL()
		if err != nil {
			log.Fatalf("TearDownMySQL() err = %s", err)
		}
		err = db.TearDownMongoDB()
		if err != nil {
			log.Fatalf("TearDownMongoDB() err = %s", err)
		}
	}
}

func testPingDatabases(t *testing.T) {
	err := ds.MySQL.Session.Ping()
	if err != nil {
		t.Fatalf("MySQL.Session.Ping() err = %s", err)
	}
	err = ds.MongoDB.Session.Ping()
	if err != nil {
		t.Fatalf("MongoDB.Session.Ping() err = %s", err)
	}
}

func testByID(t *testing.T) {
	cases := []struct {
		id  int
		doi string
	}{
		{id: 6576, doi: "https://doi.org/10.1053/j.gastro.2017.08.022"},
		{id: 6578, doi: "https://doi.org/10.1016/j.jaci.2017.03.020"},
	}

	for _, c := range cases {
		r, err := resource.ByID(ds, c.id)
		if err != nil {
			t.Errorf("resource.ByID(%d) err = %s", c.id, err)
		}
		want := c.doi
		got := r.ResourceURL
		if got != want {
			t.Errorf("resource.ByID(%d).ResourceURL = %q, want %q", c.id, got, want)
		}
	}
}

func testDocResourcesAll(t *testing.T) {
	cases := []struct {
		query       bson.M
		projection  bson.M
		resultCount int
	}{
		{query: bson.M{}, projection: bson.M{}, resultCount: 5},
		{query: bson.M{"id": 24967}, projection: bson.M{}, resultCount: 1},
		{query: bson.M{"id": 25000}, projection: bson.M{}, resultCount: 0},
		{query: bson.M{"keywords": "PD2018"}, projection: bson.M{}, resultCount: 1},
	}
	for _, c := range cases {
		xr, err := resource.DocResourcesAll(ds, c.query, c.projection)
		if err != nil {
			t.Errorf("resource.DocResourcesAll() err = %s", err)
		}
		got := len(xr)
		want := c.resultCount
		if got != want {
			t.Errorf("resource.DocResourcesAll() count = %d, want %d", got, want)
		}
	}
}

func testDocResourcesLimit(t *testing.T) {
	cases := []struct {
		query       bson.M
		projection  bson.M
		limit       int
		expectCount int
	}{
		{query: bson.M{}, projection: bson.M{}, limit: 2, expectCount: 2},
	}
	for _, c := range cases {
		xr, err := resource.DocResourcesLimit(ds, c.query, c.projection, c.limit)
		if err != nil {
			t.Errorf("resource.DocResourcesLimit() err = %s", err)
		}
		got := len(xr)
		want := c.expectCount
		if got != want {
			t.Errorf("resource.DocResourcesLimit() count =%d, want %d", got, want)
		}
	}
}

func testDocResourcesOne(t *testing.T) {
	cases := []struct {
		id  int
		doi string
	}{
		{id: 2000, doi: "https://webcast.gigtv.com.au/Mediasite/Play/bb4663e0c3b64cc58f200064bb6c03db1d"},
		{id: 10012, doi: "https://doi.org/10.1016/j.resuscitation.2017.08.218"},
	}
	for _, c := range cases {
		r, err := resource.DocResourcesOne(ds, bson.M{"id": c.id})
		if err != nil {
			t.Errorf("resource.DocResourcesOne() err = %s", err)
		}
		got := r.ResourceURL
		want := c.doi
		if got != want {
			t.Errorf("resource.DocResourcesOne().ResourceURL = %q, want %q", got, want)
		}
	}
}

func testQueryResourcesCollection(t *testing.T) {
	cases := []struct {
		query datastore.MongoQuery
		doi   string
	}{
		{
			query: datastore.MongoQuery{Find: bson.M{"id": 2000}},
			doi:   "https://webcast.gigtv.com.au/Mediasite/Play/bb4663e0c3b64cc58f200064bb6c03db1d",
		},
		{
			query: datastore.MongoQuery{Find: bson.M{"id": 10012}},
			doi:   "https://doi.org/10.1016/j.resuscitation.2017.08.218",
		},
	}
	for _, c := range cases {
		xr, err := resource.QueryResourcesCollection(ds, c.query)
		if err != nil {
			t.Errorf("resource.QueryResourcesCollection() err = %s", err)
		}
		gotCount := len(xr)
		wantCount := 1
		if gotCount != wantCount {
			t.Errorf("resource.QueryResourcesCollection() count = %d, want %d", gotCount, wantCount)
		}

		r := xr[0].(bson.M) // returns []interface{} so need to assert
		doi, ok := r["resourceUrl"]
		if !ok {
			t.Errorf("Resource.ResourceURL field missing")
		}
		got := doi
		want := c.doi
		if got != want {
			t.Errorf("Resource.ResourceURL = %q, want %q", got, want)
		}
	}
}
func testFetchResources(t *testing.T) {
	cases := []struct {
		query       map[string]interface{}
		expectCount int
	}{
		{
			query:       map[string]interface{}{"id": 2000},
			expectCount: 1,
		},
		{
			query:       map[string]interface{}{},
			expectCount: 5,
		},
	}
	for _, c := range cases {
		xr, err := resource.FetchResources(ds, c.query, 0) // limit = 0
		if err != nil {
			t.Errorf("resource.FetchResources() err = %s", err)
		}
		got := len(xr)
		want := c.expectCount
		if got != want {
			t.Errorf("resource.FetchResources() count = %d, want %d", got, want)
		}
	}
}

func testSyncResource(t *testing.T) {
	// this id is present in mysql but NOT in mongo test set
	arg := 6576 // resource id
	r, err := resource.ByID(ds, arg)
	if err != nil {
		t.Fatalf("resource.ByID(%d) err = %s", arg, err)
	}

	// sync and wait a bit :)
	resource.SyncResource(ds, r) // go routine, no error check
	time.Sleep(2 * time.Second)  // bad... but needs to time to sync to mono

	// fetch same id from mongo
	rd, err := resource.DocResourcesOne(ds, bson.M{"id": arg})
	if err != nil {
		t.Fatalf("resource.DocResourcesOne(%d) err = %s", arg, err)
	}
	got := rd.ID
	want := r.ID
	if got != want {
		t.Errorf("resource.DocResourcesOne(%d).ID = %d, want %d", arg, got, want)
	}
}

func testSaveNewResource(t *testing.T) {
	is := is.New(t)

	r := resource.Resource{
		Name:        "test resource",
		ResourceURL: "http://csanz.io/abcd1234",
	}
	newId, err := r.Save(ds)
	is.NoErr(err) // error saving resource

	r2, err := resource.ByID(ds, newId)
	is.NoErr(err)                           // error fetching new resource
	is.Equal(r.Name, r2.Name)               // New resource name does not match
	is.Equal(r.ResourceURL, r2.ResourceURL) // New resource url does not match
}

func testSaveExistingResource(t *testing.T) {

	// this one already exists in test mysql db (id 6578) so should error as nothing to change
	arg := 6578
	r := resource.Resource{
		ResourceURL: "https://doi.org/10.1016/j.jaci.2017.03.020",
	}
	id, err := r.Save(ds)
	if err != nil {
		t.Errorf("Resource.Save() err = %s", err)
	}
	got := id
	want := arg
	if got != want {
		t.Errorf("Resource.Save() = %d, want %d", got, want)
	}

	// Save() should have set r.ID to id
	got = r.ID
	if got != want {
		t.Errorf("Resource.ID = %d, want %d", got, want)
	}

	// Same again, only this time change the title
	r = resource.Resource{
		Name:        "New name",
		ResourceURL: "https://doi.org/10.1016/j.jaci.2017.03.020",
	}
	id, err = r.Save(ds)
	if err != nil {
		t.Errorf("Resource.Save() err = %s", err)
	}
	// expect existing resource id to be returned
	got = id
	if got != want {
		t.Errorf("Resource.Save() = %d, want %d", got, want)
	}

	// fetch updated resource
	r2, err := resource.ByID(ds, id) 
	if err != nil {
		t.Fatalf("resource.ByID(%d) err = %s", id, err)
	}
	gotName := r2.Name
	wantName := r.Name
	if gotName != wantName {
		t.Errorf("Resource.Name = %q, want %q", gotName, wantName)
	}
}
