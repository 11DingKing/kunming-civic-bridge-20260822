package domain

import (
	"fmt"
	"time"
)

type CampaignStatus string

const (
	CampaignDraft    CampaignStatus = "draft"
	CampaignOpen     CampaignStatus = "open"
	CampaignPaused   CampaignStatus = "paused"
	CampaignClosed   CampaignStatus = "closed"
	CampaignArchived CampaignStatus = "archived"
)

type CampaignRules struct {
	Status           CampaignStatus
	StartsAt, EndsAt time.Time
	Timezone         string
	MaxSubmissions   int
}

func (c CampaignRules) Validate() error {
	if c.StartsAt.IsZero() || c.EndsAt.IsZero() || !c.StartsAt.Before(c.EndsAt) {
		return fmt.Errorf("%w: campaign window", ErrInvalid)
	}
	if c.Timezone == "" {
		return fmt.Errorf("%w: campaign timezone", ErrInvalid)
	}
	if c.MaxSubmissions < 1 {
		return fmt.Errorf("%w: campaign capacity", ErrInvalid)
	}
	return nil
}
func (c CampaignRules) CanOpen(now time.Time) error {
	if e := c.Validate(); e != nil {
		return e
	}
	if c.Status != CampaignDraft && c.Status != CampaignPaused {
		return fmt.Errorf("%w: campaign status", ErrConflict)
	}
	if now.After(c.EndsAt) {
		return ErrExpired
	}
	return nil
}
func (c CampaignRules) CanAccept(now time.Time) error {
	if c.Status != CampaignOpen {
		return ErrClosed
	}
	if now.Before(c.StartsAt) {
		return fmt.Errorf("%w: campaign not started", ErrInvalid)
	}
	if !now.Before(c.EndsAt) {
		return ErrExpired
	}
	return nil
}
func (c CampaignRules) Pause() error {
	if c.Status != CampaignOpen {
		return fmt.Errorf("%w: pause status", ErrConflict)
	}
	return nil
}
func (c CampaignRules) Close() error {
	if c.Status != CampaignOpen && c.Status != CampaignPaused {
		return fmt.Errorf("%w: close status", ErrConflict)
	}
	return nil
}
func (c CampaignRules) Archive() error {
	if c.Status != CampaignClosed {
		return fmt.Errorf("%w: archive status", ErrConflict)
	}
	return nil
}
