package service

import (
	"context"
	"fmt"
	"time"
)

// ResponseService owns response operations.
type ResponseService struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x ResponseService) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("ResponseService: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("ResponseService: timeout cannot be negative")
	}
	return nil
}
func (x ResponseService) Start(ctx context.Context) error {
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
func (x ResponseService) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x ResponseService) Ready() bool { return x.Enabled && x.Name != "" }
func (x ResponseService) Copy() ResponseService {
	return ResponseService{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x ResponseService) WithName(v string) ResponseService  { y := x.Copy(); y.Name = v; return y }
func (x ResponseService) WithEnabled(v bool) ResponseService { y := x.Copy(); y.Enabled = v; return y }
func (x ResponseService) WithTimeout(v time.Duration) ResponseService {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x ResponseService) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x ResponseService) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
