package main

import (
	"net/http"
)

func main() {
	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// register handler for the static files
	mux.Handle("/", FileSystemHandler)

	// Register the routes and handlers
	mux.Handle("/todos", &TodosHandler{})
	mux.Handle("/todos/", &TodosHandler{})

	// Run the server
	http.ListenAndServe(":8080", mux)
}
