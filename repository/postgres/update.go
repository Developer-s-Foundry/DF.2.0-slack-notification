package postgres

import (
	"context"
	"fmt"
	"time"
)

func (p *PostgresConn) UpdateTask(t Task) error {
	query := `
		UPDATE tasks
		SET
			name         = $1,
			status       = $2,
			description  = $3,
			assigned_to  = $4,
			expires_at   = $5,
			updated_at   = $6
		WHERE id = $7
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := p.Conn.Exec(
		ctx,
		query,
		t.Name,
		t.Status,
		t.Description,
		t.AssignedTo,
		t.ExpiresAt,
		t.UpdatedAt,
		t.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	// Check if any rows were affected to see if the update was successful.
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no task found with ID %s", t.ID)
	}

	return nil
}
