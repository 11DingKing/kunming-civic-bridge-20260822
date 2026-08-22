package worker

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"github.com/11DingKing/kunming-civic-bridge/internal/service"
	"testing"
	"time"
)

func TestNotificationRunnerCancellation(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	platform.Migrate(context.Background(), db)
	ctx, cancel := context.WithCancel(context.Background())
	runner := NotificationRunner{Service: service.NotificationService{DB: db}, Interval: time.Millisecond, Send: func(context.Context, string, string) error { return nil }}
	done := make(chan struct{})
	go func() { runner.Run(ctx); close(done) }()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runner did not stop")
	}
}
func TestEscalationRunnerCancellation(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	platform.Migrate(context.Background(), db)
	ctx, cancel := context.WithCancel(context.Background())
	runner := EscalationRunner{Service: service.EscalationService{DB: db}, Interval: time.Millisecond}
	done := make(chan struct{})
	go func() { runner.Run(ctx); close(done) }()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runner did not stop")
	}
}
