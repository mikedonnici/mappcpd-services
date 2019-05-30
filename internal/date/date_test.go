package date_test

import (
	"testing"

	"github.com/cardiacsociety/web-services/internal/date"
	"github.com/matryer/is"
)

const (
	d1 = "2006-01-02 15:04:05"
	d2 = "2006-01-02"
	d3 = "2006-01-02T15:04:05+10:00"
	d4 = "January 1st 2006"
)

func TestStringToTime(t *testing.T) {
	is := is.New(t)

	_, err := date.StringToTime(d1)
	is.NoErr(err) // Error converting "2006-01-02 15:04:05" to Time value

	_, err = date.StringToTime(d2)
	is.NoErr(err) // Error converting "2006-01-02" to Time value

	_, err = date.StringToTime(d3)
	is.NoErr(err) // Error converting "2006-01-02T15:04:05Z07:00" to Time value

	_, err = date.StringToTime(d4)
	is.True(err != nil) // Expect an error for date string
}
