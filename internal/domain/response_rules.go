package domain

import (
	"fmt"
	"strings"
)

type ResponseRules struct {
	Status               ResponseStatus
	AuthorID, ReviewerID string
	Version              int
}

func (r ResponseRules) Submit(author string) error {
	if r.Status != ResponseDraft && r.Status != ResponseApproved {
		return fmt.Errorf("%w: response submit", ErrConflict)
	}
	if author == "" {
		return fmt.Errorf("%w: response author", ErrInvalid)
	}
	return nil
}
func (r ResponseRules) Approve(supervisor string) error {
	if r.Status != ResponsePending {
		return fmt.Errorf("%w: response approval", ErrConflict)
	}
	if supervisor == "" {
		return fmt.Errorf("%w: supervisor", ErrInvalid)
	}
	return nil
}
func (r ResponseRules) Publish() error {
	if r.Status != ResponseApproved {
		return fmt.Errorf("%w: response publish", ErrConflict)
	}
	return nil
}
func (r ResponseRules) Recall(note string) error {
	if r.Status != ResponsePublished {
		return fmt.Errorf("%w: response recall", ErrConflict)
	}
	if strings.TrimSpace(note) == "" {
		return fmt.Errorf("%w: recall note", ErrInvalid)
	}
	return nil
}
