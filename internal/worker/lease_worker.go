package worker

import (
	"context"
	"database/sql"
	"time"
)

type LeaseRunner struct{ DB *sql.DB }

func (w LeaseRunner) Release(ctx context.Context, now time.Time) error {
	_, e := w.DB.ExecContext(ctx, `UPDATE assignments SET status='assigned',assignee_id=NULL,lease_token=NULL,lease_until=NULL,version=version+1 WHERE status='claimed' AND lease_until IS NOT NULL AND lease_until<=?`, now.Format(time.RFC3339Nano))
	return e
}
func (w LeaseRunner) Active(ctx context.Context) (int, error) {
	var n int
	e := w.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM assignments WHERE status='claimed'`).Scan(&n)
	return n, e
}
