package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// router
	router := httprouter.New()

	// custom not found error handler
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// file server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// workspaces
	router.HandlerFunc(http.MethodGet, "/workspaces", app.workspaceListView)
	router.HandlerFunc(http.MethodGet, "/workspaces/:id", app.workspaceView)
	// router.HandlerFunc(http.MethodGet, "/workspaces/create", app.workspaceCreate)
	router.HandlerFunc(http.MethodPost, "/workspaces/create", app.workspaceCreatePost)

	// projects
	// mux.HandleFunc("/", app.dashboard)
	// mux.HandleFunc("/projects", app.projectListView)
	// mux.HandleFunc("/project/view", app.projectView)
	// mux.HandleFunc("/project/create", app.projectCreate)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
