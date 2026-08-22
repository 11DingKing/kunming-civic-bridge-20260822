package repository

import (
	"context"
	"fmt"
	"time"
)

// CampaignRepository owns campaign queries.
type CampaignRepository struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x CampaignRepository) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("CampaignRepository: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("CampaignRepository: timeout cannot be negative")
	}
	return nil
}
func (x CampaignRepository) Start(ctx context.Context) error {
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
func (x CampaignRepository) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x CampaignRepository) Ready() bool { return x.Enabled && x.Name != "" }
func (x CampaignRepository) Copy() CampaignRepository {
	return CampaignRepository{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x CampaignRepository) WithName(v string) CampaignRepository {
	y := x.Copy()
	y.Name = v
	return y
}
func (x CampaignRepository) WithEnabled(v bool) CampaignRepository {
	y := x.Copy()
	y.Enabled = v
	return y
}
func (x CampaignRepository) WithTimeout(v time.Duration) CampaignRepository {
	y := x.Copy()
	y.Timeout = v
	return y
}
func (x CampaignRepository) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x CampaignRepository) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
