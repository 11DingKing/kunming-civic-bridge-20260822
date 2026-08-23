package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/audit"
	"github.com/11DingKing/kunming-civic-bridge/internal/clock"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/notify"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"strings"
	"time"
)

type Service struct {
	DB     *sql.DB
	Store  repository.Store
	Audit  audit.Logger
	Notify notify.Queue
	Clock  clock.Clock
}

func (s Service) Submit(ctx context.Context, author, campaign, point, title, body, scope, request string) (domain.Suggestion, error) {
	if e := domain.ValidateSuggestion(title, body, scope); e != nil {
		return domain.Suggestion{}, e
	}
	now := s.Clock.Now()
	x := domain.Suggestion{ID: platform.ID(), AuthorID: author, CampaignID: campaign, IntakePointID: point, Title: title, Body: body, Scope: scope, Status: domain.StatusDraft, CreatedAt: now, UpdatedAt: now}
	tx, e := s.Store.Tx(ctx)
	if e != nil {
		return x, e
	}
	defer tx.Rollback()
	if e = s.Store.CreateSuggestion(ctx, tx, x); e != nil {
		return x, e
	}
	if e = s.Store.AddEvent(ctx, tx, x.ID, "", string(domain.StatusDraft), author, "created"); e != nil {
		return x, e
	}
	if e = s.Audit.Record(ctx, tx, author, "suggestion.create", "suggestion", x.ID, "ok", request, map[string]any{"scope": scope}); e != nil {
		return x, e
	}
	if e = tx.Commit(); e != nil {
		return x, e
	}
	return x, nil
}
func (s Service) Send(ctx context.Context, id, actor, request string) (domain.Suggestion, error) {
	x, e := s.Store.GetSuggestion(ctx, id)
	if e != nil {
		return x, e
	}
	tx, e := s.Store.Tx(ctx)
	if e != nil {
		return x, e
	}
	defer tx.Rollback()
	if e = s.Store.Transition(ctx, tx, id, x.Status, domain.StatusSubmitted, x.Version); e != nil {
		return x, e
	}
	if e = s.Store.AddEvent(ctx, tx, id, string(x.Status), string(domain.StatusSubmitted), actor, "submitted"); e != nil {
		return x, e
	}
	if e = s.Audit.Record(ctx, tx, actor, "suggestion.submit", "suggestion", id, "ok", request, nil); e != nil {
		return x, e
	}
	if e = tx.Commit(); e != nil {
		return x, e
	}
	x.Status = domain.StatusSubmitted
	x.Version++
	return x, nil
}
func (s Service) Review(ctx context.Context, id, actor, decision, note, request string) (domain.Suggestion, error) {
	x, e := s.Store.GetSuggestion(ctx, id)
	if e != nil {
		return x, e
	}
	if decision != "accept" && decision != "reject" {
		return x, fmt.Errorf("%w: decision", domain.ErrInvalid)
	}
	to := domain.StatusTriaged
	if decision == "reject" {
		to = domain.StatusRejected
	}
	tx, e := s.Store.Tx(ctx)
	if e != nil {
		return x, e
	}
	defer tx.Rollback()
	if e = s.Store.Transition(ctx, tx, id, x.Status, to, x.Version); e != nil {
		return x, e
	}
	_, e = tx.ExecContext(ctx, `INSERT INTO reviews(id,suggestion_id,reviewer_id,decision,note,version,created_at) VALUES(?,?,?,?,?,?,?)`, platform.ID(), id, actor, decision, note, x.Version+1, time.Now().UTC().Format(time.RFC3339Nano))
	if e != nil {
		return x, e
	}
	if e = s.Store.AddEvent(ctx, tx, id, string(x.Status), string(to), actor, note); e != nil {
		return x, e
	}
	if e = s.Audit.Record(ctx, tx, actor, "suggestion.review", "suggestion", id, "ok", request, map[string]any{"decision": decision}); e != nil {
		return x, e
	}
	if e = tx.Commit(); e != nil {
		return x, e
	}
	x.Status = to
	x.Version++
	return x, nil
}

func (s Service) Assign(ctx context.Context, id, actor, department, request string) error {
	x, err := s.Store.GetSuggestion(ctx, id)
	if err != nil {
		return err
	}
	tx, err := s.Store.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = s.Store.Transition(ctx, tx, id, x.Status, domain.StatusAssigned, x.Version); err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err = tx.ExecContext(ctx, `INSERT INTO assignments(id,suggestion_id,department_id,status,version,assigned_at) VALUES(?,?,?,'assigned',0,?)`, platform.ID(), id, department, now)
	if err != nil {
		return err
	}
	if err = s.Store.AddEvent(ctx, tx, id, string(x.Status), string(domain.StatusAssigned), actor, department); err != nil {
		return err
	}
	if err = s.Notify.Enqueue(ctx, tx, actor, id, "assignment", department); err != nil {
		return err
	}
	if err = s.Audit.Record(ctx, tx, actor, "suggestion.assign", "suggestion", id, "ok", request, map[string]any{"department": department}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s Service) StartWork(ctx context.Context, id, actor, request string) error {
	x, err := s.Store.GetSuggestion(ctx, id)
	if err != nil {
		return err
	}
	tx, err := s.Store.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = s.Store.Transition(ctx, tx, id, x.Status, domain.StatusInProgress, x.Version); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE assignments SET status='in_progress' WHERE suggestion_id=? AND status='assigned'`, id); err != nil {
		return err
	}
	if err = s.Store.AddEvent(ctx, tx, id, string(x.Status), string(domain.StatusInProgress), actor, "work started"); err != nil {
		return err
	}
	if err = s.Audit.Record(ctx, tx, actor, "suggestion.start", "suggestion", id, "ok", request, nil); err != nil {
		return err
	}
	return tx.Commit()
}

func (s Service) DraftResponse(ctx context.Context, id, actor, body, request string) (domain.Response, error) {
	if len(body) < 10 {
		return domain.Response{}, fmt.Errorf("%w: response body", domain.ErrInvalid)
	}
	x, err := s.Store.GetSuggestion(ctx, id)
	if err != nil {
		return domain.Response{}, err
	}
	if x.Status != domain.StatusInProgress && x.Status != domain.StatusReopened {
		return domain.Response{}, fmt.Errorf("%w: response state", domain.ErrInvalid)
	}
	r := domain.Response{ID: platform.ID(), SuggestionID: id, AuthorID: actor, Body: body, Status: "draft", CreatedAt: s.Clock.Now()}
	tx, err := s.Store.Tx(ctx)
	if err != nil {
		return r, err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO responses(id,suggestion_id,author_id,body,status,created_at) VALUES(?,?,?,?,?,?)`, r.ID, r.SuggestionID, r.AuthorID, r.Body, r.Status, r.CreatedAt.Format(time.RFC3339Nano))
	if err != nil {
		return r, err
	}
	if err = s.Audit.Record(ctx, tx, actor, "response.draft", "response", r.ID, "ok", request, nil); err != nil {
		return r, err
	}
	if err = tx.Commit(); err != nil {
		return r, err
	}
	return r, nil
}

func (s Service) PublishResponse(ctx context.Context, responseID, actor, request string) error {
	tx, err := s.Store.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var suggestionID, status string
	if err = tx.QueryRowContext(ctx, `SELECT suggestion_id,status FROM responses WHERE id=?`, responseID).Scan(&suggestionID, &status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}
	if status != "approved" {
		return fmt.Errorf("%w: response is not approved", domain.ErrInvalid)
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err = tx.ExecContext(ctx, `UPDATE responses SET status='published',published_at=? WHERE id=? AND status='approved'`, now, responseID); err != nil {
		return err
	}
	if err = s.Audit.Record(ctx, tx, actor, "response.publish", "response", responseID, "ok", request, map[string]any{"suggestion_id": suggestionID}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s Service) AddFeedback(ctx context.Context, id, author string, rating int, comment, request string) error {
	if err := domain.ValidateRating(rating); err != nil {
		return err
	}
	if strings.TrimSpace(comment) == "" {
		return fmt.Errorf("%w: comment", domain.ErrInvalid)
	}
	tx, err := s.Store.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var published int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM responses WHERE suggestion_id=? AND status='published'`, id).Scan(&published); err != nil {
		return err
	}
	if published == 0 {
		return fmt.Errorf("%w: response not published", domain.ErrInvalid)
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO feedback(id,suggestion_id,author_id,rating,comment,created_at) VALUES(?,?,?,?,?,?)`, platform.ID(), id, author, rating, comment, time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return err
	}
	if err = s.Audit.Record(ctx, tx, author, "feedback.create", "suggestion", id, "ok", request, map[string]any{"rating": rating}); err != nil {
		return err
	}
	return tx.Commit()
}
