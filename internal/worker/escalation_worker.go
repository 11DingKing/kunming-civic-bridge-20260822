package worker

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/service"
	"time"
)

type EscalationRunner struct {
	Service  service.EscalationService
	Interval time.Duration
}

func (w EscalationRunner) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			items, e := w.Service.FindDue(ctx, now)
			if e != nil {
				continue
			}
			for _, item := range items {
				_ = w.Service.Mark(ctx, item.SuggestionID)
			}
		}
	}
}

var _ *sql.DB
