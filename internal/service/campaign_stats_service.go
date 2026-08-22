package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

type CampaignStatsService struct {
	DB    *sql.DB
	Store repository.CampaignStatsStore
}

func (s CampaignStatsService) Get(ctx context.Context, id string) (repository.CampaignStats, error) {
	return s.Store.Get(ctx, id)
}
func (s CampaignStatsService) CompletionRate(ctx context.Context, id string) (float64, error) {
	return s.Store.CompletionRate(ctx, id)
}
