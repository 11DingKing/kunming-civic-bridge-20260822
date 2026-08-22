package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type ReviewStore struct{ DB *sql.DB }

func (r ReviewStore) Claim(ctx context.Context, suggestion, actor string, version int) error {
	var decision string
	if e := r.DB.QueryRowContext(ctx, `SELECT decision FROM reviews WHERE suggestion_id=?`, suggestion).Scan(&decision); e != nil {
		return e
	}
	if decision != "unassigned" {
		return domain.ErrConflict
	}
	res, e := r.DB.ExecContext(ctx, `UPDATE reviews SET reviewer_id=?,decision='claimed',version=version+1 WHERE suggestion_id=?`, actor, suggestion)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r ReviewStore) Complete(ctx context.Context, suggestion, actor string, decision domain.ReviewDecision, note string) error {
	if e := domain.ValidateReview(decision, note); e != nil {
		return e
	}
	_, e := r.DB.ExecContext(ctx, `UPDATE reviews SET reviewer_id=?,decision=?,note=?,version=version+1 WHERE suggestion_id=? AND decision='claimed'`, actor, decision, note, suggestion)
	return e
}
func (r ReviewStore) Create(ctx context.Context, suggestion string) error {
	if suggestion == "" {
		return fmt.Errorf("%w: suggestion", domain.ErrInvalid)
	}
	_, e := r.DB.ExecContext(ctx, `INSERT INTO reviews(id,suggestion_id,reviewer_id,decision,note,version,created_at) VALUES(?,?,?,'unassigned','',0,?)`, platform.ID(), suggestion, "", time.Now().UTC().Format(time.RFC3339Nano))
	return e
}
