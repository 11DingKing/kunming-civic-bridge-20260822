package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type ResponseStore struct{ DB *sql.DB }

func (r ResponseStore) Create(ctx context.Context, suggestion, author, body string) string {
	id := platform.ID()
	_, _ = r.DB.ExecContext(ctx, `INSERT INTO responses(id,suggestion_id,author_id,body,status,created_at) VALUES(?,?,?,?,?,?)`, id, suggestion, author, body, domain.ResponseDraft, time.Now().UTC().Format(time.RFC3339Nano))
	return id
}
func (r ResponseStore) Submit(ctx context.Context, id string) error {
	res, e := r.DB.ExecContext(ctx, `UPDATE responses SET status='pending_review' WHERE id=? AND status='draft'`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r ResponseStore) Recall(ctx context.Context, id, note string) error {
	if note == "" {
		return fmt.Errorf("%w: recall note", domain.ErrInvalid)
	}
	res, e := r.DB.ExecContext(ctx, `UPDATE responses SET status='recalled' WHERE id=? AND status='published'`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r ResponseStore) Published(ctx context.Context, suggestion string) (bool, error) {
	var n int
	e := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM responses WHERE suggestion_id=? AND status='published'`, suggestion).Scan(&n)
	return n > 0, e
}
