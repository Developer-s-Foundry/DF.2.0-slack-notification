package postgres

import (
	"context"
	"fmt"
	"time"
)

func (p *PostgresConn) Insert(t Task) error {
	query := `
		INSERT INTO tasks (id, name, status, description, assigned_to, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := p.Conn.QueryRow(
		ctx, query, t.ID, t.Name,
		t.Status, t.Description,
		t.AssignedTo, t.ExpiresAt,
	).Scan(&t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}
