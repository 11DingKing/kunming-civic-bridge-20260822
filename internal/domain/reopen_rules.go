package domain

import (
	"fmt"
	"time"
)

type ReopenRequest struct {
	SuggestionID, AuthorID, Reason string
	RequestedAt                    time.Time
}

func (r ReopenRequest) Validate() error {
	if r.SuggestionID == "" || r.AuthorID == "" || r.Reason == "" {
		return fmt.Errorf("%w: reopen request", ErrInvalid)
	}
	if r.RequestedAt.IsZero() {
		return fmt.Errorf("%w: reopen time", ErrInvalid)
	}
	return nil
}
func (r ReopenRequest) WithinWindow(closedAt time.Time, now time.Time) bool {
	return !closedAt.IsZero() && !now.Before(closedAt) && now.Sub(closedAt) <= 30*24*time.Hour
}
func (r ReopenRequest) CanApply(status SuggestionStatus) bool {
	return status == StatusResponded || status == StatusClosed
}
