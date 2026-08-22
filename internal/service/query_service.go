package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
)

type QueryService struct{ DB *sql.DB }
type SuggestionSummary struct {
	ID, Title, Status, Scope string
	Version                  int
}

func (q QueryService) Search(ctx context.Context, in domain.SuggestionQuery) ([]SuggestionSummary, error) {
	in = in.Normalize()
	if e := in.Validate(); e != nil {
		return nil, e
	}
	rows, e := q.DB.QueryContext(ctx, `SELECT id,title,status,scope,version FROM suggestions WHERE (?='' OR scope=?) AND (?='' OR status=?) AND (?='' OR title LIKE '%'||?||'%' OR body LIKE '%'||?||'%') ORDER BY updated_at DESC LIMIT ? OFFSET ?`, in.Scope, in.Scope, in.Status, in.Status, in.Keyword, in.Keyword, in.Keyword, in.Limit, in.Offset)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []SuggestionSummary
	for rows.Next() {
		var x SuggestionSummary
		if e = rows.Scan(&x.ID, &x.Title, &x.Status, &x.Scope, &x.Version); e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (q QueryService) Count(ctx context.Context, in domain.SuggestionQuery) (int, error) {
	in = in.Normalize()
	var n int
	e := q.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestions WHERE (?='' OR scope=?) AND (?='' OR status=?) AND (?='' OR title LIKE '%'||?||'%' OR body LIKE '%'||?||'%')`, in.Scope, in.Scope, in.Status, in.Status, in.Keyword, in.Keyword, in.Keyword).Scan(&n)
	return n, e
}
