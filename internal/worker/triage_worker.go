package worker

import (
	"context"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/service"
	"time"
)

type TriageRunner struct {
	Service  service.TriageService
	Interval time.Duration
	Actor    string
}

func (w TriageRunner) Run(ctx context.Context, scope string) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ids, e := w.Service.Pending(ctx, scope)
			if e != nil {
				continue
			}
			for _, id := range ids {
				_ = w.Service.Classify(ctx, id, w.Actor, domain.TriageResult{Category: "民生", Priority: "normal", Reason: "待研判"})
			}
		}
	}
}
