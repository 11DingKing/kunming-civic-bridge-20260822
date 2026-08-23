package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type ReopenStore struct{ DB *sql.DB }

func (s ReopenStore) Request(ctx context.Context, r domain.ReopenRequest) error {
	if e := r.Validate(); e != nil {
		return e
	}
	_, e := s.DB.ExecContext(ctx, `INSERT INTO feedback(id,suggestion_id,author_id,rating,comment,created_at) VALUES(?,?,?,?,?,?)`, platform.ID(), r.SuggestionID, r.AuthorID, 1, r.Reason, r.RequestedAt.Format(time.RFC3339Nano))
	return e
}
func (s ReopenStore) Latest(ctx context.Context, id string) (domain.ReopenRequest, error) {
	var author, reason, created string
	e := s.DB.QueryRowContext(ctx, `SELECT author_id,comment,created_at FROM feedback WHERE suggestion_id=? ORDER BY created_at DESC LIMIT 1`, id).Scan(&author, &reason, &created)
	if e != nil {
		return domain.ReopenRequest{}, e
	}
	at, _ := time.Parse(time.RFC3339Nano, created)
	return domain.ReopenRequest{SuggestionID: id, AuthorID: author, Reason: reason, RequestedAt: at}, nil
}
