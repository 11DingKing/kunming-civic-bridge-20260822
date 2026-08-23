package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
)

type SuggestionEventStore struct{ DB *sql.DB }

func (s SuggestionEventStore) Count(ctx context.Context, id string) (int, error) {
	var n int
	e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestion_events WHERE suggestion_id=?`, id).Scan(&n)
	return n, e
}
func (s SuggestionEventStore) Last(ctx context.Context, id string) (string, error) {
	var to string
	e := s.DB.QueryRowContext(ctx, `SELECT to_status FROM suggestion_events WHERE suggestion_id=? ORDER BY created_at DESC LIMIT 1`, id).Scan(&to)
	if e == sql.ErrNoRows {
		return "", domain.ErrNotFound
	}
	return to, e
}
