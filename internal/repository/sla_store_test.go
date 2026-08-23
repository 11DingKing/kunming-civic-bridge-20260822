package repository

import (
	"context"
	"testing"
	"time"
)

func TestSLAStoreDueAndSummary(t *testing.T) {
	db := storeDB(t)
	defer db.Close()
	now := time.Now()
	db.Exec(`INSERT INTO suggestions(id,campaign_id,author_id,title,body,scope,status,version,due_at,created_at,updated_at) VALUES('due','c','u','标题','正文','云南/昆明/西山/马街','in_progress',0,?,?,?)`, now.Add(-time.Minute).Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	s := SLAStore{DB: db}
	ids, e := s.Due(context.Background(), now)
	if e != nil || len(ids) != 1 {
		t.Fatalf("due=%v err=%v", ids, e)
	}
	summary, e := s.Summary(context.Background(), now)
	if e != nil || summary.Late != 1 {
		t.Fatalf("summary=%+v err=%v", summary, e)
	}
}
func TestPointReportStore(t *testing.T) {
	db := storeDB(t)
	defer db.Close()
	db.Exec(`INSERT INTO intake_points(id,name,scope,active,created_at) VALUES('p','点位','云南/昆明/西山/马街',1,?)`, time.Now().UTC().Format(time.RFC3339Nano))
	s := PointReportStore{DB: db}
	items, e := s.Reports(context.Background())
	if e != nil || len(items) != 1 || !items[0].Active {
		t.Fatalf("items=%v err=%v", items, e)
	}
	if e = s.RequireActive(context.Background(), "p"); e != nil {
		t.Fatal(e)
	}
}
