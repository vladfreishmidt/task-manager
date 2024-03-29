package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/vladfreishmidt/task-manager/internal/models"
	"github.com/vladfreishmidt/task-manager/internal/validator"
)

type workspaceCreateForm struct {
	Name                string `form:"workspace-name"`
	Description         string `form:"workspace-description"`
	validator.Validator `form:"-"`
}

func (app *application) workspaceCreatePost(w http.ResponseWriter, r *http.Request) {
	var form workspaceCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	ownerID := r.Context().Value(userID).(int)

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 100), "name", "This field cannot be more than 100 characters long")
	form.CheckField(validator.MaxChars(form.Description, 255), "description", "This field cannot be more than 255 characters long")

	if !form.Valid() {
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

	app.sessionManager.Put(r.Context(), "flash", "Workspace successfully created.")

	http.Redirect(w, r, fmt.Sprintf("/workspace/view/%d", id), http.StatusSeeOther)
}

func (app *application) workspaceCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	partials := &Partials{Sidebar: true, Header: true}

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
	userID := r.Context().Value(userID).(int)

	workspaces, err := app.workspaces.All(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	partials := &Partials{Sidebar: true, Header: true}

	data := app.newTemplateData(r)
	data.Partials = partials
	data.Workspaces = workspaces

	app.render(w, http.StatusOK, "workspace-list.tmpl", data)
}
