package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// router
	mux := http.NewServeMux()

	// file server
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// routes
	mux.HandleFunc("/", app.dashboard)
	mux.HandleFunc("/projects", app.projectListView)
	mux.HandleFunc("/project/view", app.projectView)
	mux.HandleFunc("/project/create", app.projectCreate)

	return mux
}
