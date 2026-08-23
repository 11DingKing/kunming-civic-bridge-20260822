package worker

import (
	"context"
	"database/sql"
	"time"
)

type RecoveryRunner struct{ DB *sql.DB }

func (w RecoveryRunner) Recover(ctx context.Context, now time.Time) error {
	_, e := w.DB.ExecContext(ctx, `UPDATE jobs SET status='pending',lease_until=NULL,run_after=?,updated_at=? WHERE status='running' AND lease_until<=?`, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	return e
}
func (w RecoveryRunner) Pending(ctx context.Context) (int, error) {
	var n int
	e := w.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE status='pending'`).Scan(&n)
	return n, e
}
