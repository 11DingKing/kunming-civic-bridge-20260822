package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"testing"
	"time"
)

func workflowDB(t *testing.T) *sql.DB {
	db, e := sql.Open("sqlite", ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(1)
	if e = platform.Migrate(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	now := "2026-08-19T00:00:00Z"
	_, e = db.Exec(`INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES('u','u','1','x','云南/昆明/西山/马街',?);INSERT INTO departments(id,name,scope,active,created_at) VALUES('d','部门','云南/昆明/西山/马街',1,?);INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c','活动',?,?, 'open',?)`, now, now, "2026-08-18T00:00:00Z", "2026-08-30T00:00:00Z", now)
	if e != nil {
		t.Fatal(e)
	}
	return db
}
func TestWorkflowClaimAndReopen(t *testing.T) {
	db := workflowDB(t)
	defer db.Close()
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	_, e := db.Exec(`INSERT INTO suggestions(id,campaign_id,author_id,title,body,scope,status,version,created_at,updated_at) VALUES('s','c','u','title','body','云南/昆明/西山/马街','responded',2,?,?);INSERT INTO responses(id,suggestion_id,author_id,body,status,created_at) VALUES('r','s','u','long response body here','published',?)`, now, now, now)
	if e != nil {
		t.Fatal(e)
	}
	w := Workflow{DB: db, Suggestions: repository.SuggestionRepository{DB: db}, Jobs: repository.JobStore{DB: db}, Clock: func() time.Time { return now }}
	if e = w.ReopenFromFeedback(context.Background(), "s", "u"); e != nil {
		t.Fatal(e)
	}
	var status string
	if e = db.QueryRow(`SELECT status FROM suggestions WHERE id='s'`).Scan(&status); e != nil || status != "reopened" {
		t.Fatalf("status=%s err=%v", status, e)
	}
}
