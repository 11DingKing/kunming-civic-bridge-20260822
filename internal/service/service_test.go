package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/audit"
	"github.com/11DingKing/kunming-civic-bridge/internal/clock"
	"github.com/11DingKing/kunming-civic-bridge/internal/notify"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"testing"
	"time"
)

func testService(t *testing.T) (Service, *sql.DB) {
	t.Helper()
	db, e := sql.Open("sqlite", ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(1)
	if e = platform.Migrate(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	_, e = db.Exec(`INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES('u1','市民','13800000000','x','x',?); INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c1','专题征集',?,?, 'open',?); INSERT INTO intake_points(id,name,scope,active,created_at) VALUES('p1','碧鸡广场','x',1,?)`, now.Format(time.RFC3339Nano), now.Add(-time.Hour).Format(time.RFC3339Nano), now.Add(time.Hour).Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if e != nil {
		t.Fatal(e)
	}
	return Service{DB: db, Store: repository.Store{DB: db}, Audit: audit.Logger{DB: db}, Notify: notify.Queue{DB: db}, Clock: clock.Fixed{T: now}}, db
}
func TestSubmitAndSend(t *testing.T) {
	s, db := testService(t)
	defer db.Close()
	x, e := s.Submit(context.Background(), "u1", "c1", "p1", "雨棚改造", "请改善社区雨棚", "x", "r1")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Send(context.Background(), x.ID, "u1", "r2"); e != nil {
		t.Fatal(e)
	}
}
