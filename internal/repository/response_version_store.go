package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type ResponseVersionStore struct{ DB *sql.DB }
type ResponseVersion struct {
	ID, ResponseID, AuthorID, Body string
	Version                        int
	CreatedAt                      time.Time
}

func (s ResponseVersionStore) Add(ctx context.Context, response, author, body string, version int) (string, error) {
	if e := domain.ValidateResponse(body); e != nil {
		return "", e
	}
	id := platform.ID()
	_, e := s.DB.ExecContext(ctx, `INSERT INTO suggestion_events(id,suggestion_id,from_status,to_status,actor_id,note,created_at) VALUES(?,?,?,?,?,?,?)`, id, response, "", fmt.Sprintf("response_version_%d", version), author, body, time.Now().UTC().Format(time.RFC3339Nano))
	return id, e
}
func (s ResponseVersionStore) List(ctx context.Context, response string) ([]ResponseVersion, error) {
	rows, e := s.DB.QueryContext(ctx, `SELECT id,suggestion_id,actor_id,to_status,note,created_at FROM suggestion_events WHERE suggestion_id=? AND to_status LIKE 'response_version_%' ORDER BY created_at`, response)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []ResponseVersion
	for rows.Next() {
		var x ResponseVersion
		var status, created string
		if e = rows.Scan(&x.ID, &x.ResponseID, &x.AuthorID, &status, &x.Body, &created); e != nil {
			return nil, e
		}
		fmt.Sscanf(status, "response_version_%d", &x.Version)
		x.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		out = append(out, x)
	}
	return out, rows.Err()
}
