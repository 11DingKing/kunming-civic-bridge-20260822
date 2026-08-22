package service

import (
	"context"
	"fmt"
	"time"
)

// SearchService owns filter pagination and scope.
type SearchService struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x SearchService) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("SearchService: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("SearchService: timeout cannot be negative")
	}
	return nil
}
func (x SearchService) Start(ctx context.Context) error {
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
func (x SearchService) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x SearchService) Ready() bool { return x.Enabled && x.Name != "" }
func (x SearchService) Copy() SearchService {
	return SearchService{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x SearchService) WithName(v string) SearchService  { y := x.Copy(); y.Name = v; return y }
func (x SearchService) WithEnabled(v bool) SearchService { y := x.Copy(); y.Enabled = v; return y }
func (x SearchService) WithTimeout(v time.Duration) SearchService {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x SearchService) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x SearchService) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
