package schema

import (
	"github.com/graphql-go/graphql"

	"github.com/mappcpd/web-services/internal/events"
)

// Event is slightly trimmer version of an events.Event
type Event struct {
	ID          int    `json:"id"`
	DateStart   string `json:"dateStart"`
	DateEnd     string `json:"dateEnd"`
	Location    string `json:"location"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// GetEvents returns a list of events based on supplied filters
func GetEvents(daysBack, daysForward int) ([]Event, error) {

	var xle []Event // local Event type

	xe, err := events.DaysRange(daysBack, daysForward)
	if err != nil {
		return nil, err
	}

	// map each events.Event to local Event type
	for _, v := range xe {
		e := Event{}
		e.ID = v.ID
		e.DateStart = v.DateStart
		e.DateEnd = v.DateEnd
		e.Location = v.Location
		e.Name = v.Name
		e.Description = v.Description
		e.URL = v.URL
		xle = append(xle, e)
	}

	return xle, nil
}


// EventsQuery field fetches Events
var EventsQuery = &graphql.Field{
	Description: "Fetches a list of events. Optional args can be passed to specify how many days back, or forward, " +
		"the event start date should be. Default is to show events with a start date in the past year.",
	Type:        graphql.NewList(event),
	Args: graphql.FieldConfigArgument{
		"daysBack": &graphql.ArgumentConfig{
			Type:        graphql.Int,
			Description: "Include events with start dates that fall from today, to this many days back",
		},
		"daysForward": &graphql.ArgumentConfig{
			Type:        graphql.Int,
			Description: "Include events with start dates that fall from today, to this many days forward",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// if no args these will be zero values
		db, ok := p.Args["daysBack"].(int)
		if !ok {
			db = 0
		}
		df, ok := p.Args["daysForward"].(int)
		if !ok {
			df = 0
		}

		// default behaviour to show events from the past 12 months
		if db == 0 && df == 0 {
			db = 365
		}

		return GetEvents(db, df)
	},
}

// event (object) defines the fields (properties) of an event
var event = graphql.NewObject(graphql.ObjectConfig{
	Name: "event",
	Description: "An event is an organised activity such as a conference, seminar, workshop etc.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of the event record",
		},
		"dateStart": &graphql.Field{
			Type:        graphql.String,
			Description: "The start date for the event",
		},
		"dateEnd": &graphql.Field{
			Type:        graphql.String,
			Description: "The end date for the event",
		},
		"location": &graphql.Field{
			Type:        graphql.String,
			Description: "The location of the event",
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the event",
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "A description of the event",
		},
		"url": &graphql.Field{
			Type:        graphql.String,
			Description: "A URL relevant to the event",
		},
	},
})
