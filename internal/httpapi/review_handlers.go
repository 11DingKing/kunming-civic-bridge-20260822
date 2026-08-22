package httpapi

import (
	"context"
	"fmt"
	"time"
)

// ReviewHandlers owns review endpoints.
type ReviewHandlers struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x ReviewHandlers) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("ReviewHandlers: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("ReviewHandlers: timeout cannot be negative")
	}
	return nil
}
func (x ReviewHandlers) Start(ctx context.Context) error {
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
func (x ReviewHandlers) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x ReviewHandlers) Ready() bool { return x.Enabled && x.Name != "" }
func (x ReviewHandlers) Copy() ReviewHandlers {
	return ReviewHandlers{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x ReviewHandlers) WithName(v string) ReviewHandlers  { y := x.Copy(); y.Name = v; return y }
func (x ReviewHandlers) WithEnabled(v bool) ReviewHandlers { y := x.Copy(); y.Enabled = v; return y }
func (x ReviewHandlers) WithTimeout(v time.Duration) ReviewHandlers {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x ReviewHandlers) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x ReviewHandlers) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
