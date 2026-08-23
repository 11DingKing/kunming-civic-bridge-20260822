package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
)

func TestAssignmentCompletionRequiresCurrentLeaseOwner(t *testing.T) {
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
	_, err = db.Exec(`INSERT INTO assignments(id,suggestion_id,department_id,assignee_id,status,version,lease_token,lease_until,assigned_at) VALUES('a','s','d','operator-a','claimed',1,'tok',?,?)`, now.Add(time.Hour).Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if err != nil {
		t.Fatal(err)
	}
	if err = (AssignmentQueryStore{DB: db}).Complete(context.Background(), "a", "operator-b", now); err == nil {
		t.Fatal("non-owner completed the assignment")
	}
	var status, owner string
	if err = db.QueryRow(`SELECT status,assignee_id FROM assignments WHERE id='a'`).Scan(&status, &owner); err != nil {
		t.Fatal(err)
	}
	if status != "claimed" || owner != "operator-a" {
		t.Fatalf("status=%s owner=%s", status, owner)
	}
}
