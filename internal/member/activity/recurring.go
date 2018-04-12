package activity

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/mappcpd/web-services/internal/constants"
	"github.com/mappcpd/web-services/internal/platform/datastore"
)

// Recurring maps to doc in Recurring collection. Recurring activity items are
// stored in a separate collection to avoid complications with sync from MySQL -> MongoDB - ie, recurring activities
// are NOT stored in MySQL. They are store in a single document that belongs to a member
type Recurring struct {
	//ID         string              `json:"_id" bson:"_id"`
	CreatedAt  time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time           `json:"updatedAt" bson:"updatedAt"`
	MemberID   int                 `json:"memberId" bson:"memberId" validate:"required,min=1"`
	Activities []RecurringActivity `json:"activities" bson:"activities"`
}

// RecurringActivity represents an individual recurring activity.
type RecurringActivity struct {
	ID          bson.ObjectId `json:"_id" bson:"_id"`
	ActivityID  int           `json:"activityId" bson:"activityId" validate:"required,min=1"`
	CreatedAt   time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt" bson:"updatedAt"`
	Quantity    float64       `json:"quantity" validate:"required"`
	Description string        `json:"description" validate:"required"`
	Type        string        `json:"type" validate:"required"`
	Next        time.Time     `json:"next"`
}

// MemberRecurring initialises a value of type Recurring and returns a pointer to same.
// It checks for an existing doc belonging to member (id), and if not found initialises a new one.
func MemberRecurring(id int) (*Recurring, error) {

	// Initialise value with member id...
	r := Recurring{MemberID: id}

	// Get a pointer to the collection...
	c, err := datastore.MongoDB.RecurringCol()
	if err != nil {
		return &r, errors.New("MemberRecurring() could not get a pointer to collection -" + err.Error())
	}

	// Check for an existing doc using member id, and scan into value. If no doc found just the initialised
	// value is returned
	q := c.Find(bson.M{"memberId": id})
	err = q.One(&r)
	if err == mgo.ErrNotFound {
		// New doc, so set createdAt
		r.CreatedAt = time.Now()
	} else if err != nil {
		return &r, errors.New("MemberRecurring() database error -" + err.Error())
	}

	return &r, nil
}

// Save saves the Recurring value to MongoDB
func (r *Recurring) Save() error {

	// get a pointer to the collection...
	c, err := datastore.MongoDB.RecurringCol()
	if err != nil {
		fmt.Println("Recurring.Save() could not get a pointer to collection -", err)
		return err
	}

	// Upsert the record, the selector is the member id
	mid := map[string]int{"memberId": r.MemberID}
	_, err = c.Upsert(mid, r)
	if err != nil {
		fmt.Println("Recurring.Save() upsert failed -", err)
		return err
	}

	return nil
}

// RemoveActivity removes one of the recurring activities from the Recurring.Activities and saves the resulting doc
func (r *Recurring) RemoveActivity(_id string) error {

	// The activities are stored in a doc, as sub docs in an activity array.
	// Removing one of them can be achieved with some fancy Mongo using $pull - like this:
	// db.Recurring.update({"activities._id": ObjectId("59091436a9fb6e78d8945157")}, {$pull: {"activities": {"_id": ObjectId("59091436a9fb6e78d8945157")}}})

	// get a pointer to the collection...
	c, err := datastore.MongoDB.RecurringCol()
	if err != nil {
		fmt.Println("Recurring.RemoveActivity() could not get a pointer to collection -", err)
		return err
	}

	// Selector and updater
	s := bson.M{"activities._id": bson.ObjectIdHex(_id)}
	u := bson.M{"$pull": bson.M{"activities": bson.M{"_id": bson.ObjectIdHex(_id)}}}
	err = c.Update(s, u)
	if err != nil {
		fmt.Println("Recurring.RemoveActivity() update error -", err)
		return err
	}

	// The database operation has succeeded but we haven't dropped the item from our struct slice!
	// So this means the OLD way of simply changing the value in the struct first, and then saving, is JUST as efficient!
	newMap := []RecurringActivity{}
	for _, v := range r.Activities {
		// skip the one we deleted
		if v.ID == bson.ObjectIdHex(_id) {
			continue
		}
		newMap = append(newMap, v)
	}
	r.Activities = newMap

	return nil
}

// GetActivity returns just the RecurringActivity identified by _id
func (r *Recurring) GetActivity(_id string) (RecurringActivity, error) {

	for _, v := range r.Activities {
		if v.ID == bson.ObjectIdHex(_id) {
			return v, nil
		}
	}

	return RecurringActivity{}, errors.New("No activity with id " + _id)
}

// Record writes a member activity record and sets the Next scheduled time for the recurring activity
func (r *Recurring) Record(_id string) error {

	a, err := r.GetActivity(_id)
	if err != nil {
		return err
	}

	// Make idempotent by not allowing to skip if date is in the future
	if a.Next.After(time.Now()) {
		return errors.New(".Record() cannot record a recurring activity if .Next is in the future")
	}

	ar := MemberActivityInput{}
	ar.MemberID = r.MemberID
	ar.ActivityID = a.ActivityID
	ar.Date = a.Next.Format(constants.MySQLDateFormat)
	ar.Quantity = a.Quantity
	ar.Description = a.Description

	// Add activity to database
	_, err = AddMemberActivity(ar)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Increment next
	a.UpdateNext()
	r.UpdateActivity(a)

	return nil
}

// Skip just sets the Next scheduled time for the recurring activity, and saves to db
func (r *Recurring) Skip(_id string) error {

	a, err := r.GetActivity(_id)
	if err != nil {
		return err
	}

	// Make idempotent by not allowing to skip if date is in the future
	if a.Next.After(time.Now()) {
		return errors.New(".Skip() cannot skip a recurring activity if .Next is in the future")
	}

	// Increment next
	a.UpdateNext()
	r.UpdateActivity(a)

	return nil
}

// UpdateActivity updates one RecurringActivity in the Recurring.Activities slice and saves it to the database
func (r *Recurring) UpdateActivity(a RecurringActivity) {

	// Replace the activity with a matching id
	newMap := []RecurringActivity{}
	for _, v := range r.Activities {
		if v.ID == a.ID {
			v = a
		}
		newMap = append(newMap, v)
	}
	r.Activities = newMap

	r.Save()
}

// UpdateNext pushed RecurringActivity.Next schedule forward
func (a *RecurringActivity) UpdateNext() {

	switch a.Type {
	case "daily":
		a.Next = a.Next.AddDate(0, 0, 1)
	case "weekly":
		a.Next = a.Next.AddDate(0, 0, 7)
	case "monthly":
		a.Next = a.Next.AddDate(0, 1, 0)
	}
}
