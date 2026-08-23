package worker

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
)

type Handler func(context.Context, *sql.Tx, string) error
type Worker struct {
	DB       *sql.DB
	Interval time.Duration
	Log      *slog.Logger
	Handlers map[string]Handler
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}
func (w *Worker) tick(ctx context.Context) {
	tx, e := w.DB.BeginTx(ctx, nil)
	if e != nil {
		return
	}
	defer tx.Rollback()
	row := tx.QueryRowContext(ctx, `SELECT id,kind,payload FROM jobs WHERE status='pending' AND run_after<=CURRENT_TIMESTAMP ORDER BY created_at LIMIT 1`)
	var id, kind, payload string
	if e = row.Scan(&id, &kind, &payload); e != nil {
		return
	}
	h := w.Handlers[kind]
	if h == nil {
		_, _ = tx.ExecContext(ctx, `UPDATE jobs SET status='failed',last_error='unknown job',updated_at=CURRENT_TIMESTAMP WHERE id=?`, id)
		_ = tx.Commit()
		return
	}
	if e = h(ctx, tx, payload); e != nil {
		_, _ = tx.ExecContext(ctx, `UPDATE jobs SET attempts=attempts+1,last_error=?,run_after=datetime('now','+30 seconds'),updated_at=CURRENT_TIMESTAMP WHERE id=?`, e.Error(), id)
	} else {
		_, _ = tx.ExecContext(ctx, `UPDATE jobs SET status='done',updated_at=CURRENT_TIMESTAMP WHERE id=?`, id)
	}
	_ = tx.Commit()
}
