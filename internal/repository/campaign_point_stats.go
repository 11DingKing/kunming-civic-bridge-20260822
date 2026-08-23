package repository

import (
	"context"
	"database/sql"
)

type PointStatsStore struct{ DB *sql.DB }
type PointStats struct {
	PointID                       string
	Received, Submitted, Rejected int
}

func (s PointStatsStore) Get(ctx context.Context, point string) (PointStats, error) {
	var out PointStats
	out.PointID = point
	if e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestions WHERE intake_point_id=?`, point).Scan(&out.Received); e != nil {
		return out, e
	}
	if e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestions WHERE intake_point_id=? AND status='submitted'`, point).Scan(&out.Submitted); e != nil {
		return out, e
	}
	if e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestions WHERE intake_point_id=? AND status='rejected'`, point).Scan(&out.Rejected); e != nil {
		return out, e
	}
	return out, nil
}
