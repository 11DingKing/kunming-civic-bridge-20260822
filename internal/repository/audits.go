package repository

import (
	"context"
	"fmt"
	"time"
)

// AuditRepository owns audit queries.
type AuditRepository struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x AuditRepository) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("AuditRepository: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("AuditRepository: timeout cannot be negative")
	}
	return nil
}
func (x AuditRepository) Start(ctx context.Context) error {
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
func (x AuditRepository) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x AuditRepository) Ready() bool { return x.Enabled && x.Name != "" }
func (x AuditRepository) Copy() AuditRepository {
	return AuditRepository{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x AuditRepository) WithName(v string) AuditRepository  { y := x.Copy(); y.Name = v; return y }
func (x AuditRepository) WithEnabled(v bool) AuditRepository { y := x.Copy(); y.Enabled = v; return y }
func (x AuditRepository) WithTimeout(v time.Duration) AuditRepository {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x AuditRepository) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x AuditRepository) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
