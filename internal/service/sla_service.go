package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"time"
)

type SLAService struct {
	DB    *sql.DB
	Store repository.SLAStore
}

func (s SLAService) Summary(ctx context.Context, now time.Time) (repository.SLASummary, error) {
	return s.Store.Summary(ctx, now)
}
func (s SLAService) Due(ctx context.Context, now time.Time) ([]string, error) {
	return s.Store.Due(ctx, now)
}
func (s SLAService) SetDue(ctx context.Context, id string, due time.Time) error {
	return s.Store.SetDue(ctx, id, due)
}
