package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type SuggestionRepository struct{ DB *sql.DB }

func (r SuggestionRepository) Find(ctx context.Context, id string) (domain.Suggestion, error) {
	var x domain.Suggestion
	var status, created, updated string
	var due sql.NullString
	e := r.DB.QueryRowContext(ctx, `SELECT id,campaign_id,COALESCE(intake_point_id,''),author_id,title,body,scope,status,version,due_at,created_at,updated_at FROM suggestions WHERE id=?`, id).Scan(&x.ID, &x.CampaignID, &x.IntakePointID, &x.AuthorID, &x.Title, &x.Body, &x.Scope, &status, &x.Version, &due, &created, &updated)
	if errors.Is(e, sql.ErrNoRows) {
		return x, domain.ErrNotFound
	}
	if e != nil {
		return x, e
	}
	x.Status = domain.SuggestionStatus(status)
	x.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	x.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	if due.Valid {
		v, _ := time.Parse(time.RFC3339Nano, due.String)
		x.DueAt = &v
	}
	return x, nil
}
func (r SuggestionRepository) TransitionTx(ctx context.Context, tx *sql.Tx, id string, from, to domain.SuggestionStatus, version int, now time.Time) error {
	if !from.CanTransition(to) {
		return domain.ErrInvalid
	}
	result, e := tx.ExecContext(ctx, `UPDATE suggestions SET status=?,version=version+1,updated_at=? WHERE id=? AND status=? AND version=?`, to, now.Format(time.RFC3339Nano), id, from, version)
	if e != nil {
		return e
	}
	n, _ := result.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r SuggestionRepository) EventTx(ctx context.Context, tx *sql.Tx, id, from, to, actor, note string, now time.Time) error {
	_, e := tx.ExecContext(ctx, `INSERT INTO suggestion_events(id,suggestion_id,from_status,to_status,actor_id,note,created_at) VALUES(?,?,?,?,?,?,?)`, platform.ID(), id, from, to, actor, note, now.Format(time.RFC3339Nano))
	return e
}
func (r SuggestionRepository) ListByScope(ctx context.Context, scope string, limit int) ([]domain.Suggestion, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	rows, e := r.DB.QueryContext(ctx, `SELECT id,campaign_id,COALESCE(intake_point_id,''),author_id,title,body,scope,status,version,due_at,created_at,updated_at FROM suggestions WHERE scope=? ORDER BY updated_at DESC LIMIT ?`, scope, limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.Suggestion
	for rows.Next() {
		var x domain.Suggestion
		var status, created, updated string
		var due sql.NullString
		if e = rows.Scan(&x.ID, &x.CampaignID, &x.IntakePointID, &x.AuthorID, &x.Title, &x.Body, &x.Scope, &status, &x.Version, &due, &created, &updated); e != nil {
			return nil, e
		}
		x.Status = domain.SuggestionStatus(status)
		x.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		x.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
		out = append(out, x)
	}
	return out, rows.Err()
}
