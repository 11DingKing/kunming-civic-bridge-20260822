package repository

import (
	"context"
	"fmt"
	"time"
)

// ResponseRepository owns response queries.
type ResponseRepository struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x ResponseRepository) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("ResponseRepository: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("ResponseRepository: timeout cannot be negative")
	}
	return nil
}
func (x ResponseRepository) Start(ctx context.Context) error {
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
func (x ResponseRepository) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x ResponseRepository) Ready() bool { return x.Enabled && x.Name != "" }
func (x ResponseRepository) Copy() ResponseRepository {
	return ResponseRepository{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x ResponseRepository) WithName(v string) ResponseRepository {
	y := x.Copy()
	y.Name = v
	return y
}
func (x ResponseRepository) WithEnabled(v bool) ResponseRepository {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x ResponseRepository) WithTimeout(v time.Duration) ResponseRepository {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x ResponseRepository) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x ResponseRepository) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
