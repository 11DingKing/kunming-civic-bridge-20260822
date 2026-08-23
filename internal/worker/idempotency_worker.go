package worker

import (
	"context"
	"github.com/11DingKing/kunming-civic-bridge/internal/service"
	"time"
)

type IdempotencyRunner struct {
	Service  service.IdempotencyService
	Interval time.Duration
}

func (w IdempotencyRunner) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			_, _ = w.Service.Store.Purge(ctx, now)
		}
	}
}
