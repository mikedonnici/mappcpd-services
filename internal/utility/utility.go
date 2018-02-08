package utility

import (
	"fmt"
	"reflect"
	"time"
	"errors"
	"runtime"
	"strings"

	"github.com/mappcpd/web-services/internal/constants"

)

// dateTime converts a MySQL timestamp (format "2006-01-02 15:04:05") string to a time.Time value
func DateTime(dt string) (time.Time, error) {
	return time.Parse(constants.MySQLTimestampFormat, dt)
}

// MysqlTimestamp converts a time.Time value into a MySQL timestamp - format "2006-01-02 15:04:05"
func MySQLTimestamp(t time.Time) string {
	return t.Format(constants.MySQLTimestampFormat)
}

// q is the query
// xdf - slice of date fields
func MongofyDateFilters(q map[string]interface{}, xdf []string) {

	for _, df := range xdf {
		//fmt.Println(q[df])
		d, ok := q[df]
		if ok {
			q[df] = MongoDateFilters(d)
		}
	}
}

func MongoDateFilters(i interface{}) map[string]interface{} {

	tmp := map[string]interface{}{}

	m := reflect.ValueOf(i)
	// For each item set a bson.M object in the query...
	for _, key := range m.MapKeys() {

		// m is a map, with a single key
		//fmt.Printf("key = %v, value = %v\n", key, m.MapIndex(key))

		// get strings from reflect.Value
		k := fmt.Sprintf("%s", key)
		v := fmt.Sprintf("%s", m.MapIndex(key))
		// convert the string date into a time.Time value
		t := time.Time{}
		t.UnmarshalText([]byte(v))

		tmp[k] = t
	}

	return tmp
}

// DateStringToTime will attempt to parse a date, or datetime string and convert it to a viable time.Time value.
// This function is intended for use when the format could be one of a few different types.
func DateStringToTime(dateString string) (time.Time, error) {

	var t time.Time
	var err error

	// Try MySQL timestamp format "2006-01-02 15:04:05"
	t, err = time.Parse(constants.MySQLTimestampFormat, dateString)
	if err == nil {
		return t, err
	}

	// Try MySQL date format "2006-01-02"
	t, err = time.Parse(constants.MySQLDateFormat, dateString)
	if err == nil {
		return t, err
	}

	// Try RFC3339 format "2006-01-02T15:04:05Z07:00"
	t, err = time.Parse(time.RFC3339, dateString)
	if err == nil {
		return t, err
	}

	msg := "Error parsing date string - expected one of the following layouts: '%s', '%s', OR '%s'"
	msg = fmt.Sprintf(msg, constants.MySQLTimestampFormat, constants.MySQLDateFormat, time.RFC3339)
	return t, errors.New(msg)
}

// ErrorLocation gives the filename and function and line of the caller. Useful for errors triggered deeper in the stack
func ErrorLocationMessage(function uintptr, file string, line int, ok bool, stripPaths bool) string {

	if !ok {
		return "Could not determine error location"
	}

	var funcName string
	if stripPaths {
		file = StripPath(file)
		funcName = StripPath(runtime.FuncForPC(function).Name())
	}

	return fmt.Sprintf("File: %s  Function: %s Line: %d", file, funcName, line)
}

// StripPaths removes all path information and returns the final file name
func StripPath(filePath string) string {
	i := strings.LastIndex(filePath, "/")
	if i == -1 {
		return filePath
	} else {
		return filePath[i+1:]
	}
}

