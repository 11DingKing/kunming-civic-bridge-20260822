package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"time"
)

type AuditQueryStore struct{ DB *sql.DB }
type AuditRecord struct {
	ID, ActorID, Action, ObjectType, ObjectID, Result, RequestID, Metadata string
	CreatedAt                                                              time.Time
}

func (a AuditQueryStore) ForObject(ctx context.Context, typ, id string) ([]AuditRecord, error) {
	rows, e := a.DB.QueryContext(ctx, `SELECT id,COALESCE(actor_id,''),action,object_type,object_id,result,request_id,metadata,created_at FROM audit_logs WHERE object_type=? AND object_id=? ORDER BY created_at`, typ, id)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []AuditRecord
	for rows.Next() {
		var x AuditRecord
		var created string
		if e = rows.Scan(&x.ID, &x.ActorID, &x.Action, &x.ObjectType, &x.ObjectID, &x.Result, &x.RequestID, &x.Metadata, &created); e != nil {
			return nil, e
		}
		x.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		out = append(out, x)
	}
	return out, rows.Err()
}
func (a AuditQueryStore) Append(ctx context.Context, e domain.AuditEntry) error {
	if err := e.Validate(); err != nil {
		return err
	}
	_, err := a.DB.ExecContext(ctx, `INSERT INTO audit_logs(id,actor_id,action,object_type,object_id,result,request_id,metadata,created_at) VALUES(lower(hex(randomblob(16))),?,?,?,?,?,?,?,?)`, e.ActorID, e.Action, e.ObjectType, e.ObjectID, e.Result, e.RequestID, "{}", time.Now().UTC().Format(time.RFC3339Nano))
	return err
}
