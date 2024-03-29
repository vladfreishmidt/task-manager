package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/vladfreishmidt/task-manager/internal/models"
)

type Partials struct {
	Header  bool
	Sidebar bool
	Action  bool
}

type UserInfo struct {
	Name  string
	Email string
}

type templateData struct {
	CurrentYear     int
	UserInfo        *UserInfo
	Partials        *Partials
	Workspace       *models.Workspace
	Workspaces      []*models.Workspace
	Project         *models.Project
	Projects        []*models.Project
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
