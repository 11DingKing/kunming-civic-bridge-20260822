package service

import (
	"context"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"testing"
)

func TestPointReportService(t *testing.T) {
	db := workflowSeed(t)
	defer db.Close()
	s := PointReportService{DB: db, Store: repository.PointReportStore{DB: db}}
	if _, e := s.Reports(context.Background()); e != nil {
		t.Fatal(e)
	}
	if e := s.RequireActive(context.Background(), "missing"); e == nil {
		t.Fatal("missing point accepted")
	}
}
