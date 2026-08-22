package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type NotificationService struct{ DB *sql.DB }

func (n NotificationService) Queue(ctx context.Context, user, suggestion, kind string, payload map[string]string) error {
	if user == "" || kind == "" {
		return domain.ErrInvalid
	}
	b, e := json.Marshal(payload)
	if e != nil {
		return e
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, e = n.DB.ExecContext(ctx, `INSERT INTO notifications(id,user_id,suggestion_id,kind,payload,status,attempts,next_attempt_at,created_at) VALUES(?,?,?,?,?,'pending',0,?,?)`, platform.ID(), user, suggestion, kind, string(b), now, now)
	return e
}
func (n NotificationService) Ack(ctx context.Context, id string) error {
	res, e := n.DB.ExecContext(ctx, `UPDATE notifications SET status='sent' WHERE id=? AND status='pending'`, id)
	if e != nil {
		return e
	}
	v, _ := res.RowsAffected()
	if v != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (n NotificationService) Fail(ctx context.Context, id string, reason string, now time.Time) error {
	_, e := n.DB.ExecContext(ctx, `UPDATE notifications SET attempts=attempts+1,last_error=?,next_attempt_at=? WHERE id=? AND status='pending'`, reason, now.Add(30*time.Second).Format(time.RFC3339Nano), id)
	return e
}
