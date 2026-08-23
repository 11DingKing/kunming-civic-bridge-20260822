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
	var id, kind, payload string
	e = tx.QueryRowContext(ctx, `SELECT id,kind,payload FROM jobs WHERE status IN ('pending','running') AND run_after<=? ORDER BY run_after LIMIT 1`, now.Format(time.RFC3339Nano)).Scan(&id, &kind, &payload)
	if e != nil {
		return "", "", "", e
	}
	until := now.Add(lease).Format(time.RFC3339Nano)
	if _, e = tx.ExecContext(ctx, `UPDATE jobs SET status='running',lease_until=?,updated_at=? WHERE id=?`, until, now.Format(time.RFC3339Nano), id); e != nil {
		return "", "", "", e
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
