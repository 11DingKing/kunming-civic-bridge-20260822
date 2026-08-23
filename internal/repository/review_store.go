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
	// Atomic compare-and-set: the claim only succeeds when the review is still
	// unassigned at the expected version. The WHERE clause is the sole arbiter of
	// ownership, so two concurrent claims cannot both succeed: the first commits the
	// version bump and the second matches zero rows and returns a conflict.
	res, e := r.DB.ExecContext(ctx, `UPDATE reviews SET reviewer_id=?,decision='claimed',version=version+1 WHERE suggestion_id=? AND decision='unassigned' AND version=?`, actor, suggestion, version)
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
