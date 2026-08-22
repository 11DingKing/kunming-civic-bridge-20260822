package domain

import (
	"context"
	"fmt"
	"time"
)

// IntakePolicy owns offline point authorization and scope.
type IntakePolicy struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x IntakePolicy) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("IntakePolicy: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("IntakePolicy: timeout cannot be negative")
	}
	return nil
}
func (x IntakePolicy) Start(ctx context.Context) error {
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
func (x IntakePolicy) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x IntakePolicy) Ready() bool { return x.Enabled && x.Name != "" }
func (x IntakePolicy) Copy() IntakePolicy {
	return IntakePolicy{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x IntakePolicy) WithName(v string) IntakePolicy  { y := x.Copy(); y.Name = v; return y }
func (x IntakePolicy) WithEnabled(v bool) IntakePolicy { y := x.Copy(); y.Enabled = v; return y }
func (x IntakePolicy) WithTimeout(v time.Duration) IntakePolicy {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x IntakePolicy) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x IntakePolicy) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
