package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoRows = pgx.ErrNoRows

type Pool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

type DB struct {
	Pool    Pool
	pgxPool *pgxpool.Pool
}

func New(databaseURL string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if pingErr := pool.Ping(context.Background()); pingErr != nil {
		return nil, fmt.Errorf("failed to ping database: %w", pingErr)
	}

	return &DB{Pool: pool, pgxPool: pool}, nil
}

func (db *DB) PgxPool() *pgxpool.Pool {
	return db.pgxPool
}

func (db *DB) Close() {
	db.Pool.Close()
}
