package repository

import (
	"context"
	"errors"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"sync"
	"testing"
)

// reviewRow seeds an unassigned review row (version 0) for suggestion "s".
func reviewRow(t *testing.T) *ReviewStore {
	t.Helper()
	db := storeDB(t)
	t.Cleanup(func() { db.Close() })
	s := ReviewStore{DB: db}
	if e := s.Create(context.Background(), "s"); e != nil {
		t.Fatalf("create review: %v", e)
	}
	return &s
}

func TestReviewStoreClaimSingle(t *testing.T) {
	s := reviewRow(t)
	// Normal single claim still completes.
	if e := s.Claim(context.Background(), "s", "rev1", 0); e != nil {
		t.Fatalf("claim: %v", e)
	}
	var reviewer, decision string
	var version int
	if e := s.DB.QueryRow(`SELECT reviewer_id,decision,version FROM reviews WHERE suggestion_id='s'`).Scan(&reviewer, &decision, &version); e != nil {
		t.Fatal(e)
	}
	if reviewer != "rev1" || decision != "claimed" || version != 1 {
		t.Fatalf("got reviewer=%s decision=%s version=%d", reviewer, decision, version)
	}
}

func TestReviewStoreClaimConflictOnStaleVersion(t *testing.T) {
	s := reviewRow(t)
	// First reviewer wins, bumping the version.
	if e := s.Claim(context.Background(), "s", "rev1", 0); e != nil {
		t.Fatalf("first claim: %v", e)
	}
	// Second reviewer still holds the stale version=0 snapshot. The atomic
	// compare-and-set matches zero rows (version is now 1, decision 'claimed').
	if e := s.Claim(context.Background(), "s", "rev2", 0); !errors.Is(e, domain.ErrConflict) {
		t.Fatalf("second claim want conflict, got %v", e)
	}
	var reviewer string
	s.DB.QueryRow(`SELECT reviewer_id FROM reviews WHERE suggestion_id='s'`).Scan(&reviewer)
	if reviewer != "rev1" {
		t.Fatalf("ownership leaked to %s", reviewer)
	}
}

func TestReviewStoreClaimConcurrentOneWinner(t *testing.T) {
	db := storeDB(t)
	t.Cleanup(func() { db.Close() })
	// Allow genuine parallel writers so the compare-and-set is exercised, not
	// merely serialized by a single connection.
	db.SetMaxOpenConns(4)
	s := ReviewStore{DB: db}
	if e := s.Create(context.Background(), "s"); e != nil {
		t.Fatalf("create review: %v", e)
	}
	var wg sync.WaitGroup
	results := make([]error, 4)
	for i := range results {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			actor := "rev" + string(rune('1'+idx))
			results[idx] = s.Claim(context.Background(), "s", actor, 0)
		}(i)
	}
	wg.Wait()
	wins, conflicts := 0, 0
	var winner string
	for i, e := range results {
		switch {
		case e == nil:
			wins++
			winner = "rev" + string(rune('1'+i))
		case errors.Is(e, domain.ErrConflict):
			conflicts++
		default:
			t.Fatalf("unexpected err: %v", e)
		}
	}
	if wins != 1 || conflicts != 3 {
		t.Fatalf("want 1 win 3 conflicts, got wins=%d conflicts=%d", wins, conflicts)
	}
	var reviewer string
	s.DB.QueryRow(`SELECT reviewer_id FROM reviews WHERE suggestion_id='s'`).Scan(&reviewer)
	if reviewer != winner {
		t.Fatalf("owner=%s winner=%s", reviewer, winner)
	}
}
