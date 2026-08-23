package repository

import (
	"context"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"testing"
	"time"
)

func TestCampaignAndFeedbackStats(t *testing.T) {
	db := storeDB(t)
	defer db.Close()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	db.Exec(`INSERT INTO suggestions(id,campaign_id,author_id,title,body,scope,status,version,created_at,updated_at) VALUES('s2','c','u','标题2','内容2','云南/昆明/西山/马街','closed',0,?,?);INSERT INTO feedback(id,suggestion_id,author_id,rating,comment,created_at) VALUES('f','s2','u',5,'满意',?)`, now, now, now)
	stats, e := (CampaignStatsStore{DB: db}).Get(context.Background(), "c")
	if e != nil || stats.Total < 1 || stats.Closed < 1 {
		t.Fatalf("stats=%+v err=%v", stats, e)
	}
	report, e := (FeedbackReportStore{DB: db}).Get(context.Background(), "s2")
	if e != nil || report.Average != 5 {
		t.Fatalf("report=%+v err=%v", report, e)
	}
	_ = domain.ErrInvalid
}
func TestSessionAndPointStats(t *testing.T) {
	db := storeDB(t)
	defer db.Close()
	now := time.Now()
	db.Exec(`INSERT INTO sessions(id,user_id,expires_at,created_at) VALUES('sess','u',?,?)`, now.Add(time.Hour).Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	items, e := (SessionStore{DB: db}).ActiveForUser(context.Background(), "u", now)
	if e != nil || len(items) != 1 {
		t.Fatalf("sessions=%v err=%v", items, e)
	}
	point, e := (PointStatsStore{DB: db}).Get(context.Background(), "")
	if e != nil {
		t.Fatal(e)
	}
	if point.Received != 0 {
		t.Fatal("unexpected point records")
	}
}
func TestApprovalAndAssignmentReports(t *testing.T) {
	db := storeDB(t)
	defer db.Close()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	db.Exec(`INSERT INTO responses(id,suggestion_id,author_id,body,status,created_at) VALUES('rr','s','u','回应内容足够长','pending_review',?);INSERT INTO assignments(id,suggestion_id,department_id,status,version,assigned_at) VALUES('aa','s','d','claimed',0,?)`, now, now)
	pending, e := (ApprovalStore{DB: db}).Pending(context.Background(), 10)
	if e != nil || len(pending) != 1 {
		t.Fatalf("pending=%v err=%v", pending, e)
	}
	report, e := (AssignmentReportStore{DB: db}).Get(context.Background(), now)
	if e != nil || report.Claimed != 1 {
		t.Fatalf("assignment report=%+v err=%v", report, e)
	}
}
