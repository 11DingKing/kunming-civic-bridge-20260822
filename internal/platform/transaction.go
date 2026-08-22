package platform

import (
	"context"
	"fmt"
	"time"
)

// TransactionRunner owns transaction helper.
type TransactionRunner struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x TransactionRunner) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("TransactionRunner: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("TransactionRunner: timeout cannot be negative")
	}
	return nil
}
func (x TransactionRunner) Start(ctx context.Context) error {
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
func (x TransactionRunner) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x TransactionRunner) Ready() bool { return x.Enabled && x.Name != "" }
func (x TransactionRunner) Copy() TransactionRunner {
	return TransactionRunner{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x TransactionRunner) WithName(v string) TransactionRunner { y := x.Copy(); y.Name = v; return y }
func (x TransactionRunner) WithEnabled(v bool) TransactionRunner {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x TransactionRunner) WithTimeout(v time.Duration) TransactionRunner {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x TransactionRunner) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x TransactionRunner) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
