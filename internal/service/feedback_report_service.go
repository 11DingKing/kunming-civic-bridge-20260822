package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

type FeedbackReportService struct {
	DB    *sql.DB
	Store repository.FeedbackReportStore
}

func (s FeedbackReportService) Get(ctx context.Context, id string) (repository.FeedbackReport, error) {
	return s.Store.Get(ctx, id)
}
func (s FeedbackReportService) Public(ctx context.Context, id string) (int, error) {
	return s.Store.PublicCount(ctx, id)
}
