package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

type PointReportService struct {
	DB    *sql.DB
	Store repository.PointReportStore
}

func (s PointReportService) Reports(ctx context.Context) ([]repository.PointReport, error) {
	return s.Store.Reports(ctx)
}
func (s PointReportService) RequireActive(ctx context.Context, id string) error {
	return s.Store.RequireActive(ctx, id)
}
