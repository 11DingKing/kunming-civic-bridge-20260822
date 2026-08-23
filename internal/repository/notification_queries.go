package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"time"
)

type NotificationQueryStore struct{ DB *sql.DB }
type NotificationRecord struct {
	ID, UserID, SuggestionID, Kind, Payload, Status string
	Attempts                                        int
	NextAttemptAt                                   time.Time
}

func (n NotificationQueryStore) Ready(ctx context.Context, now time.Time, limit int) ([]NotificationRecord, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	rows, e := n.DB.QueryContext(ctx, `SELECT id,user_id,COALESCE(suggestion_id,''),kind,payload,status,attempts,next_attempt_at FROM notifications WHERE status='pending' AND next_attempt_at<=? ORDER BY next_attempt_at LIMIT ?`, now.Format(time.RFC3339Nano), limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []NotificationRecord
	for rows.Next() {
		var x NotificationRecord
		var next string
		if e = rows.Scan(&x.ID, &x.UserID, &x.SuggestionID, &x.Kind, &x.Payload, &x.Status, &x.Attempts, &next); e != nil {
			return nil, e
		}
		x.NextAttemptAt, _ = time.Parse(time.RFC3339Nano, next)
		out = append(out, x)
	}
	return out, rows.Err()
}
func (n NotificationQueryStore) MarkSent(ctx context.Context, id string) error {
	res, e := n.DB.ExecContext(ctx, `UPDATE notifications SET status='sent' WHERE id=?`, id)
	if e != nil {
		return e
	}
	nRows, _ := res.RowsAffected()
	if nRows != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (n NotificationQueryStore) Reschedule(ctx context.Context, id, reason string, now time.Time) error {
	_, e := n.DB.ExecContext(ctx, `UPDATE notifications SET attempts=attempts+1,last_error=?,next_attempt_at=? WHERE id=? AND status='pending'`, reason, now.Add(30*time.Second).Format(time.RFC3339Nano), id)
	return e
}
