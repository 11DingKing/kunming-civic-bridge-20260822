package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"time"
)

type SessionQueryService struct {
	DB    *sql.DB
	Store repository.SessionStore
}

func (s SessionQueryService) Active(ctx context.Context, user string, now time.Time) ([]repository.SessionRecord, error) {
	return s.Store.ActiveForUser(ctx, user, now)
}
func (s SessionQueryService) RevokeAll(ctx context.Context, user string, now time.Time) error {
	return s.Store.RevokeUser(ctx, user, now)
}
