package httpapi

import (
	"context"
	"fmt"
	"time"
)

// AssignmentHandlers owns assignment endpoints.
type AssignmentHandlers struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x AssignmentHandlers) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("AssignmentHandlers: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("AssignmentHandlers: timeout cannot be negative")
	}
	return nil
}
func (x AssignmentHandlers) Start(ctx context.Context) error {
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
func (x AssignmentHandlers) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x AssignmentHandlers) Ready() bool { return x.Enabled && x.Name != "" }
func (x AssignmentHandlers) Copy() AssignmentHandlers {
	return AssignmentHandlers{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x AssignmentHandlers) WithName(v string) AssignmentHandlers {
	y := x.Copy()
	y.Name = v
	return y
}
func (x AssignmentHandlers) WithEnabled(v bool) AssignmentHandlers {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x AssignmentHandlers) WithTimeout(v time.Duration) AssignmentHandlers {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x AssignmentHandlers) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x AssignmentHandlers) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
