package domain

import (
	"fmt"
	"time"
)

type AssignmentRules struct {
	Status                   AssignmentStatus
	DepartmentID, AssigneeID string
	Lease                    Lease
	DueAt                    time.Time
	Version                  int
}

func (a AssignmentRules) Assign(department string) error {
	if department == "" {
		return fmt.Errorf("%w: department", ErrInvalid)
	}
	if a.Status != AssignmentAssigned && a.Status != AssignmentDeclined {
		return fmt.Errorf("%w: assignment status", ErrConflict)
	}
	return nil
}
func (a AssignmentRules) Claim(owner, token string, now time.Time) error {
	if !Claimable(a.Status, a.Lease, now) {
		return fmt.Errorf("%w: assignment claimed", ErrConflict)
	}
	if owner == "" || token == "" {
		return fmt.Errorf("%w: claim identity", ErrInvalid)
	}
	return nil
}
func (a AssignmentRules) Complete(owner string, now time.Time) error {
	if a.Status != AssignmentClaimed || a.AssigneeID != owner {
		return fmt.Errorf("%w: assignment owner", ErrForbidden)
	}
	if !a.Lease.Valid(now) {
		return ErrExpired
	}
	return nil
}
func (a AssignmentRules) Overdue(now time.Time) bool {
	return !a.DueAt.IsZero() && now.After(a.DueAt) && a.Status != AssignmentDone
}
