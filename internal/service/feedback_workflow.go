package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"time"
)

type FeedbackWorkflow struct {
	DB       *sql.DB
	Feedback repository.FeedbackStore
	Now      func() time.Time
}

func (w FeedbackWorkflow) Record(ctx context.Context, suggestion, author, comment string, rating int) error {
	published, e := w.published(ctx, suggestion)
	if e != nil {
		return e
	}
	if !published {
		return fmt.Errorf("%w: response unpublished", domain.ErrConflict)
	}
	return w.Feedback.Create(ctx, suggestion, author, comment, rating)
}
func (w FeedbackWorkflow) published(ctx context.Context, id string) (bool, error) {
	var n int
	e := w.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM responses WHERE suggestion_id=? AND status='published'`, id).Scan(&n)
	return n > 0, e
}
func (w FeedbackWorkflow) ReopenIfNeeded(ctx context.Context, suggestion string, rating int) error {
	if rating > 2 {
		return nil
	}
	now := w.Now()
	res, e := w.DB.ExecContext(ctx, `UPDATE suggestions SET status='reopened',version=version+1,updated_at=? WHERE id=? AND status='responded'`, now.Format(time.RFC3339Nano), suggestion)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
