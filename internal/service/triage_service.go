package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

type TriageService struct {
	DB    *sql.DB
	Store repository.TriageStore
	Query repository.TriageQuery
}

func (s TriageService) Classify(ctx context.Context, id, actor string, result domain.TriageResult) error {
	if e := result.Validate(); e != nil {
		return e
	}
	var status string
	if e := s.DB.QueryRowContext(ctx, `SELECT status FROM suggestions WHERE id=?`, id).Scan(&status); e != nil {
		return e
	}
	if status != string(domain.StatusSubmitted) {
		return fmt.Errorf("%w: triage state", domain.ErrConflict)
	}
	if e := s.Store.Save(ctx, id, actor, result); e != nil {
		return e
	}
	_, e := s.DB.ExecContext(ctx, `UPDATE suggestions SET status='triaged',version=version+1,updated_at=CURRENT_TIMESTAMP WHERE id=? AND status='submitted'`, id)
	return e
}
func (s TriageService) Pending(ctx context.Context, scope string) ([]string, error) {
	return s.Query.Pending(ctx, scope)
}
