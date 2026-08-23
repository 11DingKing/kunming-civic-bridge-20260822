package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"time"
)

type SLAStore struct{ DB *sql.DB }
type SLASummary struct{ Total, OnTime, Late, Open int }

func (s SLAStore) Summary(ctx context.Context, now time.Time) (SLASummary, error) {
	var out SLASummary
	if e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestions WHERE due_at IS NOT NULL`).Scan(&out.Total); e != nil {
		return out, e
	}
	if e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestions WHERE due_at IS NOT NULL AND status IN ('closed','responded') AND updated_at<=due_at`).Scan(&out.OnTime); e != nil {
		return out, e
	}
	if e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestions WHERE due_at IS NOT NULL AND status NOT IN ('closed','responded') AND due_at<=?`, now.Format(time.RFC3339Nano)).Scan(&out.Late); e != nil {
		return out, e
	}
	if e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestions WHERE status NOT IN ('closed','rejected')`).Scan(&out.Open); e != nil {
		return out, e
	}
	return out, nil
}
func (s SLAStore) Due(ctx context.Context, now time.Time) ([]string, error) {
	rows, e := s.DB.QueryContext(ctx, `SELECT id FROM suggestions WHERE due_at IS NOT NULL AND due_at<=? AND status NOT IN ('closed','rejected') ORDER BY due_at`, now.Format(time.RFC3339Nano))
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
func (s SLAStore) SetDue(ctx context.Context, id string, due time.Time) error {
	res, e := s.DB.ExecContext(ctx, `UPDATE suggestions SET due_at=?,updated_at=? WHERE id=?`, due.Format(time.RFC3339Nano), time.Now().UTC().Format(time.RFC3339Nano), id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrNotFound
	}
	return nil
}
