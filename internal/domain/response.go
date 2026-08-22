package domain

import (
	"context"
	"fmt"
	"time"
)

// ResponsePolicy owns draft, approval, publication, recall.
type ResponsePolicy struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x ResponsePolicy) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("ResponsePolicy: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("ResponsePolicy: timeout cannot be negative")
	}
	return nil
}
func (x ResponsePolicy) Start(ctx context.Context) error {
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
func (x ResponsePolicy) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x ResponsePolicy) Ready() bool { return x.Enabled && x.Name != "" }
func (x ResponsePolicy) Copy() ResponsePolicy {
	return ResponsePolicy{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x ResponsePolicy) WithName(v string) ResponsePolicy  { y := x.Copy(); y.Name = v; return y }
func (x ResponsePolicy) WithEnabled(v bool) ResponsePolicy { y := x.Copy(); y.Enabled = v; return y }
func (x ResponsePolicy) WithTimeout(v time.Duration) ResponsePolicy {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x ResponsePolicy) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x ResponsePolicy) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
