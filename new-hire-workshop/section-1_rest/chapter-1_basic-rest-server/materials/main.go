package main

import (
	"fmt"
	"html"
	"net/http"
)

func main() {
	// Create a router to handle the movement of traffic in our web server.
	router := http.NewServeMux()

	// Register our handler to the index of our web server, or "/".
	router.HandleFunc("/", indexHandler)

	// Create a HTTP server using the router we just created, and a port of 8080.
	port := 8080
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	// Broadcast the HTTP server on port 8080 of localhost, with an handler on the path "/".
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}

// Create a handler function that responds to a request with a string saying "Hello, "/"".
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello on path %q", html.EscapeString(r.URL.Path))
}
