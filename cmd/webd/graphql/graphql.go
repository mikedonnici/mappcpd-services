package graphql

import (
	"fmt"
	"io"
	"net/http"
)

func Start() {
	fmt.Println("Starting GraphQL server...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "GraphQL server!")
	})
	http.ListenAndServe(":5000", nil)
}
