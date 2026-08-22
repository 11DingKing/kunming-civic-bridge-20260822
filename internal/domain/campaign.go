package domain

import (
	"context"
	"fmt"
	"time"
)

// CampaignPolicy owns campaign lifecycle, timezone windows, close/reopen rules.
type CampaignPolicy struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x CampaignPolicy) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("CampaignPolicy: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("CampaignPolicy: timeout cannot be negative")
	}
	return nil
}
func (x CampaignPolicy) Start(ctx context.Context) error {
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
func (x CampaignPolicy) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x CampaignPolicy) Ready() bool { return x.Enabled && x.Name != "" }
func (x CampaignPolicy) Copy() CampaignPolicy {
	return CampaignPolicy{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x CampaignPolicy) WithName(v string) CampaignPolicy  { y := x.Copy(); y.Name = v; return y }
func (x CampaignPolicy) WithEnabled(v bool) CampaignPolicy { y := x.Copy(); y.Enabled = v; return y }
func (x CampaignPolicy) WithTimeout(v time.Duration) CampaignPolicy {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x CampaignPolicy) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x CampaignPolicy) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
