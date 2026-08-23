package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

type VisibilityService struct {
	DB    *sql.DB
	Store repository.VisibilityStore
}

func (s VisibilityService) CanPublish(ctx context.Context, id string) (bool, error) {
	return s.Store.CanPublish(ctx, id)
}
func (s VisibilityService) Redact(ctx context.Context, id, note string) error {
	return s.Store.SetRedacted(ctx, id, note)
}
