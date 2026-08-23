package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"time"
)

type ReviewWorkflow struct {
	DB      *sql.DB
	Reviews repository.ReviewStore
	Now     func() time.Time
}

func (w ReviewWorkflow) Claim(ctx context.Context, suggestion, actor string, version int) error {
	var status string
	if e := w.DB.QueryRowContext(ctx, `SELECT status FROM suggestions WHERE id=?`, suggestion).Scan(&status); e != nil {
		return e
	}
	if status != string(domain.StatusSubmitted) {
		return fmt.Errorf("%w: review state", domain.ErrConflict)
	}
	return w.Reviews.Claim(ctx, suggestion, actor, version)
}
func (w ReviewWorkflow) Complete(ctx context.Context, suggestion, actor string, decision domain.ReviewDecision, note string) error {
	if e := domain.ValidateReview(decision, note); e != nil {
		return e
	}
	to := domain.StatusTriaged
	if decision == domain.DecisionReject {
		to = domain.StatusRejected
	}
	tx, e := w.DB.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	var version int
	if e = tx.QueryRowContext(ctx, `SELECT version FROM suggestions WHERE id=?`, suggestion).Scan(&version); e != nil {
		return e
	}
	res, e := tx.ExecContext(ctx, `UPDATE suggestions SET status=?,version=version+1,updated_at=? WHERE id=? AND status='submitted'`, to, w.Now().Format(time.RFC3339Nano), suggestion)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	if _, e = tx.ExecContext(ctx, `UPDATE reviews SET reviewer_id=?,decision=?,note=?,version=version+1 WHERE suggestion_id=? AND decision='claimed'`, actor, decision, note, suggestion); e != nil {
		return e
	}
	_, e = tx.ExecContext(ctx, `INSERT INTO suggestion_events(id,suggestion_id,from_status,to_status,actor_id,note,created_at) VALUES(lower(hex(randomblob(16))),?,?,?,?,?,?)`, suggestion, string(domain.StatusSubmitted), string(to), actor, note, w.Now().Format(time.RFC3339Nano))
	if e != nil {
		return e
	}
	return tx.Commit()
}
