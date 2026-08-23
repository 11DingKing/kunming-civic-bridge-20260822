package worker

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"github.com/11DingKing/kunming-civic-bridge/internal/service"
	"testing"
	"time"
)

func TestTriageRunnerStops(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	platform.Migrate(context.Background(), db)
	ctx, cancel := context.WithCancel(context.Background())
	r := TriageRunner{Service: service.TriageService{DB: db}, Interval: time.Millisecond}
	done := make(chan struct{})
	go func() { r.Run(ctx, "scope"); close(done) }()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("triage runner did not stop")
	}
}
func TestIdempotencyRunnerStops(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	platform.Migrate(context.Background(), db)
	ctx, cancel := context.WithCancel(context.Background())
	r := IdempotencyRunner{Service: service.IdempotencyService{Store: repository.IdempotencyStore{DB: db}}, Interval: time.Millisecond}
	done := make(chan struct{})
	go func() { r.Run(ctx); close(done) }()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("idempotency runner did not stop")
	}
}
