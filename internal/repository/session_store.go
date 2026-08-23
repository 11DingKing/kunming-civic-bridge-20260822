package repository

import (
	"context"
	"database/sql"
	"time"
)

type SessionStore struct{ DB *sql.DB }
type SessionRecord struct {
	ID, UserID           string
	ExpiresAt, CreatedAt time.Time
	Revoked              bool
}

func (s SessionStore) Get(ctx context.Context, id string) (SessionRecord, error) {
	var x SessionRecord
	var exp, created string
	var revoked sql.NullString
	e := s.DB.QueryRowContext(ctx, `SELECT id,user_id,expires_at,created_at,revoked_at FROM sessions WHERE id=?`, id).Scan(&x.ID, &x.UserID, &exp, &created, &revoked)
	if e != nil {
		return x, e
	}
	x.ExpiresAt, _ = time.Parse(time.RFC3339Nano, exp)
	x.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	x.Revoked = revoked.Valid
	return x, nil
}
func (s SessionStore) ActiveForUser(ctx context.Context, user string, now time.Time) ([]SessionRecord, error) {
	rows, e := s.DB.QueryContext(ctx, `SELECT id,user_id,expires_at,created_at,revoked_at FROM sessions WHERE user_id=? AND revoked_at IS NULL AND expires_at>? ORDER BY created_at DESC`, user, now.Format(time.RFC3339Nano))
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []SessionRecord
	for rows.Next() {
		var x SessionRecord
		var exp, created string
		var revoked sql.NullString
		if e = rows.Scan(&x.ID, &x.UserID, &exp, &created, &revoked); e != nil {
			return nil, e
		}
		x.ExpiresAt, _ = time.Parse(time.RFC3339Nano, exp)
		x.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		x.Revoked = revoked.Valid
		out = append(out, x)
	}
	return out, rows.Err()
}
func (s SessionStore) RevokeUser(ctx context.Context, user string, now time.Time) error {
	_, e := s.DB.ExecContext(ctx, `UPDATE sessions SET revoked_at=? WHERE user_id=? AND revoked_at IS NULL`, now.Format(time.RFC3339Nano), user)
	return e
}
