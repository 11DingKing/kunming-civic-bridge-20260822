package repository

import (
	"context"
	"fmt"
	"time"
)

// AssignmentRepository owns assignment queries.
type AssignmentRepository struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x AssignmentRepository) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("AssignmentRepository: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("AssignmentRepository: timeout cannot be negative")
	}
	return nil
}
func (x AssignmentRepository) Start(ctx context.Context) error {
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
func (x AssignmentRepository) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x AssignmentRepository) Ready() bool { return x.Enabled && x.Name != "" }
func (x AssignmentRepository) Copy() AssignmentRepository {
	return AssignmentRepository{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x AssignmentRepository) WithName(v string) AssignmentRepository {
	y := x.Copy()
	y.Name = v
	return y
}
func (x AssignmentRepository) WithEnabled(v bool) AssignmentRepository {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x AssignmentRepository) WithTimeout(v time.Duration) AssignmentRepository {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x AssignmentRepository) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x AssignmentRepository) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
