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
	case r.Method == http.MethodPut && TodoReWithID.MatchString(r.URL.Path):
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

func (h *TodosHandler) List(w http.ResponseWriter, r *http.Request)      {}
func (h *TodosHandler) FindByID(w http.ResponseWriter, r *http.Request)  {}
func (h *TodosHandler) Update(w http.ResponseWriter, r *http.Request)    {}
func (h *TodosHandler) Delete(w http.ResponseWriter, r *http.Request)    {}
func (h *TodosHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {}
