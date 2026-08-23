package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"testing"
	"time"
)

func reviewWorkflowDB(t *testing.T) *sql.DB {
	t.Helper()
	db, e := sql.Open("sqlite", ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(1)
	if e = platform.Migrate(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	now := "2026-08-19T00:00:00Z"
	if _, e = db.Exec(`INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES('rev','研判员','1','x','云南/昆明',?);INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES('rev2','研判员2','2','x','云南/昆明',?);INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c','活动',?,?, 'open',?);INSERT INTO suggestions(id,campaign_id,author_id,title,body,scope,status,version,created_at,updated_at) VALUES('s','c','rev','标题','正文','云南/昆明','submitted',0,?,?);INSERT INTO reviews(id,suggestion_id,reviewer_id,decision,note,version,created_at) VALUES('rv','s','rev','claimed','',0,?)`, now, now, "2026-08-18T00:00:00Z", "2026-08-30T00:00:00Z", now, now, now, now); e != nil {
		t.Fatal(e)
	}
	return db
}

func TestReviewWorkflowCompleteRejectsDuplicate(t *testing.T) {
	db := reviewWorkflowDB(t)
	defer db.Close()
	now := func() time.Time { return time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC) }
	w := ReviewWorkflow{DB: db, Reviews: repository.ReviewStore{DB: db}, Now: now}

	if e := w.Complete(context.Background(), "s", "rev", domain.DecisionAccept, "通过"); e != nil {
		t.Fatalf("first complete: %v", e)
	}
	if e := w.Complete(context.Background(), "s", "rev2", domain.DecisionAccept, "通过"); !errors.Is(e, domain.ErrConflict) {
		t.Fatalf("second complete want conflict, got %v", e)
	}
	var events int
	if e := db.QueryRow(`SELECT COUNT(*) FROM suggestion_events WHERE suggestion_id='s'`).Scan(&events); e != nil || events != 1 {
		t.Fatalf("events=%d err=%v", events, e)
	}
	var status string
	if e := db.QueryRow(`SELECT status FROM suggestions WHERE id='s'`).Scan(&status); e != nil || status != string(domain.StatusTriaged) {
		t.Fatalf("status=%s err=%v", status, e)
	}
}
