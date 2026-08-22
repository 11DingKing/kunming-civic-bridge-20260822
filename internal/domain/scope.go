package domain

import (
	"context"
	"fmt"
	"time"
)

// ScopePolicy owns county, district, street visibility.
type ScopePolicy struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x ScopePolicy) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("ScopePolicy: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("ScopePolicy: timeout cannot be negative")
	}
	return nil
}
func (x ScopePolicy) Start(ctx context.Context) error {
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
func (x ScopePolicy) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x ScopePolicy) Ready() bool { return x.Enabled && x.Name != "" }
func (x ScopePolicy) Copy() ScopePolicy {
	return ScopePolicy{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x ScopePolicy) WithName(v string) ScopePolicy           { y := x.Copy(); y.Name = v; return y }
func (x ScopePolicy) WithEnabled(v bool) ScopePolicy          { y := x.Copy(); y.Enabled = v; return y }
func (x ScopePolicy) WithTimeout(v time.Duration) ScopePolicy { y := x.Copy(); y.Timeout = v; return y }
func (x ScopePolicy) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x ScopePolicy) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
