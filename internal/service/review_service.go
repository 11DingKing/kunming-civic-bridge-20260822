package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
)

type ReviewService struct{ DB *sql.DB }

func (r ReviewService) Return(ctx context.Context, suggestion, actor, note string) error {
	if e := (domain.ReviewRules{State: domain.ReviewClaimed, ReviewerID: actor}).Return(note); e != nil {
		return e
	}
	_, e := r.DB.ExecContext(ctx, `UPDATE reviews SET decision='returned',note=?,version=version+1 WHERE suggestion_id=? AND reviewer_id=? AND decision='claimed'`, note, suggestion, actor)
	return e
}
