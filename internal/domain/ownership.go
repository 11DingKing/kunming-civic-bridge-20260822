package domain

import (
	"fmt"
	"time"
)

type Lease struct {
	Token   string
	OwnerID string
	Until   time.Time
}

func (l Lease) Valid(now time.Time) bool {
	return l.Token != "" && l.OwnerID != "" && now.Before(l.Until)
}
func (l Lease) Expired(now time.Time) bool { return l.Token == "" || !now.Before(l.Until) }
func (l Lease) Renew(owner, token string, now time.Time, d time.Duration) (Lease, error) {
	if owner != l.OwnerID || token != l.Token {
		return Lease{}, fmt.Errorf("%w: lease owner", ErrForbidden)
	}
	if l.Expired(now) {
		return Lease{}, ErrExpired
	}
	l.Until = now.Add(d)
	return l, nil
}
func Claimable(status AssignmentStatus, lease Lease, now time.Time) bool {
	return status == AssignmentAssigned || status == AssignmentDeclined || lease.Expired(now)
}

type Escalation struct {
	SuggestionID string
	DueAt        time.Time
	Level        int
	Reason       string
}

func (e Escalation) Due(now time.Time) bool { return !e.DueAt.IsZero() && !now.Before(e.DueAt) }
func (e Escalation) Next(now time.Time) Escalation {
	e.Level++
	e.DueAt = now.Add(time.Duration(e.Level) * 24 * time.Hour)
	return e
}
