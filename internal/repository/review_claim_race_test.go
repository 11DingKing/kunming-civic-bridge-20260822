package repository

import (
	"context"
	"sync"
	"testing"
)

func TestConcurrentReviewClaimHasSingleOwner(t *testing.T) {
	db := repoDB(t)
	defer db.Close()
	if err := (ReviewStore{DB: db}).Create(context.Background(), "s"); err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, actor := range []string{"reviewer-a", "reviewer-b"} {
		wg.Add(1)
		go func(actor string) {
			defer wg.Done()
			<-start
			results <- (ReviewStore{DB: db}).Claim(context.Background(), "s", actor, 0)
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
		t.Fatalf("expected one successful reviewer claim, got %d", successes)
	}
	var owner string
	if err := db.QueryRow(`SELECT reviewer_id FROM reviews WHERE suggestion_id='s'`).Scan(&owner); err != nil {
		t.Fatal(err)
	}
	if owner != "reviewer-a" && owner != "reviewer-b" {
		t.Fatalf("unexpected owner %q", owner)
	}
}
