package repository

import (
	"context"
	"database/sql"
)

type AssignmentReportStore struct{ DB *sql.DB }
type AssignmentReport struct{ Total, Claimed, Done, Overdue int }

func (s AssignmentReportStore) Get(ctx context.Context, now string) (AssignmentReport, error) {
	var r AssignmentReport
	if e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM assignments`).Scan(&r.Total); e != nil {
		return r, e
	}
	if e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM assignments WHERE status='claimed'`).Scan(&r.Claimed); e != nil {
		return r, e
	}
	if e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM assignments WHERE status='done'`).Scan(&r.Done); e != nil {
		return r, e
	}
	e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM assignments WHERE lease_until IS NOT NULL AND lease_until<=?`, now).Scan(&r.Overdue)
	return r, e
}
