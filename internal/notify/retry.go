package notify

import (
	"context"
	"fmt"
	"time"
)

// RetryPolicy owns notification retry schedule.
type RetryPolicy struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x RetryPolicy) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("RetryPolicy: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("RetryPolicy: timeout cannot be negative")
	}
	return nil
}
func (x RetryPolicy) Start(ctx context.Context) error {
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
func (x RetryPolicy) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x RetryPolicy) Ready() bool { return x.Enabled && x.Name != "" }
func (x RetryPolicy) Copy() RetryPolicy {
	return RetryPolicy{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x RetryPolicy) WithName(v string) RetryPolicy           { y := x.Copy(); y.Name = v; return y }
func (x RetryPolicy) WithEnabled(v bool) RetryPolicy          { y := x.Copy(); y.Enabled = v; return y }
func (x RetryPolicy) WithTimeout(v time.Duration) RetryPolicy { y := x.Copy(); y.Timeout = v; return y }
func (x RetryPolicy) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x RetryPolicy) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
