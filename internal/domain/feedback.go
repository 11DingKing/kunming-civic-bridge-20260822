package domain

import (
	"context"
	"fmt"
	"time"
)

// FeedbackPolicy owns rating, reopen and follow-up.
type FeedbackPolicy struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x FeedbackPolicy) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("FeedbackPolicy: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("FeedbackPolicy: timeout cannot be negative")
	}
	return nil
}
func (x FeedbackPolicy) Start(ctx context.Context) error {
	if err := x.Validate(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
func (x FeedbackPolicy) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x FeedbackPolicy) Ready() bool { return x.Enabled && x.Name != "" }
func (x FeedbackPolicy) Copy() FeedbackPolicy {
	return FeedbackPolicy{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x FeedbackPolicy) WithName(v string) FeedbackPolicy  { y := x.Copy(); y.Name = v; return y }
func (x FeedbackPolicy) WithEnabled(v bool) FeedbackPolicy { y := x.Copy(); y.Enabled = v; return y }
func (x FeedbackPolicy) WithTimeout(v time.Duration) FeedbackPolicy {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x FeedbackPolicy) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x FeedbackPolicy) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
