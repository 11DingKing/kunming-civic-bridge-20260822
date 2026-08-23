package domain

import (
	"testing"
	"time"
)

func TestReopenAndExportRules(t *testing.T) {
	now := time.Now()
	r := ReopenRequest{SuggestionID: "s", AuthorID: "u", Reason: "未解决", RequestedAt: now}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
	if !r.WithinWindow(now.Add(-time.Hour), now) || !r.CanApply(StatusResponded) {
		t.Fatal("reopen window")
	}
	e := AuditExport{ObjectType: "suggestion", ObjectID: "s", RequestedBy: "u", From: now.Add(-time.Hour), To: now, Format: "csv"}
	if err := e.Validate(); err != nil || !e.Includes(now.Add(-time.Minute)) {
		t.Fatalf("export=%v", err)
	}
}
func TestNotificationAndVisibilityRules(t *testing.T) {
	p := NotificationPolicy{Kind: NotificationReminder, MaxAttempts: 4, RetryAfter: time.Second}
	if e := p.Validate(); e != nil {
		t.Fatal(e)
	}
	v := ResponseVisibility{Status: ResponsePublished, Visibility: VisibilityRedacted, RedactionNote: "个人信息"}
	if e := v.Validate(); e != nil {
		t.Fatal(e)
	}
	i := OfflineIntake{PointID: "p", OperatorID: "o", AuthorPhone: "1", Title: "标题", Body: "正文", Scope: "scope", Consent: true}
	if e := i.Validate(); e != nil {
		t.Fatal(e)
	}
}
