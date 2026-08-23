package domain

import (
	"fmt"
	"time"
)

type NotificationKind string

const (
	NotificationAssignment NotificationKind = "assignment"
	NotificationReview     NotificationKind = "review"
	NotificationResponse   NotificationKind = "response"
	NotificationReminder   NotificationKind = "reminder"
)

type NotificationPolicy struct {
	Kind        NotificationKind
	MaxAttempts int
	RetryAfter  time.Duration
}

func (n NotificationPolicy) Validate() error {
	if n.Kind == "" || n.MaxAttempts < 1 || n.RetryAfter <= 0 {
		return fmt.Errorf("%w: notification policy", ErrInvalid)
	}
	return nil
}
func (n NotificationPolicy) Terminal(attempt int) bool { return attempt >= n.MaxAttempts }
func (n NotificationPolicy) Next(attempt int, now time.Time) time.Time {
	if attempt < 1 {
		attempt = 1
	}
	return now.Add(n.RetryAfter * time.Duration(attempt))
}
