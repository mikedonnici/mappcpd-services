package graphql

import (
	"fmt"
	"io"
	"os"
	"net/http"
)

func Start() {
	fmt.Println("Starting GraphQL server...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "GraphQL server!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.ListenAndServe(":"+port, nil)
}
