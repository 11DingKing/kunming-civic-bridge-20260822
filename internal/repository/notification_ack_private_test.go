package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
)

func TestNotificationAckIsSingleUse(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if err = platform.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)
	if _, err = db.Exec(`INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES('u','群众','1','x','x',?); INSERT INTO notifications(id,user_id,kind,payload,status,attempts,next_attempt_at,created_at) VALUES('n','u','review','payload','pending',0,?,?)`, now, now, now); err != nil {
		t.Fatal(err)
	}
	store := NotificationQueryStore{DB: db}
	if err = store.MarkSent(context.Background(), "n"); err != nil {
		t.Fatal(err)
	}
	if err = store.MarkSent(context.Background(), "n"); err == nil {
		t.Fatal("sent notification was acknowledged twice")
	}
	var status string
	if err = db.QueryRow(`SELECT status FROM notifications WHERE id=?`, "n").Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != "sent" {
		t.Fatalf("status=%s", status)
	}
}
