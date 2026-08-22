package httpapi

import (
	"context"
	"fmt"
	"time"
)

// FeedbackHandlers owns feedback endpoints.
type FeedbackHandlers struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x FeedbackHandlers) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("FeedbackHandlers: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("FeedbackHandlers: timeout cannot be negative")
	}
	return nil
}
func (x FeedbackHandlers) Start(ctx context.Context) error {
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
func (x FeedbackHandlers) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x FeedbackHandlers) Ready() bool { return x.Enabled && x.Name != "" }
func (x FeedbackHandlers) Copy() FeedbackHandlers {
	return FeedbackHandlers{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x FeedbackHandlers) WithName(v string) FeedbackHandlers { y := x.Copy(); y.Name = v; return y }
func (x FeedbackHandlers) WithEnabled(v bool) FeedbackHandlers {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x FeedbackHandlers) WithTimeout(v time.Duration) FeedbackHandlers {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x FeedbackHandlers) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x FeedbackHandlers) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
