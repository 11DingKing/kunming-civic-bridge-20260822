package worker

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"testing"
	"time"
)

func workerDB(t *testing.T) *sql.DB {
	db, e := sql.Open("sqlite", ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(1)
	if e = platform.Migrate(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	return db
}
func TestRecoveryAndLeases(t *testing.T) {
	db := workerDB(t)
	defer db.Close()
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	_, e := db.Exec(`INSERT INTO jobs(id,kind,payload,status,attempts,run_after,lease_until,created_at,updated_at) VALUES('j','notify','x','running',1,?,?,?,?);INSERT INTO assignments(id,suggestion_id,department_id,status,version,lease_until,assigned_at) VALUES('a','s','d','claimed',0,?,?)`, now, now.Add(-time.Minute), now, now, now.Add(-time.Minute), now)
	if e != nil {
		t.Fatal(e)
	}
	r := RecoveryRunner{DB: db}
	if e = r.Recover(context.Background(), now); e != nil {
		t.Fatal(e)
	}
	l := LeaseRunner{DB: db}
	if e = l.Release(context.Background(), now); e != nil {
		t.Fatal(e)
	}
	var jobs, assignments int
	if e = db.QueryRow(`SELECT COUNT(*) FROM jobs WHERE status='pending'`).Scan(&jobs); e != nil {
		t.Fatal(e)
	}
	if e = db.QueryRow(`SELECT COUNT(*) FROM assignments WHERE status='assigned'`).Scan(&assignments); e != nil {
		t.Fatal(e)
	}
	if jobs != 1 || assignments != 1 {
		t.Fatalf("jobs=%d assignments=%d", jobs, assignments)
	}
}
