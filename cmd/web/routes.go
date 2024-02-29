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

	router.HandlerFunc(http.MethodGet, "/", app.dashboardView)

	// workspaces
	router.HandlerFunc(http.MethodGet, "/workspaces", app.workspaceListView)
	router.HandlerFunc(http.MethodGet, "/workspace/view/:id", app.workspaceView)
	router.HandlerFunc(http.MethodGet, "/workspace/create", app.workspaceCreate)
	router.HandlerFunc(http.MethodPost, "/workspace/create", app.workspaceCreatePost)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
