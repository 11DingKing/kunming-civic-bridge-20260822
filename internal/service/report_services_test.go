package service

import (
	"context"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"testing"
)

func TestReportServices(t *testing.T) {
	db := workflowSeed(t)
	defer db.Close()
	stats := CampaignStatsService{DB: db, Store: repository.CampaignStatsStore{DB: db}}
	if _, e := stats.Get(context.Background(), "missing"); e != nil {
		t.Fatal(e)
	}
	attach := AttachmentService{DB: db, Store: repository.AttachmentStore{DB: db}}
	if _, e := attach.Add(context.Background(), "s", "照片.jpg", "image/jpeg", "key", 100); e != nil {
		t.Fatal(e)
	}
	items, e := attach.List(context.Background(), "s")
	if e != nil || len(items) != 1 {
		t.Fatalf("attachments=%v err=%v", items, e)
	}
}
