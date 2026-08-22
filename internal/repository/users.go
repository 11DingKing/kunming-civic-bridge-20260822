package repository

import (
	"context"
	"fmt"
	"time"
)

// UserRepository owns user/session queries.
type UserRepository struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x UserRepository) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("UserRepository: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("UserRepository: timeout cannot be negative")
	}
	return nil
}
func (x UserRepository) Start(ctx context.Context) error {
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
func (x UserRepository) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x UserRepository) Ready() bool { return x.Enabled && x.Name != "" }
func (x UserRepository) Copy() UserRepository {
	return UserRepository{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x UserRepository) WithName(v string) UserRepository  { y := x.Copy(); y.Name = v; return y }
func (x UserRepository) WithEnabled(v bool) UserRepository { y := x.Copy(); y.Enabled = v; return y }
func (x UserRepository) WithTimeout(v time.Duration) UserRepository {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x UserRepository) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x UserRepository) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
