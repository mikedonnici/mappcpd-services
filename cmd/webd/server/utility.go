package server

import (
	"errors"
	"strings"
)

// queryParams splits search params in the form: ?q=name1:value1,name2:value2
// and returns them as a map, 's' is the string passed in.
func queryParams(s string) (map[string]interface{}, error) {

	q := map[string]interface{}{}

	// allow ?q=all to return an empty query
	if s == "all" {
		return q, nil
	}

	// Slice up query params
	xs := strings.Split(s, ",")
	// Loop through the query params and add each to the map
	for _, v := range xs {
		// Split field:value...
		fv := strings.Split(v, ":")

		//if we don't get 2 decent strings then the param is malformed
		if len(fv[0]) < 1 || len(fv[1]) < 1 {
			return q, errors.New("query parameters incorrect - should be ?q=name1:value1,name2:value2 etc")
		}

		// Also let's limit the length of each to 255 chars
		if len(fv[0]) > 255 || len(fv[1]) > 255 {
			return q, errors.New("query parameters are too long")
		}

		// Add the name-value pair to the query
		q[fv[0]] = fv[1]
	}

	return q, nil
}

// projectParams takes the ?p=field1,field2 string and creates a map
func projectParams(s string) map[string]interface{} {

	p := map[string]interface{}{}

	// Projection
	// The fields we want to project can be set up the same way as the Query,
	// ie. we create a map out of the name value pairs which will eventually
	// be passed to MongoDB as a JSON projection doc. We're using the param name
	// 'fields' as it is friendlier than 'projection', and only need to specify the
	// // fields required - no values:
	// ?....&fields=name1,name2,name3...
	//fParams := r.FormValue("f")
	xs := strings.Split(s, ",")

	for _, v := range xs {
		if len(v) > 0 {
			p[v] = "1"
		}
	}

	return p
}
