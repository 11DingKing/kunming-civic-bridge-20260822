package service

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
)

// TestClaimAssignmentSerialSucceeds verifies the single-operator path is
// unaffected by the concurrency guard: one operator takes the lease cleanly.
func TestClaimAssignmentSerialSucceeds(t *testing.T) {
	db := claimAssignmentDB(t)
	defer db.Close()
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	w := Workflow{DB: db, Clock: func() time.Time { return now }}
	if e := w.ClaimAssignment(context.Background(), "s", "op-a", "token-a", time.Hour); e != nil {
		t.Fatalf("single claim err=%v", e)
	}
	var assignee, status, token string
	db.QueryRow(`SELECT COALESCE(assignee_id,''),status,COALESCE(lease_token,'') FROM assignments WHERE suggestion_id='s'`).Scan(&assignee, &status, &token)
	if assignee != "op-a" || status != "claimed" || token != "token-a" {
		t.Fatalf("assignment not claimed by op-a: assignee=%s status=%s token=%s", assignee, status, token)
	}
}

// TestClaimAssignmentConcurrentSingleWinner verifies that when two operators
// race on the same assignable dispatch, only one obtains a valid lease and the
// loser receives an explicit conflict instead of silently overwriting it.
func TestClaimAssignmentConcurrentSingleWinner(t *testing.T) {
	db := claimAssignmentDB(t)
	defer db.Close()
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	w := Workflow{DB: db, Clock: func() time.Time { return now }}

	var aErr, bErr error
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); aErr = w.ClaimAssignment(context.Background(), "s", "op-a", "token-a", time.Hour) }()
	go func() { defer wg.Done(); bErr = w.ClaimAssignment(context.Background(), "s", "op-b", "token-b", time.Hour) }()
	wg.Wait()

	if aErr == nil && bErr == nil {
		t.Fatal("both concurrent claims succeeded; expected exactly one conflict")
	}
	if aErr != nil && bErr != nil {
		t.Fatalf("both concurrent claims failed; expected one winner, got a=%v b=%v", aErr, bErr)
	}
	loser := aErr
	if bErr != nil {
		loser = bErr
	}
	if !errors.Is(loser, domain.ErrConflict) {
		t.Fatalf("loser expected conflict, got %v", loser)
	}

	var assignee, token string
	var wins int
	db.QueryRow(`SELECT COALESCE(assignee_id,''),COALESCE(lease_token,'') FROM assignments WHERE suggestion_id='s'`).Scan(&assignee, &token)
	db.QueryRow(`SELECT COUNT(*) FROM assignments WHERE suggestion_id='s' AND assignee_id IS NOT NULL`).Scan(&wins)
	if wins != 1 {
		t.Fatalf("expected exactly one assignment with an owner, got %d", wins)
	}
	if assignee != "op-a" && assignee != "op-b" {
		t.Fatalf("unexpected winner assignee=%s", assignee)
	}
	if token == "" {
		t.Fatal("winner has empty lease token")
	}
}

// claimAssignmentDB seeds an assignable dispatch that ClaimAssignment reads.
// A shared in-memory database keeps all connections on the same backing store so
// two goroutines can hold transactions concurrently and exercise the race.
func claimAssignmentDB(t *testing.T) *sql.DB {
	db, e := sql.Open("sqlite", "file::memory:?cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	if e = platform.Migrate(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(8)
	now := "2026-08-19T00:00:00Z"
	if _, e = db.Exec(`INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES('u','u','1','x','云南/昆明/西山/马街',?);INSERT INTO departments(id,name,scope,active,created_at) VALUES('d','部门','云南/昆明/西山/马街',1,?)`, now, now); e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec(`INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c','活动',?,?, 'open',?);INSERT INTO suggestions(id,campaign_id,author_id,title,body,scope,status,version,created_at,updated_at) VALUES('s','c','u','标题','正文','云南/昆明/西山/马街','assigned',0,?,?)`, now, now, now, now, now); e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec(`INSERT INTO assignments(id,suggestion_id,department_id,status,version,assigned_at) VALUES('a','s','d','assigned',0,?)`, now); e != nil {
		t.Fatal(e)
	}
	return db
}
