package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"strings"
)

type SuggestionQueryStore struct{ DB *sql.DB }
type SuggestionFilter struct {
	Scope, Status, Keyword string
	Limit, Offset          int
	IncludeClosed          bool
}

func (q SuggestionQueryStore) Search(ctx context.Context, f SuggestionFilter) ([]domain.Suggestion, error) {
	if f.Limit < 1 {
		f.Limit = 20
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	if f.Status == string(domain.StatusClosed) && !f.IncludeClosed {
		return nil, fmt.Errorf("%w: closed filter", domain.ErrForbidden)
	}
	where := []string{"1=1"}
	args := []any{}
	if f.Scope != "" {
		where = append(where, "scope=?")
		args = append(args, f.Scope)
	}
	if f.Status != "" {
		where = append(where, "status=?")
		args = append(args, f.Status)
	}
	if strings.TrimSpace(f.Keyword) != "" {
		where = append(where, "(title LIKE ? OR body LIKE ?)")
		kw := "%" + strings.TrimSpace(f.Keyword) + "%"
		args = append(args, kw, kw)
	}
	args = append(args, f.Limit, f.Offset)
	rows, e := q.DB.QueryContext(ctx, `SELECT id,campaign_id,COALESCE(intake_point_id,''),author_id,title,body,scope,status,version,due_at,created_at,updated_at FROM suggestions WHERE `+strings.Join(where, " AND ")+` ORDER BY updated_at DESC LIMIT ? OFFSET ?`, args...)
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
		out = append(out, x)
	}
	return out, rows.Err()
}
func (q SuggestionQueryStore) Count(ctx context.Context, f SuggestionFilter) (int, error) {
	var n int
	e := q.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestions WHERE (?='' OR scope=?) AND (?='' OR status=?)`, f.Scope, f.Scope, f.Status, f.Status).Scan(&n)
	return n, e
}
