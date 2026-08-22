package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"testing"
	"time"
)

func workflowSeed(t *testing.T) *sql.DB {
	db, e := sql.Open("sqlite", ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(1)
	if e = platform.Migrate(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	now := "2026-08-19T00:00:00Z"
	db.Exec(`INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES('u','u','1','x','云南/昆明/西山/马街',?);INSERT INTO intake_points(id,name,scope,active,created_at) VALUES('p','点位','云南/昆明/西山/马街',1,?)`, now, now)
	return db
}
func TestIntakeServiceSubmit(t *testing.T) {
	db := workflowSeed(t)
	defer db.Close()
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	db.Exec(`INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c','活动',?,?, 'open',?)`, now.Add(-time.Hour).Format(time.RFC3339Nano), now.Add(time.Hour).Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	s := IntakeService{DB: db, Campaigns: repository.CampaignQueryStore{DB: db}, Points: repository.IntakePointStore{DB: db}, Now: func() time.Time { return now }}
	x, e := s.Submit(context.Background(), domain.SuggestionDraft{CampaignID: "c", AuthorID: "u", IntakePointID: "p", Title: "建议标题", Body: "建议内容足够长", Scope: "云南/昆明/西山/马街", Mode: domain.ModePoint}, "point")
	if e != nil || x.Status != domain.StatusDraft {
		t.Fatalf("suggestion=%+v err=%v", x, e)
	}
	var events int
	db.QueryRow(`SELECT COUNT(*) FROM suggestion_events WHERE suggestion_id=?`, x.ID).Scan(&events)
	if events != 1 {
		t.Fatal("creation event missing")
	}
}
func TestQueryServiceFiltering(t *testing.T) {
	db := workflowSeed(t)
	defer db.Close()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	db.Exec(`INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c','活动',?,?, 'open',?);INSERT INTO suggestions(id,campaign_id,author_id,title,body,scope,status,version,created_at,updated_at) VALUES('s','c','u','充电停车场','闲置用地建议','云南/昆明/西山/马街','submitted',0,?,?)`, now, now, now, now, now)
	q := QueryService{DB: db}
	out, e := q.Search(context.Background(), domain.SuggestionQuery{Keyword: "停车场", Limit: 10})
	if e != nil || len(out) != 1 {
		t.Fatalf("out=%d err=%v", len(out), e)
	}
	count, e := q.Count(context.Background(), domain.SuggestionQuery{Scope: "云南/昆明/西山/马街"})
	if e != nil || count != 1 {
		t.Fatalf("count=%d err=%v", count, e)
	}
}
func TestFeedbackWorkflowReopens(t *testing.T) {
	db := workflowSeed(t)
	defer db.Close()
	now := time.Now()
	db.Exec(`INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c','活动',?,?, 'open',?);INSERT INTO suggestions(id,campaign_id,author_id,title,body,scope,status,version,created_at,updated_at) VALUES('s','c','u','标题','正文','云南/昆明/西山/马街','responded',1,?,?);INSERT INTO responses(id,suggestion_id,author_id,body,status,created_at) VALUES('r','s','u','这是足够长的回应内容','published',?)`, now.Add(-time.Hour).Format(time.RFC3339Nano), now.Add(time.Hour).Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	w := FeedbackWorkflow{DB: db, Feedback: repository.FeedbackStore{DB: db}, Now: func() time.Time { return now }}
	if e := w.Record(context.Background(), "s", "u", "仍需改进", 2); e != nil {
		t.Fatal(e)
	}
	if e := w.ReopenIfNeeded(context.Background(), "s", 2); e != nil {
		t.Fatal(e)
	}
	var status string
	db.QueryRow(`SELECT status FROM suggestions WHERE id='s'`).Scan(&status)
	if status != "reopened" {
		t.Fatalf("status=%s", status)
	}
}
