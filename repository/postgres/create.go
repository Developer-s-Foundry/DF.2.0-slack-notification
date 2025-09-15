package postgres

import (
	"context"
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Description string    `json:"description,omitempty"`
	AssignedTo  string    `json:"assigned_to,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (p *PostgresConn) Create() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			status VARCHAR(100) NOT NULL DEFAULT 'pending',
			description TEXT NOT NULL,
			assigned_to VARCHAR(100),
			expires_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
	)	
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := p.Conn.Exec(ctx, query)
	return err
}
