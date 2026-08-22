package domain

import (
	"context"
	"fmt"
	"time"
)

// IdempotencyPolicy owns request hash and replay expiry.
type IdempotencyPolicy struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x IdempotencyPolicy) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("IdempotencyPolicy: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("IdempotencyPolicy: timeout cannot be negative")
	}
	return nil
}
func (x IdempotencyPolicy) Start(ctx context.Context) error {
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
func (x IdempotencyPolicy) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x IdempotencyPolicy) Ready() bool { return x.Enabled && x.Name != "" }
func (x IdempotencyPolicy) Copy() IdempotencyPolicy {
	return IdempotencyPolicy{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x IdempotencyPolicy) WithName(v string) IdempotencyPolicy { y := x.Copy(); y.Name = v; return y }
func (x IdempotencyPolicy) WithEnabled(v bool) IdempotencyPolicy {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x IdempotencyPolicy) WithTimeout(v time.Duration) IdempotencyPolicy {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x IdempotencyPolicy) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x IdempotencyPolicy) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
