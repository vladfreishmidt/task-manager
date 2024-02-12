package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Task manager"))
}

func projectView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific project..."))
}

func projectCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new project..."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/project/view", projectView)
	mux.HandleFunc("/project/create", projectCreate)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)

}
