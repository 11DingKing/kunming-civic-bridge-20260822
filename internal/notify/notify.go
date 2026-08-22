package notify

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type Queue struct{ DB *sql.DB }

func (q Queue) Enqueue(ctx context.Context, tx *sql.Tx, user, suggestion, kind, payload string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, e := tx.ExecContext(ctx, `INSERT INTO notifications(id,user_id,suggestion_id,kind,payload,status,attempts,next_attempt_at,created_at) VALUES(?,?,?,?,?,'pending',0,?,?)`, platform.ID(), user, suggestion, kind, payload, now, now)
	return e
}
