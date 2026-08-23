package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"time"
)

type JobStore struct{ DB *sql.DB }

func (j JobStore) Enqueue(ctx context.Context, kind, payload string, runAt time.Time) error {
	if e := domain.ValidateJob(kind, payload); e != nil {
		return e
	}
	_, e := j.DB.ExecContext(ctx, `INSERT INTO jobs(id,kind,payload,status,attempts,run_after,created_at,updated_at) VALUES(lower(hex(randomblob(16))),?,?,?,0,?,?,?)`, kind, payload, domain.JobPending, runAt.Format(time.RFC3339Nano), runAt.Format(time.RFC3339Nano), runAt.Format(time.RFC3339Nano))
	return e
}
func (j JobStore) Claim(ctx context.Context, now time.Time, lease time.Duration) (string, string, string, error) {
	tx, e := j.DB.BeginTx(ctx, nil)
	if e != nil {
		return "", "", "", e
	}
	defer tx.Rollback()
	nowStr := now.Format(time.RFC3339Nano)
	// A running job still inside its lease is owned by another worker and must
	// not be re-claimed; only pending jobs and running jobs whose lease has
	// expired are eligible. The UPDATE re-checks the same condition so that a
	// concurrent claim on the same row loses atomically.
	var id, kind, payload string
	e = tx.QueryRowContext(ctx, `SELECT id,kind,payload FROM jobs WHERE run_after<=? AND (status='pending' OR (status='running' AND (lease_until IS NULL OR lease_until<=?))) ORDER BY run_after LIMIT 1`, nowStr, nowStr).Scan(&id, &kind, &payload)
	if e != nil {
		return "", "", "", e
	}
	until := now.Add(lease).Format(time.RFC3339Nano)
	res, e := tx.ExecContext(ctx, `UPDATE jobs SET status='running',lease_until=?,updated_at=? WHERE id=? AND (status='pending' OR (status='running' AND (lease_until IS NULL OR lease_until<=?)))`, until, nowStr, id, nowStr)
	if e != nil {
		return "", "", "", e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return "", "", "", domain.ErrConflict
	}
	if e = tx.Commit(); e != nil {
		return "", "", "", e
	}
	return id, kind, payload, nil
}
func (j JobStore) Finish(ctx context.Context, id string, runErr error, attempt, max int, now time.Time) error {
	state := domain.NextJobState(domain.JobRunning, runErr, attempt, max)
	if runErr == nil {
		_, e := j.DB.ExecContext(ctx, `UPDATE jobs SET status='done',lease_until=NULL,updated_at=? WHERE id=?`, now.Format(time.RFC3339Nano), id)
		return e
	}
	if state == domain.JobDead {
		_, e := j.DB.ExecContext(ctx, `UPDATE jobs SET status='dead',lease_until=NULL,last_error=?,updated_at=? WHERE id=?`, runErr.Error(), now.Format(time.RFC3339Nano), id)
		return e
	}
	_, e := j.DB.ExecContext(ctx, `UPDATE jobs SET status='pending',attempts=attempts+1,lease_until=NULL,last_error=?,run_after=?,updated_at=? WHERE id=?`, runErr.Error(), now.Add(domain.Backoff(attempt)).Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), id)
	return e
}
