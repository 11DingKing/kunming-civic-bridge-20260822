package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"testing"
	"time"
)

func storeDB(t *testing.T) *sql.DB {
	db, e := sql.Open("sqlite", ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(1)
	if e = platform.Migrate(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	now := "2026-08-19T00:00:00Z"
	if _, e = db.Exec(`INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES('u','u','1','x','云南/昆明/西山/马街',?)`, now); e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec(`INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c','活动',?,?, 'draft',?)`, "2026-08-18T00:00:00Z", "2026-08-30T00:00:00Z", now); e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec(`INSERT INTO departments(id,name,scope,active,created_at) VALUES('d','部门','云南/昆明/西山/马街',1,?)`, now); e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec(`INSERT INTO suggestions(id,campaign_id,author_id,title,body,scope,status,version,created_at,updated_at) VALUES('s','c','u','雨棚','社区雨棚改造建议','云南/昆明/西山/马街','submitted',0,?,?)`, now, now); e != nil {
		t.Fatal(e)
	}
	return db
}
func TestCampaignQueryStore(t *testing.T) {
	db := storeDB(t)
	defer db.Close()
	s := CampaignQueryStore{DB: db}
	r, e := s.Get(context.Background(), "c")
	if e != nil || r.Name != "活动" {
		t.Fatalf("campaign=%+v err=%v", r, e)
	}
	if e = s.Open(context.Background(), "c", time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)); e != nil {
		t.Fatal(e)
	}
	open, e := s.ListOpen(context.Background(), time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC), 10)
	if e != nil || len(open) != 1 {
		t.Fatalf("open=%d err=%v", len(open), e)
	}
}
func TestSuggestionQueryStore(t *testing.T) {
	db := storeDB(t)
	defer db.Close()
	s := SuggestionQueryStore{DB: db}
	out, e := s.Search(context.Background(), SuggestionFilter{Scope: "云南/昆明/西山/马街", Keyword: "雨棚", Limit: 10})
	if e != nil || len(out) != 1 {
		t.Fatalf("search=%d err=%v", len(out), e)
	}
	if _, e = s.Search(context.Background(), SuggestionFilter{Status: string(domain.StatusClosed)}); e == nil {
		t.Fatal("closed filter should be forbidden")
	}
}
func TestAssignmentQueryStore(t *testing.T) {
	db := storeDB(t)
	defer db.Close()
	s := AssignmentQueryStore{DB: db}
	id, e := s.Create(context.Background(), "s", "d", nil)
	if e != nil {
		t.Fatal(e)
	}
	a, e := s.Get(context.Background(), id)
	if e != nil {
		t.Fatal(e)
	}
	until := time.Now().Add(time.Hour)
	if e = s.Claim(context.Background(), id, "op", "token", until, a.Version); e != nil {
		t.Fatal(e)
	}
	if e = s.Complete(context.Background(), id, "op", time.Now()); e != nil {
		t.Fatal(e)
	}
	a, e = s.Get(context.Background(), id)
	if e != nil || a.Status != "done" {
		t.Fatalf("assignment=%+v err=%v", a, e)
	}
}
func TestNotificationAndAuditStores(t *testing.T) {
	db := storeDB(t)
	defer db.Close()
	n := NotificationQueryStore{DB: db}
	now := time.Now()
	_, e := db.Exec(`INSERT INTO notifications(id,user_id,suggestion_id,kind,payload,status,attempts,next_attempt_at,created_at) VALUES('n','u','s','review','{}','pending',0,?,?)`, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if e != nil {
		t.Fatal(e)
	}
	ready, e := n.Ready(context.Background(), now.Add(time.Second), 5)
	if e != nil || len(ready) != 1 {
		t.Fatalf("ready=%d err=%v", len(ready), e)
	}
	if e = n.MarkSent(context.Background(), "n"); e != nil {
		t.Fatal(e)
	}
	a := AuditQueryStore{DB: db}
	if e = a.Append(context.Background(), domain.AuditEntry{ActorID: "u", Action: "read", ObjectType: "suggestion", ObjectID: "s", Result: "ok", RequestID: "r"}); e != nil {
		t.Fatal(e)
	}
	events, e := a.ForObject(context.Background(), "suggestion", "s")
	if e != nil || len(events) != 1 {
		t.Fatalf("events=%d err=%v", len(events), e)
	}
}
