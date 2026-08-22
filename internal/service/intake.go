package service

import (
	"context"
	"fmt"
	"time"
)

// IntakeService owns intake operations.
type IntakeService struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x IntakeService) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("IntakeService: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("IntakeService: timeout cannot be negative")
	}
	return nil
}
func (x IntakeService) Start(ctx context.Context) error {
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
func (x IntakeService) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x IntakeService) Ready() bool { return x.Enabled && x.Name != "" }
func (x IntakeService) Copy() IntakeService {
	return IntakeService{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x IntakeService) WithName(v string) IntakeService  { y := x.Copy(); y.Name = v; return y }
func (x IntakeService) WithEnabled(v bool) IntakeService { y := x.Copy(); y.Enabled = v; return y }
func (x IntakeService) WithTimeout(v time.Duration) IntakeService {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x IntakeService) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x IntakeService) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
