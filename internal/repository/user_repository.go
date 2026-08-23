package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/auth"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type IdentityStore struct{ DB *sql.DB }

func (s IdentityStore) Create(ctx context.Context, name, phone, password, scope string, role domain.Role) (string, error) {
	hash, e := auth.HashPassword(password)
	if e != nil {
		return "", e
	}
	id := platform.ID()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	tx, e := s.DB.BeginTx(ctx, nil)
	if e != nil {
		return "", e
	}
	defer tx.Rollback()
	if _, e = tx.ExecContext(ctx, `INSERT INTO users(id,name,phone,password_hash,scope,created_at) VALUES(?,?,?,?,?,?)`, id, name, phone, hash, scope, now); e != nil {
		return "", e
	}
	if _, e = tx.ExecContext(ctx, `INSERT INTO roles(user_id,role) VALUES(?,?)`, id, role); e != nil {
		return "", e
	}
	if e = tx.Commit(); e != nil {
		return "", e
	}
	return id, nil
}
func (s IdentityStore) Authenticate(ctx context.Context, phone, password string) (domain.User, error) {
	var u domain.User
	var hash, created string
	e := s.DB.QueryRowContext(ctx, `SELECT id,name,phone,password_hash,scope,created_at FROM users WHERE phone=?`, phone).Scan(&u.ID, &u.Name, &u.Phone, &hash, &u.Scope, &created)
	if e == sql.ErrNoRows {
		return u, domain.ErrForbidden
	}
	if e != nil {
		return u, e
	}
	if !auth.CheckPassword(hash, password) {
		return u, domain.ErrForbidden
	}
	u.PasswordHash = hash
	u.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	rows, e := s.DB.QueryContext(ctx, `SELECT role FROM roles WHERE user_id=?`, u.ID)
	if e != nil {
		return u, e
	}
	defer rows.Close()
	for rows.Next() {
		var role string
		if e = rows.Scan(&role); e != nil {
			return u, e
		}
		u.Roles = append(u.Roles, domain.Role(role))
	}
	return u, rows.Err()
}
func (s IdentityStore) UpdateScope(ctx context.Context, id, scope string) error {
	res, e := s.DB.ExecContext(ctx, `UPDATE users SET scope=? WHERE id=?`, scope, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrNotFound
	}
	return nil
}
