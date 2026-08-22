package service

import (
	"context"
	"fmt"
	"time"
)

// CampaignService owns campaign operations.
type CampaignService struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x CampaignService) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("CampaignService: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("CampaignService: timeout cannot be negative")
	}
	return nil
}
func (x CampaignService) Start(ctx context.Context) error {
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
func (x CampaignService) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x CampaignService) Ready() bool { return x.Enabled && x.Name != "" }
func (x CampaignService) Copy() CampaignService {
	return CampaignService{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x CampaignService) WithName(v string) CampaignService  { y := x.Copy(); y.Name = v; return y }
func (x CampaignService) WithEnabled(v bool) CampaignService { y := x.Copy(); y.Enabled = v; return y }
func (x CampaignService) WithTimeout(v time.Duration) CampaignService {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x CampaignService) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x CampaignService) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
