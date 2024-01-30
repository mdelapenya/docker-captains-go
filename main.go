package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Order     int    `json:"order"`
}

func (t *Todo) URL() string {
	return fmt.Sprintf("/todos/%s", t.ID)
}

func main() {
	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// register handler for the static files
	mux.Handle("/", FileSystemHandler)

	store, err := NewTodosRepository(context.Background(), os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Fatalf("Cannot create a Todos repository. Exiting")
	}

	todosHandler := NewTodosHandler(store)

	// Register the routes and handlers
	mux.Handle("/todos", todosHandler)
	mux.Handle("/todos/", todosHandler)

	// Run the server
	http.ListenAndServe(":8080", mux)
}
