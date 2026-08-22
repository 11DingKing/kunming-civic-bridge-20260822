package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
)

type UserStore struct{ DB *sql.DB }

func (u UserStore) Roles(ctx context.Context, user string) ([]domain.Role, error) {
	rows, e := u.DB.QueryContext(ctx, `SELECT role FROM roles WHERE user_id=? ORDER BY role`, user)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.Role
	for rows.Next() {
		var role string
		if e = rows.Scan(&role); e != nil {
			return nil, e
		}
		out = append(out, domain.Role(role))
	}
	return out, rows.Err()
}
func (u UserStore) Scope(ctx context.Context, user string) (string, error) {
	var scope string
	e := u.DB.QueryRowContext(ctx, `SELECT scope FROM users WHERE id=?`, user).Scan(&scope)
	if e == sql.ErrNoRows {
		return "", domain.ErrNotFound
	}
	return scope, e
}
