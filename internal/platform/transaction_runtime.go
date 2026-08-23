package platform

import (
	"context"
	"database/sql"
	"fmt"
)

type TransactionRuntime struct{ DB *sql.DB }

func (t TransactionRuntime) Run(ctx context.Context, fn func(context.Context, *sql.Tx) error) error {
	tx, e := t.DB.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	if e = fn(ctx, tx); e != nil {
		return e
	}
	if e = tx.Commit(); e != nil {
		return fmt.Errorf("commit transaction: %w", e)
	}
	return nil
}
func (t TransactionRuntime) Read(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return t.DB.QueryContext(ctx, query, args...)
}
