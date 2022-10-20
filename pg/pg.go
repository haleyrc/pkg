package pg

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Row interface {
	Scan(args ...interface{}) error
}

type Tx interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
}

func Connect(ctx context.Context, path string) (*Database, error) {
	conn, err := sqlx.ConnectContext(ctx, "postgres", path)
	if err != nil {
		return nil, fmt.Errorf("connect failed: %w", err)
	}

	return &Database{conn: conn}, nil
}

type Database struct {
	conn *sqlx.DB
}

func (db *Database) WithTx(ctx context.Context, f func(ctx context.Context, tx Tx) error) error {
	tx, err := db.conn.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}

	if err := f(ctx, ttx{tx: tx}); err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func (db *Database) Ping() error {
	if err := db.conn.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

type ttx struct {
	tx *sqlx.Tx
}

func (tx ttx) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := tx.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}

func (tx ttx) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	return tx.tx.QueryRowContext(ctx, query, args...)
}
