package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

func TestSubmitRollsBackSuggestionAndEvent(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if err = platform.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	if _, err = db.Exec(`INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES('u','群众','1','x','云南/昆明/西山/马街',?)`, now.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(`INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c','七夕征集',?,?, 'open',?)`, now.Add(-time.Hour).Format(time.RFC3339Nano), now.Add(time.Hour).Format(time.RFC3339Nano), now.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(`INSERT INTO intake_points(id,name,scope,active,created_at) VALUES('p','碧鸡广场','云南/昆明/西山/马街',1,?)`, now.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(`CREATE TRIGGER fail_suggestion_event BEFORE INSERT ON suggestion_events BEGIN SELECT RAISE(ABORT, 'audit event sink unavailable'); END`); err != nil {
		t.Fatal(err)
	}
	svc := IntakeService{DB: db, Campaigns: repository.CampaignQueryStore{DB: db}, Points: repository.IntakePointStore{DB: db}, Now: func() time.Time { return now }}
	x, err := svc.Submit(context.Background(), domain.SuggestionDraft{CampaignID: "c", AuthorID: "u", IntakePointID: "p", Title: "雨棚改造", Body: "请改善社区雨棚", Scope: "云南/昆明/西山/马街", Mode: domain.ModePoint}, "point")
	if err == nil {
		t.Fatal("expected event write failure")
	}
	var suggestions, events int
	if err = db.QueryRow(`SELECT COUNT(*) FROM suggestions WHERE id=?`, x.ID).Scan(&suggestions); err != nil {
		t.Fatal(err)
	}
	if err = db.QueryRow(`SELECT COUNT(*) FROM suggestion_events WHERE suggestion_id=?`, x.ID).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if suggestions != 0 || events != 0 {
		t.Fatalf("failed submission left persistent state: suggestions=%d events=%d", suggestions, events)
	}
}
