package httpapi

import (
	"context"
	"fmt"
	"time"
)

// ResponseHandlers owns response endpoints.
type ResponseHandlers struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x ResponseHandlers) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("ResponseHandlers: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("ResponseHandlers: timeout cannot be negative")
	}
	return nil
}
func (x ResponseHandlers) Start(ctx context.Context) error {
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
func (x ResponseHandlers) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x ResponseHandlers) Ready() bool { return x.Enabled && x.Name != "" }
func (x ResponseHandlers) Copy() ResponseHandlers {
	return ResponseHandlers{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x ResponseHandlers) WithName(v string) ResponseHandlers { y := x.Copy(); y.Name = v; return y }
func (x ResponseHandlers) WithEnabled(v bool) ResponseHandlers {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x ResponseHandlers) WithTimeout(v time.Duration) ResponseHandlers {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x ResponseHandlers) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x ResponseHandlers) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
