package worker

import (
	"context"
	"github.com/11DingKing/kunming-civic-bridge/internal/service"
	"time"
)

type NotificationRunner struct {
	Service  service.NotificationService
	Interval time.Duration
	Send     func(context.Context, string, string) error
}

func (w NotificationRunner) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			w.tick(ctx, now)
		}
	}
}
func (w NotificationRunner) tick(ctx context.Context, now time.Time) {
	rows, e := w.Service.DB.QueryContext(ctx, `SELECT id,user_id,payload FROM notifications WHERE status='pending' AND next_attempt_at<=? ORDER BY next_attempt_at LIMIT 20`, now.Format(time.RFC3339Nano))
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id, user, payload string
		if e = rows.Scan(&id, &user, &payload); e != nil {
			continue
		}
		if e = w.Send(ctx, user, payload); e != nil {
			_ = w.Service.Fail(ctx, id, e.Error(), now)
		} else {
			_ = w.Service.Ack(ctx, id)
		}
	}
}
