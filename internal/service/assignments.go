package service

import (
	"context"
	"fmt"
	"time"
)

// AssignmentService owns assignment operations.
type AssignmentService struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x AssignmentService) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("AssignmentService: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("AssignmentService: timeout cannot be negative")
	}
	return nil
}
func (x AssignmentService) Start(ctx context.Context) error {
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
func (x AssignmentService) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x AssignmentService) Ready() bool { return x.Enabled && x.Name != "" }
func (x AssignmentService) Copy() AssignmentService {
	return AssignmentService{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x AssignmentService) WithName(v string) AssignmentService { y := x.Copy(); y.Name = v; return y }
func (x AssignmentService) WithEnabled(v bool) AssignmentService {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x AssignmentService) WithTimeout(v time.Duration) AssignmentService {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x AssignmentService) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x AssignmentService) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
