package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"github.com/vladfreishmidt/task-manager/internal/models"
)

type workspaceCreateForm struct {
	Name        string
	Description string
	FieldErrors map[string]string
}

func (app *application) workspaceCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := workspaceCreateForm{
		Name:        r.PostForm.Get("workspace-name"),
		Description: r.PostForm.Get("workspace-description"),
		FieldErrors: map[string]string{},
	}

	ownerID := 1 // hardcoded userID

	if strings.TrimSpace(form.Name) == "" {
		form.FieldErrors["name"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Name) > 100 {
		form.FieldErrors["name"] = "This field cannot be more than 100 characters long"
	}

	if utf8.RuneCountInString(form.Description) > 250 {
		form.FieldErrors["description"] = "This field cannot be more than 250 characters long"
	}

	if len(form.FieldErrors) > 0 {
		partials := &Partials{Sidebar: true}
		data := app.newTemplateData(r)
		data.Form = form
		data.Partials = partials
		app.render(w, http.StatusUnprocessableEntity, "workspace-create.tmpl", data)
		return
	}

	id, err := app.workspaces.Insert(ownerID, form.Name, form.Description)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/workspace/view/%d", id), http.StatusSeeOther)
}

func (app *application) workspaceCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	partials := &Partials{Sidebar: false}
	data.Partials = partials
	data.Form = workspaceCreateForm{}

	app.render(w, http.StatusOK, "workspace-create.tmpl", data)
}

func (app *application) workspaceView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	workspace, err := app.workspaces.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	projects, err := app.projects.All(workspace.WorkspaceID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	partials := &Partials{Sidebar: true}

	data := app.newTemplateData(r)
	data.Partials = partials
	data.Workspace = workspace
	data.Projects = projects

	app.render(w, http.StatusOK, "workspace-details.tmpl", data)
}

func (app *application) workspaceListView(w http.ResponseWriter, r *http.Request) {

	hardcodedUserID := 1 // hardcoded user id for dev purpose

	workspaces, err := app.workspaces.All(hardcodedUserID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	partials := &Partials{Sidebar: false}

	data := app.newTemplateData(r)
	data.Partials = partials
	data.Workspaces = workspaces

	app.render(w, http.StatusOK, "workspace-list.tmpl", data)
}
