package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

func TestConcurrentReviewCompletionHasSingleWinner(t *testing.T) {
	db := workflowSeed(t)
	defer db.Close()
	now := time.Now().UTC()
	if _, err := db.Exec(`INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c','活动',?,?, 'open',?)`, now.Add(-time.Hour).Format(time.RFC3339Nano), now.Add(time.Hour).Format(time.RFC3339Nano), now.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO suggestions(id,campaign_id,author_id,title,body,scope,status,version,created_at,updated_at) VALUES('s','c','u','标题','正文','云南/昆明/西山/马街','submitted',0,?,?)`, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}
	if err := (repository.ReviewStore{DB: db}).Create(context.Background(), "s"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`UPDATE reviews SET reviewer_id='reviewer',decision='claimed' WHERE suggestion_id='s'`); err != nil {
		t.Fatal(err)
	}
	w := ReviewWorkflow{DB: db, Reviews: repository.ReviewStore{DB: db}, Now: func() time.Time { return now }}
	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, actor := range []string{"reviewer-a", "reviewer-b"} {
		wg.Add(1)
		go func(actor string) {
			defer wg.Done()
			<-start
			results <- w.Complete(context.Background(), "s", actor, domain.DecisionAccept, "处理意见")
		}(actor)
	}
	close(start)
	wg.Wait()
	close(results)
	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("expected one completion winner, got %d", successes)
	}
	var events int
	if err := db.QueryRow(`SELECT COUNT(*) FROM suggestion_events WHERE suggestion_id='s'`).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if events != 1 {
		t.Fatalf("expected one completion event, got %d", events)
	}
}
