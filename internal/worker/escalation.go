package worker

import (
	"context"
	"fmt"
	"time"
)

// EscalationWorker owns deadline escalation worker.
type EscalationWorker struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x EscalationWorker) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("EscalationWorker: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("EscalationWorker: timeout cannot be negative")
	}
	return nil
}
func (x EscalationWorker) Start(ctx context.Context) error {
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
func (x EscalationWorker) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x EscalationWorker) Ready() bool { return x.Enabled && x.Name != "" }
func (x EscalationWorker) Copy() EscalationWorker {
	return EscalationWorker{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x EscalationWorker) WithName(v string) EscalationWorker { y := x.Copy(); y.Name = v; return y }
func (x EscalationWorker) WithEnabled(v bool) EscalationWorker {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x EscalationWorker) WithTimeout(v time.Duration) EscalationWorker {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x EscalationWorker) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x EscalationWorker) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
