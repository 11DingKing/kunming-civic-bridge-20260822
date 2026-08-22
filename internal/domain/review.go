package domain

import (
	"context"
	"fmt"
	"time"
)

// ReviewPolicy owns review decisions and reassignment rules.
type ReviewPolicy struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x ReviewPolicy) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("ReviewPolicy: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("ReviewPolicy: timeout cannot be negative")
	}
	return nil
}
func (x ReviewPolicy) Start(ctx context.Context) error {
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
func (x ReviewPolicy) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x ReviewPolicy) Ready() bool { return x.Enabled && x.Name != "" }
func (x ReviewPolicy) Copy() ReviewPolicy {
	return ReviewPolicy{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x ReviewPolicy) WithName(v string) ReviewPolicy  { y := x.Copy(); y.Name = v; return y }
func (x ReviewPolicy) WithEnabled(v bool) ReviewPolicy { y := x.Copy(); y.Enabled = v; return y }
func (x ReviewPolicy) WithTimeout(v time.Duration) ReviewPolicy {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x ReviewPolicy) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x ReviewPolicy) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
