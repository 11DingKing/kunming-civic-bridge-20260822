package service

import (
	"context"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"testing"
	"time"
)

func TestSLAServiceDelegates(t *testing.T) {
	db := workflowSeed(t)
	defer db.Close()
	s := SLAService{DB: db, Store: repository.SLAStore{DB: db}}
	if _, e := s.Summary(context.Background(), time.Now()); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Due(context.Background(), time.Now()); e != nil {
		t.Fatal(e)
	}
}
