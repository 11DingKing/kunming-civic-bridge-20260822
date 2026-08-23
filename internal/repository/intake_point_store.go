package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type IntakePointStore struct{ DB *sql.DB }
type IntakePointRecord struct {
	ID, Name, Scope string
	Active          bool
	CreatedAt       time.Time
}

func (s IntakePointStore) Create(ctx context.Context, name, scope string) (string, error) {
	if name == "" || scope == "" {
		return "", domain.ErrInvalid
	}
	id := platform.ID()
	_, e := s.DB.ExecContext(ctx, `INSERT INTO intake_points(id,name,scope,active,created_at) VALUES(?,?,?,1,?)`, id, name, scope, time.Now().UTC().Format(time.RFC3339Nano))
	return id, e
}
func (s IntakePointStore) Get(ctx context.Context, id string) (IntakePointRecord, error) {
	var x IntakePointRecord
	var active int
	var created string
	e := s.DB.QueryRowContext(ctx, `SELECT id,name,scope,active,created_at FROM intake_points WHERE id=?`, id).Scan(&x.ID, &x.Name, &x.Scope, &active, &created)
	if e != nil {
		return x, e
	}
	x.Active = active == 1
	x.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return x, nil
}
func (s IntakePointStore) SetActive(ctx context.Context, id string, active bool) error {
	v := 0
	if active {
		v = 1
	}
	res, e := s.DB.ExecContext(ctx, `UPDATE intake_points SET active=? WHERE id=?`, v, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrNotFound
	}
	return nil
}
func (s IntakePointStore) ForScope(ctx context.Context, scope string) ([]IntakePointRecord, error) {
	rows, e := s.DB.QueryContext(ctx, `SELECT id,name,scope,active,created_at FROM intake_points WHERE scope=? ORDER BY name`, scope)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []IntakePointRecord
	for rows.Next() {
		var x IntakePointRecord
		var active int
		var created string
		if e = rows.Scan(&x.ID, &x.Name, &x.Scope, &active, &created); e != nil {
			return nil, e
		}
		x.Active = active == 1
		x.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		out = append(out, x)
	}
	return out, rows.Err()
}
