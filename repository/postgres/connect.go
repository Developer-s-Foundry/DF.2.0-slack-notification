package postgres

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConn struct {
	Conn *pgxpool.Pool
}

func NewPostgresConn(conn *pgxpool.Pool) *PostgresConn {
	return &PostgresConn{Conn: conn}
}

func ConnectPostgres(uri, password, port, host, database, user, sslmode string) (*PostgresConn, error) {
	var err error
	var portInt int
	if port != "" {
		portInt, err = strconv.Atoi(port)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}

	if uri == "" {
		uri = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", user, password, host, portInt, database, sslmode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pgx, err := pgxpool.New(ctx, uri)

	if err != nil {
		log.Printf("Database connection failed: %v", err)
		return nil, err
	}

	if err := pgx.Ping(ctx); err != nil {
		log.Printf("unable to ping database: %v", err)
		return nil, err
	}
	log.Println("Database connected successfully")

	conn := NewPostgresConn(pgx)

	err = conn.Create()

	if err != nil {
		log.Printf("unable to create table: %v", err)
		return nil, err
	}

	log.Println("task table create successfully")
	return conn, nil
}
