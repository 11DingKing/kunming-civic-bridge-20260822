package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type Store struct{ DB *sql.DB }

func (s Store) Tx(ctx context.Context) (*sql.Tx, error) { return s.DB.BeginTx(ctx, nil) }
func (s Store) CreateSuggestion(ctx context.Context, tx *sql.Tx, x domain.Suggestion) error {
	_, e := tx.ExecContext(ctx, `INSERT INTO suggestions(id,campaign_id,intake_point_id,author_id,title,body,scope,status,version,due_at,created_at,updated_at) VALUES(?,?,?,?,?,?,?,'draft',0,?,?,?)`, x.ID, x.CampaignID, x.IntakePointID, x.AuthorID, x.Title, x.Body, x.Scope, x.DueAt, x.CreatedAt.Format(time.RFC3339Nano), x.UpdatedAt.Format(time.RFC3339Nano))
	return e
}
func (s Store) GetSuggestion(ctx context.Context, id string) (domain.Suggestion, error) {
	var x domain.Suggestion
	var st string
	var due, created, updated sql.NullString
	err := s.DB.QueryRowContext(ctx, `SELECT id,campaign_id,COALESCE(intake_point_id,''),author_id,title,body,scope,status,version,due_at,created_at,updated_at FROM suggestions WHERE id=?`, id).Scan(&x.ID, &x.CampaignID, &x.IntakePointID, &x.AuthorID, &x.Title, &x.Body, &x.Scope, &st, &x.Version, &due, &created, &updated)
	if errors.Is(err, sql.ErrNoRows) {
		return x, domain.ErrNotFound
	}
	if err != nil {
		return x, err
	}
	x.Status = domain.SuggestionStatus(st)
	if created.Valid {
		x.CreatedAt, _ = time.Parse(time.RFC3339Nano, created.String)
	}
	if updated.Valid {
		x.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated.String)
	}
	if due.Valid {
		parsed, _ := time.Parse(time.RFC3339Nano, due.String)
		x.DueAt = &parsed
	}
	return x, nil
}
func (s Store) Transition(ctx context.Context, tx *sql.Tx, id string, from, to domain.SuggestionStatus, version int) error {
	if !from.CanTransition(to) {
		return fmt.Errorf("%w: %s to %s", domain.ErrInvalid, from, to)
	}
	r, e := tx.ExecContext(ctx, `UPDATE suggestions SET status=?,version=version+1,updated_at=? WHERE id=? AND status=? AND version=?`, to, time.Now().UTC().Format(time.RFC3339Nano), id, from, version)
	if e != nil {
		return e
	}
	n, _ := r.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (s Store) AddEvent(ctx context.Context, tx *sql.Tx, suggestion, from, to, actor, note string) error {
	_, e := tx.ExecContext(ctx, `INSERT INTO suggestion_events(id,suggestion_id,from_status,to_status,actor_id,note,created_at) VALUES(?,?,?,?,?,?,?)`, platform.ID(), suggestion, from, to, actor, note, time.Now().UTC().Format(time.RFC3339Nano))
	return e
}
func (s Store) CreateJob(ctx context.Context, tx *sql.Tx, kind, payload string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, e := tx.ExecContext(ctx, `INSERT INTO jobs(id,kind,payload,status,attempts,run_after,created_at,updated_at) VALUES(?,?,?,'pending',0,?,?,?)`, platform.ID(), kind, payload, now, now, now)
	return e
}
