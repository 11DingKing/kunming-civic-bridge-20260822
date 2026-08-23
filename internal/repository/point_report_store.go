package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
)

type PointReportStore struct{ DB *sql.DB }
type PointReport struct {
	PointID, Name, Scope        string
	Active                      bool
	Received, Submitted, Closed int
}

func (s PointReportStore) Reports(ctx context.Context) ([]PointReport, error) {
	rows, e := s.DB.QueryContext(ctx, `SELECT id,name,scope,active FROM intake_points ORDER BY name`)
	if e != nil {
		return nil, e
	}
	var out []PointReport
	for rows.Next() {
		var x PointReport
		var active int
		if e = rows.Scan(&x.PointID, &x.Name, &x.Scope, &active); e != nil {
			return nil, e
		}
		x.Active = active == 1
		out = append(out, x)
	}
	if e = rows.Err(); e != nil {
		rows.Close()
		return nil, e
	}
	rows.Close()
	for i := range out {
		if e = s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestions WHERE intake_point_id=?`, out[i].PointID).Scan(&out[i].Received); e != nil {
			return nil, e
		}
		if e = s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestions WHERE intake_point_id=? AND status='submitted'`, out[i].PointID).Scan(&out[i].Submitted); e != nil {
			return nil, e
		}
		if e = s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestions WHERE intake_point_id=? AND status='closed'`, out[i].PointID).Scan(&out[i].Closed); e != nil {
			return nil, e
		}
	}
	return out, nil
}
func (s PointReportStore) RequireActive(ctx context.Context, id string) error {
	var active int
	e := s.DB.QueryRowContext(ctx, `SELECT active FROM intake_points WHERE id=?`, id).Scan(&active)
	if e == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	if e != nil {
		return e
	}
	if active != 1 {
		return domain.ErrConflict
	}
	return nil
}
