/*
	Package date provides various date functions used in the application.
*/
package date

import (
	"fmt"
	"time"

	"github.com/mappcpd/web-services/internal/constants"
	"github.com/pkg/errors"
)

// StringToTime returns a time.Time value from a number of date string formats.
func StringToTime(dateString string) (time.Time, error) {

	var t time.Time
	var err error

	t, err = time.Parse(constants.MySQLTimestampFormat, dateString)
	if err == nil {
		return t, err
	}

	t, err = time.Parse(constants.MySQLDateFormat, dateString)
	if err == nil {
		return t, err
	}

	t, err = time.Parse(time.RFC3339, dateString)
	if err == nil {
		return t, err
	}

	msg := "Error parsing date string - expected one of the following layouts: '%s', '%s', OR '%s'"
	msg = fmt.Sprintf(msg, constants.MySQLTimestampFormat, constants.MySQLDateFormat, time.RFC3339)
	return t, errors.New(msg)
}
