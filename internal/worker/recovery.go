package worker

import (
	"context"
	"fmt"
	"time"
)

// RecoveryWorker owns restart recovery.
type RecoveryWorker struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x RecoveryWorker) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("RecoveryWorker: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("RecoveryWorker: timeout cannot be negative")
	}
	return nil
}
func (x RecoveryWorker) Start(ctx context.Context) error {
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
func (x RecoveryWorker) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x RecoveryWorker) Ready() bool { return x.Enabled && x.Name != "" }
func (x RecoveryWorker) Copy() RecoveryWorker {
	return RecoveryWorker{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x RecoveryWorker) WithName(v string) RecoveryWorker  { y := x.Copy(); y.Name = v; return y }
func (x RecoveryWorker) WithEnabled(v bool) RecoveryWorker { y := x.Copy(); y.Enabled = v; return y }
func (x RecoveryWorker) WithTimeout(v time.Duration) RecoveryWorker {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x RecoveryWorker) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x RecoveryWorker) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
