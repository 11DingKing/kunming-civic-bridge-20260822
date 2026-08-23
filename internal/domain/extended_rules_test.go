package domain

import (
	"testing"
	"time"
)

func TestAssignmentEvents(t *testing.T) {
	e := AssignmentEvent{AssignmentID: "a", ActorID: "u", From: "assigned", To: "claimed", At: time.Now()}
	if err := e.Validate(); err != nil || !e.LeaseRequired() || e.Terminal() {
		t.Fatalf("event=%+v err=%v", e, err)
	}
	e.To = "done"
	if !e.Terminal() {
		t.Fatal("done should be terminal")
	}
}
func TestNotificationPolicy(t *testing.T) {
	p := NotificationPolicy{Kind: NotificationReview, MaxAttempts: 3, RetryAfter: time.Second}
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
	if p.Terminal(2) || !p.Terminal(3) {
		t.Fatal("attempt policy")
	}
	if !p.Next(2, time.Now()).After(time.Now()) {
		t.Fatal("retry next")
	}
}
func TestTriageAndDepartmentRules(t *testing.T) {
	r := TriageResult{Category: "交通", Priority: "urgent", Reason: "影响群众", Sensitive: true}
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
	if !r.RequiresSupervisor() || r.PublicCategory() == "交通" {
		t.Fatal("sensitive policy")
	}
	d := DepartmentScope{DepartmentID: "d", Name: "部门", Scope: "云南/昆明/西山", Active: true}
	if !d.Eligible("云南/昆明/西山/马街") {
		t.Fatal("department should receive")
	}
}
func TestResponseVisibility(t *testing.T) {
	v := ResponseVisibility{Status: ResponseDraft, Visibility: VisibilityPublic}
	if err := v.Validate(); err == nil {
		t.Fatal("draft published")
	}
	v.Status = ResponsePublished
	v.Sensitive = true
	if v.Audience() != "author" {
		t.Fatal("sensitive audience")
	}
}
func TestOfflineIntakeAndExport(t *testing.T) {
	i := OfflineIntake{PointID: "p", OperatorID: "o", AuthorPhone: "138 0000", Title: "标题", Body: "正文", Scope: "云南/昆明/西山/马街", Consent: true}
	if err := i.Validate(); err != nil {
		t.Fatal(err)
	}
	if i.Normalized().AuthorPhone != "1380000" {
		t.Fatal("phone normalize")
	}
	e := AuditExport{ObjectType: "suggestion", ObjectID: "s", RequestedBy: "u", From: time.Now(), To: time.Now().Add(time.Hour), Format: "json"}
	if err := e.Validate(); err != nil {
		t.Fatal(err)
	}
}
