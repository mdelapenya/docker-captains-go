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
	Completed *bool  `json:"completed"`
	Order     *int   `json:"order_number"`
	Url       string `json:"url"`
}

func (t *Todo) URL() string {
	return fmt.Sprintf("/todos/%s", t.ID)
}

type TodoApp struct {
	Name            string
	Version         string
	UsersConnection string
}

var App = &TodoApp{
	Name:            "Todos",
	Version:         "1.0.0",
	UsersConnection: os.Getenv("POSTGRESQL_URL"),
}

func main() {
	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// register handler for the static files
	mux.Handle("/", FileSystemHandler)

	store, err := NewTodosRepository(context.Background(), App.UsersConnection)
	if err != nil {
		log.Fatalf("Cannot create a Todos repository. Exiting")
	}

	todosHandler := NewTodosHandler(store)

	handlers := logMiddleware(todosHandler)

	// Register the routes and handlers
	mux.Handle("/todos", handlers)
	mux.Handle("/todos/", handlers)

	// Run the server
	http.ListenAndServe(":8080", mux)
}
