package utility

import (
	"fmt"
	"reflect"
	"time"

	"github.com/mappcpd/web-services/internal/constants"
)

// dateTime converts a MySQL timestamp (format "2006-01-02 15:04:05") string to a time.Time value
func DateTime(dt string) (time.Time, error) {
	return time.Parse(constants.MySQLTimestampFormat, dt)
}

// mysqlTimestamp converts a time.Time value into a MySQL timestamp - format "2006-01-02 15:04:05"
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
