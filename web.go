package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

var (
	TodoRe       = regexp.MustCompile(`^/todos/*$`)
	TodoReWithID = regexp.MustCompile(`^/todos/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

var FileSystemHandler http.Handler = http.StripPrefix("/", http.FileServer(http.Dir("public")))

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

type TodosHandler struct {
	store todoStore
}

func NewTodosHandler(s todoStore) *TodosHandler {
	return &TodosHandler{
		store: s,
	}
}

func (h *TodosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && TodoRe.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	case r.Method == http.MethodGet && TodoRe.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodGet && TodoReWithID.MatchString(r.URL.Path):
		h.FindByID(w, r)
		return
	case r.Method == http.MethodPatch && TodoReWithID.MatchString(r.URL.Path):
		h.Update(w, r)
		return
	case r.Method == http.MethodDelete && TodoReWithID.MatchString(r.URL.Path):
		h.Delete(w, r)
		return
	case r.Method == http.MethodDelete && TodoRe.MatchString(r.URL.Path):
		h.DeleteAll(w, r)
		return
	default:
		return
	}
}

func (h *TodosHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Todo object that will be populated from JSON payload
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		log.Printf("Cannot decode the request body: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

	// Call the store to add the recipe
	if err := h.store.Create(r.Context(), &todo); err != nil {
		log.Printf("Cannot create the todo: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

	// Set the status code to 200
	w.WriteHeader(http.StatusOK)
}

func (h *TodosHandler) List(w http.ResponseWriter, r *http.Request) {
	resources, err := h.store.List(r.Context())
	if err != nil {
		log.Printf("Cannot retrieve todos: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(resources)
	if err != nil {
		log.Printf("Cannot encode todos to JSON: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *TodosHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	// Extract the resource ID using a regex
	matches := TodoReWithID.FindStringSubmatch(r.URL.Path)
	// Expect matches to be length >= 2 (full string + 1 matching group)
	if len(matches) < 2 {
		log.Printf("Cannot parse the request URL: %v", r.URL.Path)
		InternalServerErrorHandler(w, r)
		return
	}

	// Retrieve todo from the store
	todo, err := h.store.FindByID(r.Context(), matches[1])
	if err != nil {
		// Special case of not-found Error
		if err == ErrNotFound {
			NotFoundHandler(w, r)
			return
		}

		// Every other error
		log.Printf("Cannot retrieve the todo: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

	// Convert the struct into JSON payload
	jsonBytes, err := json.Marshal(todo)
	if err != nil {
		log.Printf("Cannot encode todo to JSON: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

	// Write the results
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *TodosHandler) Update(w http.ResponseWriter, r *http.Request) {}

func (h *TodosHandler) Delete(w http.ResponseWriter, r *http.Request) {
	matches := TodoReWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		log.Printf("Cannot parse the request URL: %v", r.URL.Path)
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.DeleteByID(r.Context(), matches[1]); err != nil {
		// we do not want to return a 404 error if the todo is not found
		// to avoid leaking information about the existence of a resource
		log.Printf("Cannot delete the todo: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TodosHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if err := h.store.DeleteAll(r.Context()); err != nil {
		log.Printf("Cannot delete all todos: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
