package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"time"
)

type ResponseWorkflow struct {
	DB        *sql.DB
	Responses repository.ResponseStore
	Now       func() time.Time
}

func (w ResponseWorkflow) Submit(ctx context.Context, id, actor, body string) error {
	if e := domain.ValidateResponse(body); e != nil {
		return e
	}
	var status string
	if e := w.DB.QueryRowContext(ctx, `SELECT status FROM suggestions WHERE id=(SELECT suggestion_id FROM responses WHERE id=?)`, id).Scan(&status); e != nil {
		return e
	}
	if status != string(domain.StatusInProgress) && status != string(domain.StatusReopened) {
		return fmt.Errorf("%w: suggestion state", domain.ErrConflict)
	}
	return w.Responses.Submit(ctx, id)
}
func (w ResponseWorkflow) Approve(ctx context.Context, id, actor string) error {
	tx, e := w.DB.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	res, e := tx.ExecContext(ctx, `UPDATE responses SET status='approved',reviewed_by=? WHERE id=? AND status='pending_review'`, actor, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return tx.Commit()
}
func (w ResponseWorkflow) Publish(ctx context.Context, id string) error {
	tx, e := w.DB.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	var suggestion, status string
	if e = tx.QueryRowContext(ctx, `SELECT suggestion_id,status FROM responses WHERE id=?`, id).Scan(&suggestion, &status); e != nil {
		return e
	}
	if status != "approved" {
		return domain.ErrConflict
	}
	if _, e = tx.ExecContext(ctx, `UPDATE responses SET status='published',published_at=? WHERE id=?`, w.Now().Format(time.RFC3339Nano), id); e != nil {
		return e
	}
	if _, e = tx.ExecContext(ctx, `UPDATE suggestions SET status='responded',version=version+1,updated_at=? WHERE id=? AND status='in_progress'`, w.Now().Format(time.RFC3339Nano), suggestion); e != nil {
		return e
	}
	return tx.Commit()
}
