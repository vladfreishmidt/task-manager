package main

import (
	"net/http"
)

func (app *application) dashboardView(w http.ResponseWriter, r *http.Request) {

	partials := &Partials{Sidebar: true}

	data := app.newTemplateData(r)
	data.Partials = partials

	app.render(w, http.StatusOK, "dashboard-new.tmpl", data)
}
