package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type AuditStore struct{ DB *sql.DB }

func (a AuditStore) Append(ctx context.Context, e domain.AuditEntry) error {
	if err := e.Validate(); err != nil {
		return err
	}
	_, err := a.DB.ExecContext(ctx, `INSERT INTO audit_logs(id,actor_id,action,object_type,object_id,result,request_id,metadata,created_at) VALUES(?,?,?,?,?,?,?,?,?)`, platform.ID(), e.ActorID, e.Action, e.ObjectType, e.ObjectID, e.Result, e.RequestID, "{}", time.Now().UTC().Format(time.RFC3339Nano))
	return err
}
func (a AuditStore) Count(ctx context.Context, objectType, objectID string) (int, error) {
	var n int
	e := a.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_logs WHERE object_type=? AND object_id=?`, objectType, objectID).Scan(&n)
	return n, e
}
