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

func TestReopenApplyIsIdempotent(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if err = platform.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	if _, err = db.Exec(`INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES('u','群众','1','x','x',?); INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c','活动',?,?, 'open',?); INSERT INTO suggestions(id,campaign_id,author_id,title,body,scope,status,version,created_at,updated_at) VALUES('s','c','u','标题','正文','x','responded',2,?,?); INSERT INTO feedback(id,suggestion_id,author_id,rating,comment,created_at) VALUES('f','s','u',1,'仍需改进',?)`, now.Format(time.RFC3339Nano), now.Add(-time.Hour).Format(time.RFC3339Nano), now.Add(time.Hour).Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}
	svc := ReopenService{DB: db, Store: repository.ReopenStore{DB: db}, Now: func() time.Time { return now }}
	if err = svc.Apply(context.Background(), "s"); err != nil {
		t.Fatal(err)
	}
	if err = svc.Apply(context.Background(), "s"); err == nil {
		t.Fatal("reopen was applied twice")
	}
	var status string
	var version int
	if err = db.QueryRow(`SELECT status,version FROM suggestions WHERE id='s'`).Scan(&status, &version); err != nil {
		t.Fatal(err)
	}
	if status != string(domain.StatusReopened) || version != 3 {
		t.Fatalf("status=%s version=%d", status, version)
	}
}
