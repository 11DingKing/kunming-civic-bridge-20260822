package domain

import (
	"context"
	"fmt"
	"time"
)

// JobPolicy owns retry, backoff, lease and terminal failure.
type JobPolicy struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x JobPolicy) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("JobPolicy: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("JobPolicy: timeout cannot be negative")
	}
	return nil
}
func (x JobPolicy) Start(ctx context.Context) error {
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
func (x JobPolicy) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x JobPolicy) Ready() bool { return x.Enabled && x.Name != "" }
func (x JobPolicy) Copy() JobPolicy {
	return JobPolicy{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x JobPolicy) WithName(v string) JobPolicy           { y := x.Copy(); y.Name = v; return y }
func (x JobPolicy) WithEnabled(v bool) JobPolicy          { y := x.Copy(); y.Enabled = v; return y }
func (x JobPolicy) WithTimeout(v time.Duration) JobPolicy { y := x.Copy(); y.Timeout = v; return y }
func (x JobPolicy) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x JobPolicy) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
