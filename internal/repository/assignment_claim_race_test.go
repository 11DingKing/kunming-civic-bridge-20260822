package repository

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestConcurrentAssignmentClaimHasSingleLease(t *testing.T) {
	db := storeDB(t)
	defer db.Close()
	now := time.Now().UTC()
	if _, err := db.Exec(`INSERT INTO assignments(id,suggestion_id,department_id,status,version,assigned_at) VALUES('a','s','d','assigned',0,?)`, now.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, owner := range []string{"operator-a", "operator-b"} {
		wg.Add(1)
		go func(owner string) {
			defer wg.Done()
			results <- (AssignmentQueryStore{DB: db}).Claim(context.Background(), "a", owner, owner+"-token", now.Add(time.Hour), 0)
		}(owner)
		wg.Wait()
	}
	close(results)
	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("expected stale assignment claim to be rejected, got %d successes", successes)
	}
}
