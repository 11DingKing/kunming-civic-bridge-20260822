package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
)

type ApprovalStore struct{ DB *sql.DB }

func (s ApprovalStore) Pending(ctx context.Context, limit int) ([]string, error) {
	if limit < 1 {
		limit = 20
	}
	rows, e := s.DB.QueryContext(ctx, `SELECT id FROM responses WHERE status='pending_review' ORDER BY created_at LIMIT ?`, limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var id string
		if e = rows.Scan(&id); e != nil {
			return nil, e
		}
		out = append(out, id)
	}
	return out, rows.Err()
}
func (s ApprovalStore) Reject(ctx context.Context, id, note string) error {
	if note == "" {
		return domain.ErrInvalid
	}
	res, e := s.DB.ExecContext(ctx, `UPDATE responses SET status='draft' WHERE id=? AND status='pending_review'`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
