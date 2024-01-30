package main

import "net/http"

var FileSystemHandler http.Handler = http.StripPrefix("/", http.FileServer(http.Dir("public")))

type TodosHandler struct{}

func (h *TodosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my todos page"))
}
