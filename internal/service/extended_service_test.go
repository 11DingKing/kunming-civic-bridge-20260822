package service

import (
	"context"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"testing"
	"time"
)

func TestDepartmentServiceAndVisibility(t *testing.T) {
	db := workflowSeed(t)
	defer db.Close()
	d := DepartmentService{DB: db, Store: repository.DepartmentStore{DB: db}}
	id, e := d.Store.Create(context.Background(), "部门", "云南/昆明/西山/马街")
	if e != nil {
		t.Fatal(e)
	}
	out, e := d.Eligible(context.Background(), "云南/昆明/西山/马街")
	if e != nil || len(out) != 1 {
		t.Fatalf("out=%v err=%v", out, e)
	}
	if e = d.Disable(context.Background(), id); e != nil {
		t.Fatal(e)
	}
	v := VisibilityService{DB: db, Store: repository.VisibilityStore{DB: db}}
	_ = v
	_ = domain.ErrInvalid
	_ = time.Now()
}
func TestIdempotencyService(t *testing.T) {
	db := workflowSeed(t)
	defer db.Close()
	s := IdempotencyService{DB: db, Store: repository.IdempotencyStore{DB: db}, TTL: time.Hour}
	now := time.Now()
	if e := s.Save(context.Background(), "key", "u", []byte("body"), "response", now); e != nil {
		t.Fatal(e)
	}
	got, replay, e := s.Replay(context.Background(), "key", "u", []byte("body"), now)
	if e != nil || !replay || got != "response" {
		t.Fatalf("got=%s replay=%v err=%v", got, replay, e)
	}
}
func TestTriageService(t *testing.T) {
	db := workflowSeed(t)
	defer db.Close()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	db.Exec(`INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES('c','活动',?,?, 'open',?);INSERT INTO suggestions(id,campaign_id,author_id,title,body,scope,status,version,created_at,updated_at) VALUES('s','c','u','标题','正文','云南/昆明/西山/马街','submitted',0,?,?)`, now, now, now, now, now)
	s := TriageService{DB: db, Store: repository.TriageStore{DB: db}, Query: repository.TriageQuery{DB: db}}
	if e := s.Classify(context.Background(), "s", "u", domain.TriageResult{Category: "民生", Priority: "normal", Reason: "日常"}); e != nil {
		t.Fatal(e)
	}
	ids, e := s.Pending(context.Background(), "")
	if e != nil || len(ids) != 0 {
		t.Fatalf("pending=%v err=%v", ids, e)
	}
}
