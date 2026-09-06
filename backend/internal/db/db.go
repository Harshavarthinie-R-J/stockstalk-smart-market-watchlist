package db

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Database handle.
type DB struct {
	*sql.DB
}

func Open(url string) (*DB, error) {
	conn, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	if err := conn.PingContext(context.Background()); err != nil {
		conn.Close()
		return nil, err
	}

	return &DB{DB: conn}, nil
}
