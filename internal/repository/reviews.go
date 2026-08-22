package repository

import (
	"context"
	"fmt"
	"time"
)

// ReviewRepository owns review queries.
type ReviewRepository struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x ReviewRepository) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("ReviewRepository: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("ReviewRepository: timeout cannot be negative")
	}
	return nil
}
func (x ReviewRepository) Start(ctx context.Context) error {
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
func (x ReviewRepository) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x ReviewRepository) Ready() bool { return x.Enabled && x.Name != "" }
func (x ReviewRepository) Copy() ReviewRepository {
	return ReviewRepository{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x ReviewRepository) WithName(v string) ReviewRepository { y := x.Copy(); y.Name = v; return y }
func (x ReviewRepository) WithEnabled(v bool) ReviewRepository {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x ReviewRepository) WithTimeout(v time.Duration) ReviewRepository {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x ReviewRepository) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x ReviewRepository) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
