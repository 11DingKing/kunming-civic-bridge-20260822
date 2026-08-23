package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"time"
)

type ResponseHistoryStore struct{ DB *sql.DB }
type ResponseHistory struct {
	ID, ResponseID, ActorID, FromStatus, ToStatus, Note string
	CreatedAt                                           time.Time
}

func (s ResponseHistoryStore) Append(ctx context.Context, response, actor, from, to, note string) error {
	if !domain.ResponseTransition(domain.ResponseStatus(from), domain.ResponseStatus(to)) && from != to {
		return domain.ErrInvalid
	}
	_, e := s.DB.ExecContext(ctx, `INSERT INTO suggestion_events(id,suggestion_id,from_status,to_status,actor_id,note,created_at) VALUES(lower(hex(randomblob(16))),?,?,?,?,?,?)`, response, from, to, actor, note, time.Now().UTC().Format(time.RFC3339Nano))
	return e
}
func (s ResponseHistoryStore) List(ctx context.Context, response string) ([]ResponseHistory, error) {
	rows, e := s.DB.QueryContext(ctx, `SELECT id,suggestion_id,actor_id,COALESCE(from_status,''),to_status,COALESCE(note,''),created_at FROM suggestion_events WHERE suggestion_id=? ORDER BY created_at`, response)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []ResponseHistory
	for rows.Next() {
		var x ResponseHistory
		var created string
		if e = rows.Scan(&x.ID, &x.ResponseID, &x.ActorID, &x.FromStatus, &x.ToStatus, &x.Note, &created); e != nil {
			return nil, e
		}
		x.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		out = append(out, x)
	}
	return out, rows.Err()
}
