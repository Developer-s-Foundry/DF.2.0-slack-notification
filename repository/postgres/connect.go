package postgres

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

type PostgresConn struct {
	Port     int
	Password string
	Host     string
	Database string
	User     string
	Uri      string
	SSLMode  string
}

func NewPostgresConn(uri, password, port, host, database, user, sslmode string) (*PostgresConn, error) {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &PostgresConn{Port: portInt,
		Password: password,
		Host:     host,
		Database: database,
		SSLMode:  sslmode,
		User:     user,
		Uri:      uri,
	}, nil
}

func ConnectPostgres(p *PostgresConn) (*pgx.Conn, error) {
	if p.Uri == "" {
		p.Uri = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", p.User, p.Password, p.Host, p.Port, p.Database, p.SSLMode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pgx, err := pgx.Connect(ctx, p.Uri)

	if err != nil {
		log.Printf("Database connection failed: %v", err)
		return nil, err
	}

	if err := pgx.Ping(ctx); err != nil {
		log.Printf("unable to ping database: %v", err)
		return nil, err
	}
	log.Println("Database connected successfully")
	return pgx, nil
}
