package main

import "net/http"

var FileSystemHandler http.Handler = http.StripPrefix("/", http.FileServer(http.Dir("public")))

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
