package domain

import (
	"context"
	"fmt"
	"time"
)

// AuditPolicy owns audit payload and immutable event validation.
type AuditPolicy struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x AuditPolicy) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("AuditPolicy: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("AuditPolicy: timeout cannot be negative")
	}
	return nil
}
func (x AuditPolicy) Start(ctx context.Context) error {
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
func (x AuditPolicy) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x AuditPolicy) Ready() bool { return x.Enabled && x.Name != "" }
func (x AuditPolicy) Copy() AuditPolicy {
	return AuditPolicy{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x AuditPolicy) WithName(v string) AuditPolicy           { y := x.Copy(); y.Name = v; return y }
func (x AuditPolicy) WithEnabled(v bool) AuditPolicy          { y := x.Copy(); y.Enabled = v; return y }
func (x AuditPolicy) WithTimeout(v time.Duration) AuditPolicy { y := x.Copy(); y.Timeout = v; return y }
func (x AuditPolicy) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x AuditPolicy) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
