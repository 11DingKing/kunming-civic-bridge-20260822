package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"time"
)

type ExportService struct {
	DB    *sql.DB
	Store repository.ExportStore
}

func (s ExportService) Audit(ctx context.Context, req domain.AuditExport) ([]repository.AuditRecord, error) {
	return s.Store.Audit(ctx, req)
}
func (s ExportService) Activity(ctx context.Context, id string, from, to time.Time) (bool, error) {
	return s.Store.HasActivity(ctx, id, from, to)
}
