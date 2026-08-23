package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type DepartmentStore struct{ DB *sql.DB }
type DepartmentRecord struct {
	ID, Name, Scope string
	Active          bool
	CreatedAt       time.Time
}

func (d DepartmentStore) Create(ctx context.Context, name, scope string) (string, error) {
	if name == "" || scope == "" {
		return "", domain.ErrInvalid
	}
	id := platform.ID()
	_, e := d.DB.ExecContext(ctx, `INSERT INTO departments(id,name,scope,active,created_at) VALUES(?,?,?,1,?)`, id, name, scope, time.Now().UTC().Format(time.RFC3339Nano))
	return id, e
}
func (d DepartmentStore) ActiveForScope(ctx context.Context, scope string) ([]DepartmentRecord, error) {
	rows, e := d.DB.QueryContext(ctx, `SELECT id,name,scope,active,created_at FROM departments WHERE scope=? AND active=1 ORDER BY name`, scope)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []DepartmentRecord
	for rows.Next() {
		var x DepartmentRecord
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
func (d DepartmentStore) Disable(ctx context.Context, id string) error {
	res, e := d.DB.ExecContext(ctx, `UPDATE departments SET active=0 WHERE id=?`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrNotFound
	}
	return nil
}
