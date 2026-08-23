package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"testing"
	"time"
)

func repoDB(t *testing.T) *sql.DB {
	db, e := sql.Open("sqlite", ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(1)
	if e = platform.Migrate(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	now := "2026-08-19T00:00:00Z"
	_, e = db.Exec(`INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES('u','u','1','x','云南/昆明/西山/马街',?);INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c','活动',?,?, 'open',?);INSERT INTO suggestions(id,campaign_id,author_id,title,body,scope,status,version,created_at,updated_at) VALUES('s','c','u','title','body','云南/昆明/西山/马街','submitted',0,?,?)`, now, "2026-08-18T00:00:00Z", "2026-08-30T00:00:00Z", now, now, now)
	if e != nil {
		t.Fatal(e)
	}
	return db
}
func TestSuggestionRepositoryVersionConflict(t *testing.T) {
	db := repoDB(t)
	defer db.Close()
	r := SuggestionRepository{DB: db}
	x, e := r.Find(context.Background(), "s")
	if e != nil {
		t.Fatal(e)
	}
	tx, e := db.Begin()
	if e != nil {
		t.Fatal(e)
	}
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	if e = r.TransitionTx(context.Background(), tx, x.ID, x.Status, "triaged", x.Version, now); e != nil {
		t.Fatal(e)
	}
	if e = tx.Commit(); e != nil {
		t.Fatal(e)
	}
	tx, e = db.Begin()
	if e != nil {
		t.Fatal(e)
	}
	if e = r.TransitionTx(context.Background(), tx, x.ID, x.Status, "triaged", x.Version, now); e == nil {
		t.Fatal("stale version accepted")
	}
	tx.Rollback()
}
func TestJobStoreRetryLifecycle(t *testing.T) {
	db := repoDB(t)
	defer db.Close()
	j := JobStore{DB: db}
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	if e := j.Enqueue(context.Background(), "notify", "payload", now); e != nil {
		t.Fatal(e)
	}
	id, kind, payload, e := j.Claim(context.Background(), now, time.Minute)
	if e != nil || kind != "notify" || payload != "payload" {
		t.Fatalf("claim %s %s %v", kind, payload, e)
	}
	if e = j.Finish(context.Background(), id, sql.ErrNoRows, 1, 3, now); e != nil {
		t.Fatal(e)
	}
	var status string
	if e = db.QueryRow(`SELECT status FROM jobs WHERE id=?`, id).Scan(&status); e != nil || status != "pending" {
		t.Fatalf("status=%s err=%v", status, e)
	}
}
func TestJobStoreClaimRejectsActiveLease(t *testing.T) {
	db := repoDB(t)
	defer db.Close()
	j := JobStore{DB: db}
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	if e := j.Enqueue(context.Background(), "notify", "payload", now); e != nil {
		t.Fatal(e)
	}
	if _, _, _, e := j.Claim(context.Background(), now, time.Minute); e != nil {
		t.Fatalf("first claim: %v", e)
	}
	// While the lease is still valid, a second worker must not re-claim the job.
	if _, _, _, e := j.Claim(context.Background(), now.Add(30*time.Second), time.Minute); !errors.Is(e, sql.ErrNoRows) && !errors.Is(e, domain.ErrConflict) {
		t.Fatalf("expected no rows/conflict, got %v", e)
	}
	var status, lease string
	if e := db.QueryRow(`SELECT status,lease_until FROM jobs WHERE kind='notify'`).Scan(&status, &lease); e != nil {
		t.Fatal(e)
	}
	if status != "running" {
		t.Fatalf("status=%s want running", status)
	}
	// After the lease expires the normal recovery path must re-claim it.
	later := now.Add(2 * time.Minute)
	id, kind, payload, e := j.Claim(context.Background(), later, time.Minute)
	if e != nil || kind != "notify" || payload != "payload" {
		t.Fatalf("reclaim after expiry: id=%s kind=%s payload=%s err=%v", id, kind, payload, e)
	}
	if id == "" {
		t.Fatal("empty job id reclaimed")
	}
}
