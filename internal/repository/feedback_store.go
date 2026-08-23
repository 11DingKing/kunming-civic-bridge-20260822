package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type FeedbackStore struct{ DB *sql.DB }

func (f FeedbackStore) Create(ctx context.Context, suggestion, author, comment string, rating int) error {
	if e := (domain.FeedbackRules{Rating: rating, Comment: comment, CreatedAt: time.Now().UTC()}).Validate(); e != nil {
		return e
	}
	_, e := f.DB.ExecContext(ctx, `INSERT INTO feedback(id,suggestion_id,author_id,rating,comment,created_at) VALUES(?,?,?,?,?,?)`, platform.ID(), suggestion, author, rating, comment, time.Now().UTC().Format(time.RFC3339Nano))
	return e
}
func (f FeedbackStore) Average(ctx context.Context, suggestion string) (float64, error) {
	var v sql.NullFloat64
	e := f.DB.QueryRowContext(ctx, `SELECT AVG(rating) FROM feedback WHERE suggestion_id=?`, suggestion).Scan(&v)
	if e != nil {
		return 0, e
	}
	if !v.Valid {
		return 0, nil
	}
	return v.Float64, nil
}
