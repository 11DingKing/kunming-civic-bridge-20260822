package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"time"
)

type LeaseAuditStore struct{ DB *sql.DB }

func (s LeaseAuditStore) Record(ctx context.Context, event domain.AssignmentEvent) error {
	if e := event.Validate(); e != nil {
		return e
	}
	_, e := s.DB.ExecContext(ctx, `INSERT INTO suggestion_events(id,suggestion_id,from_status,to_status,actor_id,note,created_at) VALUES(lower(hex(randomblob(16))),?,?,?,?,?,?)`, event.AssignmentID, event.From, event.To, event.ActorID, event.Reason, event.At.Format(time.RFC3339Nano))
	return e
}
func (s LeaseAuditStore) History(ctx context.Context, assignment string) (int, error) {
	var n int
	e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM suggestion_events WHERE suggestion_id=?`, assignment).Scan(&n)
	return n, e
}
