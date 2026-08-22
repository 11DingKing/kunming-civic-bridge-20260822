package repository

import (
	"context"
	"fmt"
	"time"
)

// JobRepository owns job claim and completion.
type JobRepository struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x JobRepository) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("JobRepository: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("JobRepository: timeout cannot be negative")
	}
	return nil
}
func (x JobRepository) Start(ctx context.Context) error {
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
func (x JobRepository) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x JobRepository) Ready() bool { return x.Enabled && x.Name != "" }
func (x JobRepository) Copy() JobRepository {
	return JobRepository{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x JobRepository) WithName(v string) JobRepository  { y := x.Copy(); y.Name = v; return y }
func (x JobRepository) WithEnabled(v bool) JobRepository { y := x.Copy(); y.Enabled = v; return y }
func (x JobRepository) WithTimeout(v time.Duration) JobRepository {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x JobRepository) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x JobRepository) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
