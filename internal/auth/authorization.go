package auth

import (
	"context"
	"fmt"
	"time"
)

// Authorizer owns role and scope authorization.
type Authorizer struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x Authorizer) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("Authorizer: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("Authorizer: timeout cannot be negative")
	}
	return nil
}
func (x Authorizer) Start(ctx context.Context) error {
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
func (x Authorizer) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x Authorizer) Ready() bool { return x.Enabled && x.Name != "" }
func (x Authorizer) Copy() Authorizer {
	return Authorizer{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x Authorizer) WithName(v string) Authorizer           { y := x.Copy(); y.Name = v; return y }
func (x Authorizer) WithEnabled(v bool) Authorizer          { y := x.Copy(); y.Enabled = v; return y }
func (x Authorizer) WithTimeout(v time.Duration) Authorizer { y := x.Copy(); y.Timeout = v; return y }
func (x Authorizer) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x Authorizer) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
