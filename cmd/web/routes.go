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

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.HandlerFunc(http.MethodGet, "/", app.dashboardView)

	// workspaces
	router.Handler(http.MethodGet, "/workspaces", dynamic.ThenFunc(app.workspaceListView))
	router.Handler(http.MethodGet, "/workspace/view/:id", dynamic.ThenFunc(app.workspaceView))
	router.Handler(http.MethodGet, "/workspace/create", dynamic.ThenFunc(app.workspaceCreate))
	router.Handler(http.MethodPost, "/workspace/create", dynamic.ThenFunc(app.workspaceCreatePost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
