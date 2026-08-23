package domain

import (
	"testing"
	"time"
)

func TestScopeContainment(t *testing.T) {
	parent, e := ParseScope("云南/昆明/西山/")
	if e == nil {
		t.Fatal("empty street should be rejected")
	}
	_ = parent
	s := Scope{Province: "云南", City: "昆明", District: "西山"}
	child := Scope{Province: "云南", City: "昆明", District: "西山", Street: "马街"}
	if !s.Contains(child) || child.Contains(s) {
		t.Fatal("scope hierarchy broken")
	}
}
func TestLeaseRenewal(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	l := Lease{Token: "t", OwnerID: "u", Until: now.Add(time.Hour)}
	next, e := l.Renew("u", "t", now, 2*time.Hour)
	if e != nil || !next.Valid(now.Add(time.Hour)) {
		t.Fatalf("renewal=%+v err=%v", next, e)
	}
	if _, e = l.Renew("other", "t", now, time.Hour); e == nil {
		t.Fatal("foreign owner renewed lease")
	}
}
func TestIdempotencyHashAndExpiry(t *testing.T) {
	now := time.Now()
	r := IdempotencyRecord{UserID: "u", RequestHash: HashRequest([]byte("a")), ExpiresAt: now.Add(time.Hour)}
	if e := r.Usable("u", HashRequest([]byte("a")), now); e != nil {
		t.Fatal(e)
	}
	if e := r.Usable("u", HashRequest([]byte("b")), now); e == nil {
		t.Fatal("changed request accepted")
	}
	if e := r.Usable("u", r.RequestHash, now.Add(2*time.Hour)); e == nil {
		t.Fatal("expired key accepted")
	}
}
func TestJobBackoffAndTerminalState(t *testing.T) {
	if Backoff(1) != time.Second || Backoff(4) != 8*time.Second {
		t.Fatal("backoff")
	}
	if NextJobState(JobRunning, nil, 1, 3) != JobDone {
		t.Fatal("success state")
	}
	if NextJobState(JobRunning, ErrInvalid, 3, 3) != JobDead {
		t.Fatal("dead state")
	}
}
