package repository

import (
	"context"
	"database/sql"
)

type FeedbackReportStore struct{ DB *sql.DB }
type FeedbackReport struct {
	Count         int
	Average       float64
	NeedsFollowUp int
}

func (s FeedbackReportStore) Get(ctx context.Context, suggestion string) (FeedbackReport, error) {
	var out FeedbackReport
	var avg sql.NullFloat64
	e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*),AVG(rating),SUM(CASE WHEN rating<=2 THEN 1 ELSE 0 END) FROM feedback WHERE suggestion_id=?`, suggestion).Scan(&out.Count, &avg, &out.NeedsFollowUp)
	if e != nil {
		return out, e
	}
	if avg.Valid {
		out.Average = avg.Float64
	}
	return out, nil
}
func (s FeedbackReportStore) PublicCount(ctx context.Context, suggestion string) (int, error) {
	var n int
	e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM feedback WHERE suggestion_id=? AND rating>=3`, suggestion).Scan(&n)
	return n, e
}
