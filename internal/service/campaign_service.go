package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"time"
)

type CampaignService struct {
	DB    *sql.DB
	Store repository.CampaignStore
	Now   func() time.Time
}

func (c CampaignService) Open(ctx context.Context, id string) error {
	r, e := c.Store.Window(ctx, id)
	if e != nil {
		return e
	}
	if e = r.CanOpen(c.Now()); e != nil {
		return e
	}
	return c.Store.SetStatus(ctx, id, domain.CampaignOpen)
}
func (c CampaignService) Pause(ctx context.Context, id string) error {
	r, e := c.Store.Window(ctx, id)
	if e != nil {
		return e
	}
	if e = r.Pause(); e != nil {
		return e
	}
	return c.Store.SetStatus(ctx, id, domain.CampaignPaused)
}
func (c CampaignService) Close(ctx context.Context, id string) error {
	r, e := c.Store.Window(ctx, id)
	if e != nil {
		return e
	}
	if e = r.Close(); e != nil {
		return e
	}
	return c.Store.SetStatus(ctx, id, domain.CampaignClosed)
}
func (c CampaignService) Accepting(ctx context.Context, id string) error {
	r, e := c.Store.Window(ctx, id)
	if e != nil {
		return e
	}
	if e = r.CanAccept(c.Now()); e != nil {
		return fmt.Errorf("campaign %s: %w", id, e)
	}
	return nil
}
