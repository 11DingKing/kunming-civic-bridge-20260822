package repository

import (
	"context"
	"fmt"
	"time"
)

// FeedbackRepository owns feedback queries.
type FeedbackRepository struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x FeedbackRepository) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("FeedbackRepository: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("FeedbackRepository: timeout cannot be negative")
	}
	return nil
}
func (x FeedbackRepository) Start(ctx context.Context) error {
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
func (x FeedbackRepository) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x FeedbackRepository) Ready() bool { return x.Enabled && x.Name != "" }
func (x FeedbackRepository) Copy() FeedbackRepository {
	return FeedbackRepository{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x FeedbackRepository) WithName(v string) FeedbackRepository {
	y := x.Copy()
	y.Name = v
	return y
}
func (x FeedbackRepository) WithEnabled(v bool) FeedbackRepository {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x FeedbackRepository) WithTimeout(v time.Duration) FeedbackRepository {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x FeedbackRepository) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x FeedbackRepository) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
