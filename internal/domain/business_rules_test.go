package domain

import (
	"errors"
	"testing"
	"time"
)

func TestCampaignRulesLifecycle(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	r := CampaignRules{Status: CampaignDraft, StartsAt: now.Add(-time.Hour), EndsAt: now.Add(time.Hour), Timezone: "Asia/Shanghai", MaxSubmissions: 100}
	if e := r.CanOpen(now); e != nil {
		t.Fatal(e)
	}
	r.Status = CampaignOpen
	if e := r.CanAccept(now); e != nil {
		t.Fatal(e)
	}
	r.Status = CampaignClosed
	if e := r.Archive(); e != nil {
		t.Fatal(e)
	}
	r.Status = CampaignDraft
	r.EndsAt = now.Add(-time.Minute)
	if e := r.CanOpen(now); !errors.Is(e, ErrExpired) {
		t.Fatalf("expired open=%v", e)
	}
}
func TestReviewRules(t *testing.T) {
	r := ReviewRules{State: ReviewUnassigned}
	if e := r.Claim("reviewer"); e != nil {
		t.Fatal(e)
	}
	r.State = ReviewClaimed
	r.ReviewerID = "reviewer"
	if !r.CanEdit("reviewer") || r.CanEdit("other") {
		t.Fatal("review ownership")
	}
	if e := r.Complete(DecisionAccept, "通过"); e != nil {
		t.Fatal(e)
	}
	if e := r.Return(""); e == nil {
		t.Fatal("empty return note accepted")
	}
}
func TestAssignmentRules(t *testing.T) {
	now := time.Now()
	a := AssignmentRules{Status: AssignmentAssigned, DueAt: now.Add(-time.Hour)}
	if e := a.Assign("dept"); e != nil {
		t.Fatal(e)
	}
	if !a.Overdue(now) {
		t.Fatal("overdue missing")
	}
	a.Status = AssignmentClaimed
	a.AssigneeID = "op"
	a.Lease = Lease{Token: "t", OwnerID: "op", Until: now.Add(time.Hour)}
	if e := a.Complete("op", now); e != nil {
		t.Fatal(e)
	}
	if e := a.Complete("other", now); e == nil {
		t.Fatal("foreign completion")
	}
}
func TestResponseRules(t *testing.T) {
	r := ResponseRules{Status: ResponseDraft}
	if e := r.Submit("op"); e != nil {
		t.Fatal(e)
	}
	r.Status = ResponsePending
	if e := r.Approve("sup"); e != nil {
		t.Fatal(e)
	}
	r.Status = ResponseApproved
	if e := r.Publish(); e != nil {
		t.Fatal(e)
	}
	r.Status = ResponsePublished
	if e := r.Recall("纠正内容"); e != nil {
		t.Fatal(e)
	}
}
func TestFeedbackRules(t *testing.T) {
	f := FeedbackRules{Rating: 2, Comment: "需要继续跟进", CreatedAt: time.Now(), CanReopen: true}
	if e := f.Validate(); e != nil {
		t.Fatal(e)
	}
	if !f.RequiresFollowUp() || f.Public() {
		t.Fatal("feedback policy")
	}
}
func TestQueryAndAuditRules(t *testing.T) {
	q := SuggestionQuery{Limit: 500, Offset: -2}
	q = q.Normalize()
	if q.Limit != 100 || q.Offset != 0 {
		t.Fatal("query normalize")
	}
	if e := (AuditEntry{ActorID: "u", Action: "read", ObjectType: "suggestion", ObjectID: "s", Result: "ok", RequestID: "r"}).Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestBusinessCalendar(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	c := BusinessCalendar{Location: loc, WorkingStart: 9, WorkingEnd: 18}
	if e := c.Validate(); e != nil {
		t.Fatal(e)
	}
	weekday := time.Date(2026, 8, 19, 10, 0, 0, 0, loc)
	if !c.InHours(weekday) {
		t.Fatal("weekday should be open")
	}
	weekend := time.Date(2026, 8, 22, 10, 0, 0, 0, loc)
	if c.InHours(weekend) {
		t.Fatal("weekend should be closed")
	}
}
