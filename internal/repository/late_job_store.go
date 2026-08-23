package repository

import (
	"context"
	"database/sql"
	"time"
)

type LateJobStore struct{ DB *sql.DB }

func (s LateJobStore) Count(ctx context.Context, now time.Time) (int, error) {
	var n int
	e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE status IN ('pending','running') AND run_after<?`, now.Format(time.RFC3339Nano)).Scan(&n)
	return n, e
}
func (s LateJobStore) Reschedule(ctx context.Context, now time.Time) error {
	_, e := s.DB.ExecContext(ctx, `UPDATE jobs SET run_after=?,updated_at=? WHERE status='pending' AND run_after<?`, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	return e
}
