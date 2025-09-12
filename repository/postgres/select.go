package postgres

import (
	"context"
	"fmt"
)

func (p *PostgresConn) CheckDataExists() (bool, error) {
	var count int
	err := p.Conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check existing tasks: %w", err)
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}
