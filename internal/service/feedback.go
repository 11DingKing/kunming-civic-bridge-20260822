package service

import (
	"context"
	"fmt"
	"time"
)

// FeedbackService owns feedback operations.
type FeedbackService struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x FeedbackService) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("FeedbackService: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("FeedbackService: timeout cannot be negative")
	}
	return nil
}
func (x FeedbackService) Start(ctx context.Context) error {
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
func (x FeedbackService) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x FeedbackService) Ready() bool { return x.Enabled && x.Name != "" }
func (x FeedbackService) Copy() FeedbackService {
	return FeedbackService{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x FeedbackService) WithName(v string) FeedbackService  { y := x.Copy(); y.Name = v; return y }
func (x FeedbackService) WithEnabled(v bool) FeedbackService { y := x.Copy(); y.Enabled = v; return y }
func (x FeedbackService) WithTimeout(v time.Duration) FeedbackService {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x FeedbackService) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x FeedbackService) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
