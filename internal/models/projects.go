package models

import (
	"database/sql"
	"time"
)

type Project struct {
	ID          int
	Title       string
	Description string
	Created     time.Time
}

type ProjectModel struct {
	DB *sql.DB
}

func (m *ProjectModel) Insert(title string, description string) (int, error) {
	stmt := "INSERT INTO projects (title, description, created) VALUES(?, ?, UTC_TIMESTAMP())"

	result, err := m.DB.Exec(stmt, title, description)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (m *ProjectModel) Get(id int) (*Project, error) {
	return nil, nil
}
func (m *ProjectModel) Latest() ([]*Project, error) {
	return nil, nil
}
