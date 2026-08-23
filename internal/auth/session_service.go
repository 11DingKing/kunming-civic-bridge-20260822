package auth

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type SessionService struct {
	DB  *sql.DB
	TTL time.Duration
	Now func() time.Time
}

func (s SessionService) Issue(ctx context.Context, user string) (string, error) {
	if user == "" {
		return "", fmt.Errorf("%w: user", domain.ErrInvalid)
	}
	now := s.Now()
	id := platform.ID()
	_, e := s.DB.ExecContext(ctx, `INSERT INTO sessions(id,user_id,expires_at,created_at) VALUES(?,?,?,?)`, id, user, now.Add(s.TTL).Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	return id, e
}
func (s SessionService) Revoke(ctx context.Context, id string) error {
	r, e := s.DB.ExecContext(ctx, `UPDATE sessions SET revoked_at=? WHERE id=? AND revoked_at IS NULL`, s.Now().Format(time.RFC3339Nano), id)
	if e != nil {
		return e
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
func (s SessionService) Resolve(ctx context.Context, id string) (string, error) {
	var user, expires string
	var revoked sql.NullString
	if e := s.DB.QueryRowContext(ctx, `SELECT user_id,expires_at,revoked_at FROM sessions WHERE id=?`, id).Scan(&user, &expires, &revoked); e != nil {
		if e == sql.ErrNoRows {
			return "", domain.ErrNotFound
		}
		return "", e
	}
	if revoked.Valid {
		return "", domain.ErrForbidden
	}
	until, e := time.Parse(time.RFC3339Nano, expires)
	if e != nil {
		return "", e
	}
	if !s.Now().Before(until) {
		return "", domain.ErrExpired
	}
	return user, nil
}
