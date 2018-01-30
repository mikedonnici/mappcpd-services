package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/mappcpd/web-services/cmd/webd/graphql/data"
	"fmt"
)

// Events query field fetches Events
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

		fmt.Println(db, df)

		return data.GetEvents(db, df)
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
