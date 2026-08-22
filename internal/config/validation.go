package config

import (
	"context"
	"fmt"
	"time"
)

// ConfigValidator owns runtime configuration validation.
type ConfigValidator struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x ConfigValidator) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("ConfigValidator: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("ConfigValidator: timeout cannot be negative")
	}
	return nil
}
func (x ConfigValidator) Start(ctx context.Context) error {
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
func (x ConfigValidator) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x ConfigValidator) Ready() bool { return x.Enabled && x.Name != "" }
func (x ConfigValidator) Copy() ConfigValidator {
	return ConfigValidator{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x ConfigValidator) WithName(v string) ConfigValidator  { y := x.Copy(); y.Name = v; return y }
func (x ConfigValidator) WithEnabled(v bool) ConfigValidator { y := x.Copy(); y.Enabled = v; return y }
func (x ConfigValidator) WithTimeout(v time.Duration) ConfigValidator {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x ConfigValidator) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x ConfigValidator) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
