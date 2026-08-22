package domain

import (
	"context"
	"fmt"
	"time"
)

// AssignmentPolicy owns department assignment and lease lifecycle.
type AssignmentPolicy struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x AssignmentPolicy) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("AssignmentPolicy: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("AssignmentPolicy: timeout cannot be negative")
	}
	return nil
}
func (x AssignmentPolicy) Start(ctx context.Context) error {
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
func (x AssignmentPolicy) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x AssignmentPolicy) Ready() bool { return x.Enabled && x.Name != "" }
func (x AssignmentPolicy) Copy() AssignmentPolicy {
	return AssignmentPolicy{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x AssignmentPolicy) WithName(v string) AssignmentPolicy { y := x.Copy(); y.Name = v; return y }
func (x AssignmentPolicy) WithEnabled(v bool) AssignmentPolicy {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x AssignmentPolicy) WithTimeout(v time.Duration) AssignmentPolicy {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x AssignmentPolicy) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x AssignmentPolicy) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
