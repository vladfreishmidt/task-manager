package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/vladfreishmidt/task-manager/internal/models"
)

func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	projects, err := app.projects.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	partials := &Partials{Sidebar: false}

	data := app.newTemplateData(r)
	data.Projects = projects
	data.Partials = partials

	app.render(w, http.StatusOK, "dashboard.tmpl", data)
}

func (app *application) projectListView(w http.ResponseWriter, r *http.Request) {
	projects, err := app.projects.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	partials := &Partials{Sidebar: true, Action: true}

	app.render(w, http.StatusOK, "project-list.tmpl", &templateData{
		Projects: projects,
		Partials: partials,
	})
}

func (app *application) projectView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	project, err := app.projects.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	partials := &Partials{Sidebar: true}

	app.render(w, http.StatusOK, "project-details.tmpl", &templateData{
		Project:  project,
		Partials: partials,
	})
}

func (app *application) projectCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "Cabin Booking app"
	description := "A web app for booking an A-frame cabins"

	id, err := app.projects.Insert(title, description)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/project/view?id=%d", id), http.StatusSeeOther)
}
