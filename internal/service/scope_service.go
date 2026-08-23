package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

type ScopeService struct {
	DB    *sql.DB
	Store repository.ScopeStore
}

func (s ScopeService) CanRead(ctx context.Context, user, suggestion string) (bool, error) {
	return s.Store.UserCanSee(ctx, user, suggestion)
}
func (s ScopeService) VisibleIDs(ctx context.Context, user string) ([]string, error) {
	return s.Store.SuggestionsForUser(ctx, user)
}
