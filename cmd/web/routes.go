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

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.authenticate)

	router.HandlerFunc(http.MethodGet, "/", app.dashboardView)

	// user
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	protected := dynamic.Append(app.requireAuthentication)

	// workspaces
	router.Handler(http.MethodGet, "/workspaces", protected.ThenFunc(app.workspaceListView))
	router.Handler(http.MethodGet, "/workspace/view/:id", protected.ThenFunc(app.workspaceView))
	router.Handler(http.MethodGet, "/workspace/create", protected.ThenFunc(app.workspaceCreate))
	router.Handler(http.MethodPost, "/workspace/create", protected.ThenFunc(app.workspaceCreatePost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
