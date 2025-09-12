package postgres

import (
	"context"
	"time"
)

func (p *PostgresConn) Create() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			name TEXT,
			status VARCHAR(100) DEFAULT 'pending',
			description TEXT,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
	)	
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := p.Conn.Exec(ctx, query)
	return err
}
