package main

import "github.com/vladfreishmidt/task-manager/internal/models"

type templateData struct {
	Project  *models.Project
	Projects []*models.Project
}
