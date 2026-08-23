package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
)

func TestRunningJobCannotBeClaimedAgain(t *testing.T) {
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
	if _, err = db.Exec(`INSERT INTO jobs(id,kind,payload,status,attempts,run_after,created_at,updated_at) VALUES('j','notify','payload','pending',0,?,?,?)`, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}
	store := JobStore{DB: db}
	if _, _, _, err = store.Claim(context.Background(), now, time.Minute); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err = store.Claim(context.Background(), now, time.Minute); err == nil {
		t.Fatal("running job was claimed a second time")
	}
}
