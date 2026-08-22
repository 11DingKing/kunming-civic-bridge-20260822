package domain

import (
	"testing"
	"time"
)

func TestSLAAndPointRules(t *testing.T) {
	r := ReopenRequest{SuggestionID: "s", AuthorID: "u", Reason: "follow up", RequestedAt: time.Now()}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
	p := NotificationPolicy{Kind: NotificationAssignment, MaxAttempts: 3, RetryAfter: time.Second}
	if e := p.Validate(); e != nil {
		t.Fatal(e)
	}
	if p.Terminal(2) || !p.Terminal(3) {
		t.Fatal("terminal policy")
	}
}
func TestAuditExportWindows(t *testing.T) {
	from := time.Now().Add(-time.Hour)
	to := time.Now()
	e := AuditExport{ObjectType: "suggestion", ObjectID: "s", RequestedBy: "u", From: from, To: to, Format: "json"}
	if err := e.Validate(); err != nil {
		t.Fatal(err)
	}
	if !e.Includes(from.Add(time.Minute)) || e.Includes(to) {
		t.Fatal("window")
	}
}
