package domain

import (
	"testing"
	"time"
)

func TestTransitionRules(t *testing.T) {
	if !StatusSubmitted.CanTransition(StatusTriaged) {
		t.Fatal("submitted should be triaged")
	}
	if StatusDraft.CanTransition(StatusClosed) {
		t.Fatal("draft cannot close")
	}
}
func TestCampaignOpenWindow(t *testing.T) {
	n := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	c := Campaign{Status: "open", StartsAt: n, EndsAt: n.Add(time.Hour)}
	if c.OpenAt(n.Add(30*time.Minute)) != nil {
		t.Fatal("window should be open")
	}
	if c.OpenAt(n.Add(2*time.Hour)) == nil {
		t.Fatal("window should close")
	}
}
