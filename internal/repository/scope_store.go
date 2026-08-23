package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
)

type ScopeStore struct{ DB *sql.DB }

func (s ScopeStore) UserCanSee(ctx context.Context, user, suggestion string) (bool, error) {
	var userScope, suggestionScope string
	if e := s.DB.QueryRowContext(ctx, `SELECT scope FROM users WHERE id=?`, user).Scan(&userScope); e != nil {
		return false, e
	}
	if e := s.DB.QueryRowContext(ctx, `SELECT scope FROM suggestions WHERE id=?`, suggestion).Scan(&suggestionScope); e != nil {
		return false, e
	}
	a, e := domain.ParseScope(userScope)
	if e != nil {
		return false, e
	}
	b, e := domain.ParseScope(suggestionScope)
	if e != nil {
		return false, e
	}
	return a.Contains(b), nil
}
func (s ScopeStore) SuggestionsForUser(ctx context.Context, user string) ([]string, error) {
	var scope string
	if e := s.DB.QueryRowContext(ctx, `SELECT scope FROM users WHERE id=?`, user).Scan(&scope); e != nil {
		return nil, e
	}
	rows, e := s.DB.QueryContext(ctx, `SELECT id FROM suggestions WHERE scope=? OR scope LIKE ? ORDER BY updated_at DESC`, scope, scope+"/%")
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
