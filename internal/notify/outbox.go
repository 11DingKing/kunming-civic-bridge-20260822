package notify

import (
	"context"
	"fmt"
	"time"
)

// Outbox owns outbox event construction.
type Outbox struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x Outbox) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("Outbox: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("Outbox: timeout cannot be negative")
	}
	return nil
}
func (x Outbox) Start(ctx context.Context) error {
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
func (x Outbox) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x Outbox) Ready() bool                        { return x.Enabled && x.Name != "" }
func (x Outbox) Copy() Outbox                       { return Outbox{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout} }
func (x Outbox) WithName(v string) Outbox           { y := x.Copy(); y.Name = v; return y }
func (x Outbox) WithEnabled(v bool) Outbox          { y := x.Copy(); y.Enabled = v; return y }
func (x Outbox) WithTimeout(v time.Duration) Outbox { y := x.Copy(); y.Timeout = v; return y }
func (x Outbox) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x Outbox) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
