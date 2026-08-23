package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"time"
)

type ReopenService struct {
	DB    *sql.DB
	Store repository.ReopenStore
	Now   func() time.Time
}

func (s ReopenService) Request(ctx context.Context, r domain.ReopenRequest) error {
	var status string
	if e := s.DB.QueryRowContext(ctx, `SELECT status FROM suggestions WHERE id=?`, r.SuggestionID).Scan(&status); e != nil {
		return e
	}
	if !r.CanApply(domain.SuggestionStatus(status)) {
		return fmt.Errorf("%w: reopen status", domain.ErrConflict)
	}
	return s.Store.Request(ctx, r)
}
func (s ReopenService) Apply(ctx context.Context, id string) error {
	r, e := s.Store.Latest(ctx, id)
	if e != nil {
		return e
	}
	now := s.Now()
	if !r.WithinWindow(now.Add(-24*time.Hour), now) {
		return domain.ErrExpired
	}
	var status string
	if e = s.DB.QueryRowContext(ctx, `SELECT status FROM suggestions WHERE id=?`, id).Scan(&status); e != nil {
		return e
	}
	if !r.CanApply(domain.SuggestionStatus(status)) {
		return fmt.Errorf("%w: reopen status", domain.ErrConflict)
	}
	_, e = s.DB.ExecContext(ctx, `UPDATE suggestions SET status='reopened',version=version+1,updated_at=? WHERE id=? AND status=?`, now.Format(time.RFC3339Nano), id, status)
	return e
}
