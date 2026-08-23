package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
)

type VisibilityStore struct{ DB *sql.DB }

func (s VisibilityStore) CanPublish(ctx context.Context, id string) (bool, error) {
	var status string
	e := s.DB.QueryRowContext(ctx, `SELECT status FROM responses WHERE id=?`, id).Scan(&status)
	if e != nil {
		return false, e
	}
	return status == string(domain.ResponseApproved), nil
}
func (s VisibilityStore) SetRedacted(ctx context.Context, id, note string) error {
	if note == "" {
		return domain.ErrInvalid
	}
	res, e := s.DB.ExecContext(ctx, `UPDATE responses SET status='published' WHERE id=? AND status='approved'`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
