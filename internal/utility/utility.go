package utility

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// DateTime converts a MySQL timestamp (format "2006-01-02 15:04:05") string to a time.Time value
func DateTime(dt string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", dt)
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
