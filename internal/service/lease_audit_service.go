package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

type LeaseAuditService struct {
	DB    *sql.DB
	Store repository.LeaseAuditStore
}

func (s LeaseAuditService) Record(ctx context.Context, event domain.AssignmentEvent) error {
	return s.Store.Record(ctx, event)
}
func (s LeaseAuditService) Count(ctx context.Context, id string) (int, error) {
	return s.Store.History(ctx, id)
}
