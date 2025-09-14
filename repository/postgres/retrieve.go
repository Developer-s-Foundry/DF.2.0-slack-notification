package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

func (p *PostgresConn) GetTaskByID(context context.Context, taskID string) (*Task, error) {
	query := `
		SELECT id, name, status, description, assigned_to, expires_at, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	t := &Task{}

	err := p.Conn.QueryRow(
		context, query, taskID,
	).Scan(&t.ID,
		&t.Name,
		&t.Status,
		&t.Description,
		&t.AssignedTo,
		&t.Expires_at,
		&t.CreatedAt,
		&t.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with ID %s not found", taskID)
		}
		return nil, fmt.Errorf("failed to retrieve task by ID: %w", err)
	}
	return t, nil
}
