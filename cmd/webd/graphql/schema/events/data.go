package events

import "github.com/mappcpd/web-services/internal/events"

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

