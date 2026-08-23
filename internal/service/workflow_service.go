package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"time"
)

type Workflow struct {
	DB          *sql.DB
	Suggestions repository.SuggestionRepository
	Jobs        repository.JobStore
	Clock       func() time.Time
}

func (w Workflow) ClaimAssignment(ctx context.Context, suggestion, operator, token string, lease time.Duration) error {
	// Claim and lease updates share one transaction so ownership changes are atomic.
	now := w.Clock()
	tx, e := w.DB.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	var id, status, oldToken string
	var version int
	var oldUntil sql.NullString
	e = tx.QueryRowContext(ctx, `SELECT id,status,COALESCE(lease_token,''),lease_until,version FROM assignments WHERE suggestion_id=? ORDER BY assigned_at DESC LIMIT 1`, suggestion).Scan(&id, &status, &oldToken, &oldUntil, &version)
	if e == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	if e != nil {
		return e
	}
	var until time.Time
	if oldUntil.Valid {
		until, _ = time.Parse(time.RFC3339Nano, oldUntil.String)
	}
	if !domain.Claimable(domain.AssignmentStatus(status), domain.Lease{Token: oldToken, OwnerID: operator, Until: until}, now) {
		return domain.ErrConflict
	}
	leaseUntil := now.Add(lease)
	// Guard against concurrent claims with optimistic versioning: the UPDATE only
	// matches the version read above, so the second concurrent writer affects zero
	// rows and is reported as a conflict instead of silently overwriting the winner.
	result, e := tx.ExecContext(ctx, `UPDATE assignments SET status='claimed',assignee_id=?,lease_token=?,lease_until=?,version=version+1 WHERE id=? AND version=?`, operator, token, leaseUntil.Format(time.RFC3339Nano), id, version)
	if e != nil {
		return platform.TranslateSQLiteBusy(e)
	}
	n, _ := result.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return tx.Commit()
}
func (w Workflow) ApproveResponse(ctx context.Context, response, supervisor string) error {
	tx, e := w.DB.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	var suggestionID, status string
	if e = tx.QueryRowContext(ctx, `SELECT suggestion_id,status FROM responses WHERE id=?`, response).Scan(&suggestionID, &status); e != nil {
		if e == sql.ErrNoRows {
			return domain.ErrNotFound
		}
		return e
	}
	if status != string(domain.ResponsePending) {
		return fmt.Errorf("%w: response state", domain.ErrConflict)
	}
	if _, e = tx.ExecContext(ctx, `UPDATE responses SET status='approved',reviewed_by=? WHERE id=? AND status='pending_review'`, supervisor, response); e != nil {
		return e
	}
	if _, e = tx.ExecContext(ctx, `INSERT INTO audit_logs(id,actor_id,action,object_type,object_id,result,request_id,metadata,created_at) VALUES(?,?,?,?,?,?,?,?,?)`, platform.ID(), supervisor, "response.approve", "response", response, "ok", "", `{"suggestion_id":"`+suggestionID+`"}`, w.Clock().Format(time.RFC3339Nano)); e != nil {
		return e
	}
	return tx.Commit()
}
func (w Workflow) ReopenFromFeedback(ctx context.Context, suggestion, author string) error {
	tx, e := w.DB.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	var status string
	var version int
	if e = tx.QueryRowContext(ctx, `SELECT status,version FROM suggestions WHERE id=?`, suggestion).Scan(&status, &version); e != nil {
		if e == sql.ErrNoRows {
			return domain.ErrNotFound
		}
		return e
	}
	if domain.SuggestionStatus(status) != domain.StatusResponded {
		return fmt.Errorf("%w: only responded suggestions reopen", domain.ErrInvalid)
	}
	if _, e = tx.ExecContext(ctx, `UPDATE suggestions SET status='reopened',version=version+1,updated_at=? WHERE id=? AND status='responded' AND version=?`, w.Clock().Format(time.RFC3339Nano), suggestion, version); e != nil {
		return e
	}
	_, e = tx.ExecContext(ctx, `INSERT INTO suggestion_events(id,suggestion_id,from_status,to_status,actor_id,note,created_at) VALUES(?,?,?,?,?,?,?)`, platform.ID(), suggestion, status, "reopened", author, "feedback requested follow-up", w.Clock().Format(time.RFC3339Nano))
	if e != nil {
		return e
	}
	return tx.Commit()
}
func (w Workflow) Archive(ctx context.Context, suggestion, actor string) error {
	tx, e := w.DB.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	var status string
	if e = tx.QueryRowContext(ctx, `SELECT status FROM suggestions WHERE id=?`, suggestion).Scan(&status); e != nil {
		return e
	}
	if status != "closed" {
		return fmt.Errorf("%w: archive requires closed", domain.ErrConflict)
	}
	if _, e = tx.ExecContext(ctx, `UPDATE suggestions SET updated_at=? WHERE id=?`, w.Clock().Format(time.RFC3339Nano), suggestion); e != nil {
		return e
	}
	_, e = tx.ExecContext(ctx, `INSERT INTO audit_logs(id,actor_id,action,object_type,object_id,result,request_id,metadata,created_at) VALUES(?,?,?,?,?,?,?,?,?)`, platform.ID(), actor, "suggestion.archive", "suggestion", suggestion, "ok", "", "{}", w.Clock().Format(time.RFC3339Nano))
	if e != nil {
		return e
	}
	return tx.Commit()
}
