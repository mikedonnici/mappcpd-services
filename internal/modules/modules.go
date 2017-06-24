package modules

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/mappcpd/api/db"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/mappcpd/web-services/internal/utility"
)

// Module defines struct for a CPD module
type Module struct {
	_id             string    `json:"_id" bson:"_id"`
	ID              int       `json:"id" bson:"id"`
	CreatedAt       time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt" bson:"updatedAt"`
	PublishedAt     time.Time `json:"publishedAt" bson:"publishedAt"`
	Name            string    `json:"name" bson:"name"`
	Description     string    `json:"description" bson:"description"`
	DurationMinutes int       `json:"durationMinutes" bson:"durationMinutes"`
	Started         int       `json:"started" bson:"started"`
	Finished        int       `json:"finished" bson:"finished"`
	Current         bool      `json:"current" bson:"current"`
}

type Modules []Module

// ModuleByID fetches a module by id, from the MySQL db
func ModuleByID(id int) (*Module, error) {

	// Set up a new empty Member
	m := Module{ID: id}

	// Coalesce any NULL-able fields
	query := `
	SELECT
	olm.created_at,
	olm.updated_at,
	olm.published_at,
	COALESCE(olm.name, ''),
	COALESCE(olm.description, ''),
	olm.started, olm.finished, olm.current
	FROM ol_module olm
	WHERE active = 1 AND
	olm.id = ?`

	// Hold these until we fix them up
	var createdAt string
	var updatedAt string
	var publishedAt string

	err := datastore.MySQL.Session.QueryRow(query, id).Scan(
		&createdAt,
		&updatedAt,
		&publishedAt,
		&m.Name,
		&m.Description,
		&m.Started,
		&m.Finished,
		&m.Current,
	)
	if err != nil {
		return &m, err
	}
	// Convert MySQL date time strings to time.Time
	m.CreatedAt, _ = utility.DateTime(createdAt)
	m.UpdatedAt, _ = utility.DateTime(updatedAt)
	m.PublishedAt, _ = utility.DateTime(publishedAt)

	return &m, nil
}

// DocModulesAll searches the Modules collection.
func DocModulesAll(q map[string]interface{}, p map[string]interface{}) ([]interface{}, error) {

	modules, err := datastore.MongoDB.ModulesCol()
	if err != nil {
		return nil, err
	}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt", "publishedAt"})

	// Run query and return results TODO remove limit here!!
	var m []interface{}
	err = modules.Find(q).Select(p).All(&m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// DocModulesLimit returns n modules
func DocModulesLimit(q map[string]interface{}, p map[string]interface{}, l int) ([]interface{}, error) {

	m := []interface{}{}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt", "publishedAt"})

	modules, err := db.MongoDB.ModulesCol()
	if err != nil {
		return m, err
	}
	err = modules.Find(q).Select(p).Limit(l).All(&m)
	if err != nil {
		return m, err
	}

	return m, nil
}

// DocModuleOne returns one module, unmarshaled into the proper struct
// so no projection allowed here
func DocModulesOne(q map[string]interface{}) (Module, error) {

	m := Module{}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(q, []string{"updatedAt", "createdAt", "publishedAt"})

	modules, err := db.MongoDB.ModulesCol()
	if err != nil {
		return m, err
	}
	err = modules.Find(q).One(&m)
	if err != nil {
		return m, err
	}

	return m, nil
}

// QueryModulesCollection ... queries the modules collection :)
func QueryModulesCollection(mq db.MongoQuery) ([]interface{}, error) {

	// results
	r := []interface{}{}

	// Convert string date filters to time.Time
	utility.MongofyDateFilters(mq.Find, []string{"updatedAt", "createdAt"})

	// get a pointer to the modules collection
	c, err := db.MongoDB.ModulesCol()
	if err != nil {
		return r, err
	}

	// execute query, scan results into r
	err = mq.Do(c, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// SyncModule synchronises the Module record from MySQL -> MongoDB
func SyncModule(m *Module) {
	// Fetch the current Doc (if there) and compare updatedAt
	m2, err := DocModulesOne(bson.M{"id": m.ID})
	fmt.Println(m2)
	if err != nil {
		log.Println("Target document error: ", err, "- so do an upsert")
	}

	msg := fmt.Sprintf("MySQL record updated at %s, MongoDB record updated at %s", m.UpdatedAt, m2.UpdatedAt)
	if m.UpdatedAt.Equal(m2.UpdatedAt) {
		msg += " - NO need to sync"
		log.Println(msg)
		return
	}
	msg += " - syncing..."
	log.Println(msg)

	// Update the document in the Members collection
	var w sync.WaitGroup
	w.Add(1)
	go UpdateModuleDoc(m, &w)
	w.Wait()
}

// UpdateModuleDoc updates a document in the Modules collection
func UpdateModuleDoc(m *Module, w *sync.WaitGroup) {

	// Make the selector for Upsert
	id := map[string]int{"id": m.ID}

	// Get pointer to the Modules collection
	mc, err := db.MongoDB.ModulesCol()
	if err != nil {
		log.Printf("Error getting pointer to Modules collection: %s\n", err.Error())
		return
	}

	// Upsert
	_, err = mc.Upsert(id, &m)
	if err != nil {
		log.Printf("Error updating document in Modules collection: %s\n", err.Error())
	}

	// Tell wait group we're done, if it was passed in
	if w != nil {
		w.Done()
	}
	log.Println("Updated Modules document")
}
