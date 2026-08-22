package auth

import (
	"context"
	"fmt"
	"time"
)

// SessionStore owns session issue, revoke, expiry.
type SessionStore struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x SessionStore) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("SessionStore: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("SessionStore: timeout cannot be negative")
	}
	return nil
}
func (x SessionStore) Start(ctx context.Context) error {
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
func (x SessionStore) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x SessionStore) Ready() bool { return x.Enabled && x.Name != "" }
func (x SessionStore) Copy() SessionStore {
	return SessionStore{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x SessionStore) WithName(v string) SessionStore  { y := x.Copy(); y.Name = v; return y }
func (x SessionStore) WithEnabled(v bool) SessionStore { y := x.Copy(); y.Enabled = v; return y }
func (x SessionStore) WithTimeout(v time.Duration) SessionStore {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x SessionStore) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x SessionStore) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
