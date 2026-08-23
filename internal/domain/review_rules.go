package domain

import (
	"fmt"
	"strings"
)

type ReviewState string

const (
	ReviewUnassigned ReviewState = "unassigned"
	ReviewClaimed    ReviewState = "claimed"
	ReviewCompleted  ReviewState = "completed"
	ReviewReturned   ReviewState = "returned"
)

type ReviewRules struct {
	State      ReviewState
	ReviewerID string
	Version    int
	Note       string
}

func (r ReviewRules) Claim(actor string) error {
	if actor == "" {
		return fmt.Errorf("%w: reviewer", ErrInvalid)
	}
	if r.State != ReviewUnassigned && r.State != ReviewReturned {
		return fmt.Errorf("%w: review claim", ErrConflict)
	}
	return nil
}
func (r ReviewRules) Complete(decision ReviewDecision, note string) error {
	if r.State != ReviewClaimed {
		return fmt.Errorf("%w: review completion", ErrConflict)
	}
	if e := ValidateReview(decision, note); e != nil {
		return e
	}
	return nil
}
func (r ReviewRules) Return(note string) error {
	if r.State != ReviewClaimed {
		return fmt.Errorf("%w: review return", ErrConflict)
	}
	if strings.TrimSpace(note) == "" {
		return fmt.Errorf("%w: return note", ErrInvalid)
	}
	return nil
}
func (r ReviewRules) CanEdit(actor string) bool {
	return actor != "" && r.State == ReviewClaimed && actor == r.ReviewerID
}
