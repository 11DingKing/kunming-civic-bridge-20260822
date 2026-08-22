package platform

import (
	"context"
	"fmt"
	"time"
)

// ErrorMapper owns platform errors.
type ErrorMapper struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x ErrorMapper) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("ErrorMapper: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("ErrorMapper: timeout cannot be negative")
	}
	return nil
}
func (x ErrorMapper) Start(ctx context.Context) error {
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
func (x ErrorMapper) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x ErrorMapper) Ready() bool { return x.Enabled && x.Name != "" }
func (x ErrorMapper) Copy() ErrorMapper {
	return ErrorMapper{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x ErrorMapper) WithName(v string) ErrorMapper           { y := x.Copy(); y.Name = v; return y }
func (x ErrorMapper) WithEnabled(v bool) ErrorMapper          { y := x.Copy(); y.Enabled = v; return y }
func (x ErrorMapper) WithTimeout(v time.Duration) ErrorMapper { y := x.Copy(); y.Timeout = v; return y }
func (x ErrorMapper) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x ErrorMapper) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
