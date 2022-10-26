package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Repository struct {
	conn *sql.DB
}

func New(ctx context.Context, dbURI string) (*Repository, error) {
	conn, err := sql.Open("postgres", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	return &Repository{
		conn: conn,
	}, nil
}
