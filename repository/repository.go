package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"

	"github.com/Drozd0f/csv-app/db"
)

const uniqueConstraintCode = "23505"

type Repository struct {
	conn *pgx.Conn
	q    *db.Queries
}

func New(ctx context.Context, dbURI string) (*Repository, error) {
	conn, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	q := db.New(conn)
	if err != nil {
		return nil, fmt.Errorf("db prepare: %w", err)
	}

	return &Repository{
		conn: conn,
		q:    q,
	}, nil
}
