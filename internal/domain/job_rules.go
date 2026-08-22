package domain

import (
	"fmt"
	"time"
)

type JobState string

const (
	JobPending JobState = "pending"
	JobRunning JobState = "running"
	JobDone    JobState = "done"
	JobFailed  JobState = "failed"
	JobDead    JobState = "dead"
)

func (s JobState) CanRetry() bool { return s == JobPending || s == JobFailed }
func Backoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if attempt > 8 {
		attempt = 8
	}
	return time.Duration(1<<uint(attempt-1)) * time.Second
}
func NextJobState(current JobState, err error, attempt, max int) JobState {
	if err == nil {
		return JobDone
	}
	if attempt >= max {
		return JobDead
	}
	if !current.CanRetry() {
		return JobFailed
	}
	return JobPending
}
func ValidateJob(kind, payload string) error {
	if kind == "" || payload == "" {
		return fmt.Errorf("%w: job", ErrInvalid)
	}
	return nil
}
