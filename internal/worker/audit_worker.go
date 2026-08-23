package worker

import (
	"context"
	"github.com/11DingKing/kunming-civic-bridge/internal/service"
	"time"
)

type AuditRunner struct {
	Service  service.ExportService
	Interval time.Duration
}

func (w AuditRunner) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			continue
		}
	}
}
