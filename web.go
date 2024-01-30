package main

import "net/http"

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
	w.Write([]byte("This is my todos page"))
}
