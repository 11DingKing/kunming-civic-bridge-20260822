package repository

import (
	"context"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"testing"
	"time"
)

func TestMigrationStorePending(t *testing.T) {
	s := MigrationStore{}
	out := s.Pending([]string{"002_x.sql", "001_x.sql", "003_x.sql"}, []int{1})
	if len(out) != 2 || out[0] != "002_x.sql" {
		t.Fatalf("pending=%v", out)
	}
}
func TestIdempotencyStoreRoundTrip(t *testing.T) {
	db := storeDB(t)
	defer db.Close()
	s := IdempotencyStore{DB: db}
	r := domain.IdempotencyRecord{Key: "k", UserID: "u", RequestHash: "hash", ResponseJSON: "{}", ExpiresAt: time.Now().Add(time.Hour)}
	if err := s.Put(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	got, err := s.Find(context.Background(), "k")
	if err != nil || got.RequestHash != "hash" {
		t.Fatalf("got=%+v err=%v", got, err)
	}
	n, err := s.Purge(context.Background(), time.Now().Add(2*time.Hour))
	if err != nil || n != 1 {
		t.Fatalf("purged=%d err=%v", n, err)
	}
}
func TestScopeStoreVisibility(t *testing.T) {
	db := storeDB(t)
	defer db.Close()
	s := ScopeStore{DB: db}
	ok, err := s.UserCanSee(context.Background(), "u", "s")
	if err != nil || !ok {
		t.Fatalf("visible=%v err=%v", ok, err)
	}
	ids, err := s.SuggestionsForUser(context.Background(), "u")
	if err != nil || len(ids) != 1 {
		t.Fatalf("ids=%v err=%v", ids, err)
	}
}
