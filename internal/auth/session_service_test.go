package auth

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"testing"
	"time"
)

func TestSessionLifecycle(t *testing.T) {
	db, e := sql.Open("sqlite", ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(1)
	if e = platform.Migrate(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	if _, e = db.Exec(`INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES('u','用户','1','x','x',?)`, now); e != nil {
		t.Fatal(e)
	}
	s := SessionService{DB: db, TTL: time.Hour, Now: func() time.Time { return now }}
	id, e := s.Issue(context.Background(), "u")
	if e != nil {
		t.Fatal(e)
	}
	got, e := s.Resolve(context.Background(), id)
	if e != nil || got != "u" {
		t.Fatalf("resolve %s %v", got, e)
	}
	if e = s.Revoke(context.Background(), id); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Resolve(context.Background(), id); e == nil {
		t.Fatal("revoked session accepted")
	}
	_ = domain.ErrExpired
}
