package models

import (
	"database/sql"
	"errors"
	"time"
)

type Workspace struct {
	WorkspaceID int
	OwnerID     int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type WorkspaceModel struct {
	DB *sql.DB
}

func (m *WorkspaceModel) Insert(ownerID int, name string, description string) (int, error) {
	stmt := "INSERT INTO workspaces (owner_id, name, description) VALUES(?, ?, ?)"

	result, err := m.DB.Exec(stmt, ownerID, name, description)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (m *WorkspaceModel) All(id int) ([]*Workspace, error) {
	stmt := `SELECT w.workspace_id, w.owner_id, w.name, w.description, w.created_at, w.updated_at FROM workspaces w 
	LEFT JOIN user_workspace uw ON w.workspace_id = uw.workspace_id 
	WHERE (uw.user_id = ? OR w.owner_id = ?);`

	rows, err := m.DB.Query(stmt, id, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	workspaces := []*Workspace{}

	for rows.Next() {
		w := &Workspace{}

		err = rows.Scan(&w.WorkspaceID, &w.OwnerID, &w.Name, &w.Description, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return nil, err
		}

		workspaces = append(workspaces, w)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (m *WorkspaceModel) Get(id int) (*Workspace, error) {
	stmt := `SELECT workspace_id, name, description, created_at FROM workspaces WHERE workspace_id = ?`

	row := m.DB.QueryRow(stmt, id)

	w := &Workspace{}

	err := row.Scan(&w.WorkspaceID, &w.Name, &w.Description, &w.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return w, nil
}
