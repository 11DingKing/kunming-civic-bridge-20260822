package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"strings"
)

type TriageQuery struct{ DB *sql.DB }

func (q TriageQuery) Pending(ctx context.Context, scope string) ([]string, error) {
	rows, e := q.DB.QueryContext(ctx, `SELECT id FROM suggestions WHERE status='submitted' AND (?='' OR scope=?) ORDER BY created_at`, scope, scope)
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
func (q TriageQuery) Category(ctx context.Context, suggestion string) (string, error) {
	var note string
	e := q.DB.QueryRowContext(ctx, `SELECT note FROM reviews WHERE suggestion_id=? AND decision='triage' ORDER BY created_at DESC LIMIT 1`, suggestion).Scan(&note)
	if e != nil {
		return "", e
	}
	parts := strings.SplitN(note, ":", 2)
	if len(parts) != 2 {
		return "", domain.ErrInvalid
	}
	return parts[0], nil
}
