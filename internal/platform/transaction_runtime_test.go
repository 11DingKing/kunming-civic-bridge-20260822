package platform

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func TestTransactionRuntimeRollsBack(t *testing.T) {
	db, e := sql.Open("sqlite", ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec(`CREATE TABLE values_table(value TEXT)`); e != nil {
		t.Fatal(e)
	}
	runner := TransactionRuntime{DB: db}
	sentinel := errors.New("stop")
	if e = runner.Run(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, e := tx.ExecContext(ctx, `INSERT INTO values_table(value) VALUES('x')`)
		if e != nil {
			return e
		}
		return sentinel
	}); !errors.Is(e, sentinel) {
		t.Fatal(e)
	}
	var n int
	if e = db.QueryRow(`SELECT COUNT(*) FROM values_table`).Scan(&n); e != nil || n != 0 {
		t.Fatalf("rows=%d err=%v", n, e)
	}
}
func TestTransactionRuntimeCommits(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	db.Exec(`CREATE TABLE values_table(value TEXT)`)
	runner := TransactionRuntime{DB: db}
	if e := runner.Run(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, e := tx.ExecContext(ctx, `INSERT INTO values_table(value) VALUES('ok')`)
		return e
	}); e != nil {
		t.Fatal(e)
	}
	var n int
	db.QueryRow(`SELECT COUNT(*) FROM values_table`).Scan(&n)
	if n != 1 {
		t.Fatal("commit missing")
	}
}
