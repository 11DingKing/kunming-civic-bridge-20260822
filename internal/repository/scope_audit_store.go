package repository

import (
	"context"
	"database/sql"
)

type ScopeAuditStore struct{ DB *sql.DB }

func (s ScopeAuditStore) CountByScope(ctx context.Context) (map[string]int, error) {
	rows, e := s.DB.QueryContext(ctx, `SELECT scope,COUNT(*) FROM suggestions GROUP BY scope ORDER BY scope`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var scope string
		var n int
		if e = rows.Scan(&scope, &n); e != nil {
			return nil, e
		}
		out[scope] = n
	}
	return out, rows.Err()
}
