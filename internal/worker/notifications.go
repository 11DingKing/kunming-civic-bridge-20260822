package worker

import (
	"context"
	"fmt"
	"time"
)

// NotificationWorker owns notification worker.
type NotificationWorker struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x NotificationWorker) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("NotificationWorker: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("NotificationWorker: timeout cannot be negative")
	}
	return nil
}
func (x NotificationWorker) Start(ctx context.Context) error {
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
func (x NotificationWorker) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x NotificationWorker) Ready() bool { return x.Enabled && x.Name != "" }
func (x NotificationWorker) Copy() NotificationWorker {
	return NotificationWorker{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x NotificationWorker) WithName(v string) NotificationWorker {
	y := x.Copy()
	y.Name = v
	return y
}
func (x NotificationWorker) WithEnabled(v bool) NotificationWorker {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x NotificationWorker) WithTimeout(v time.Duration) NotificationWorker {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x NotificationWorker) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x NotificationWorker) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
