package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type Logger struct{ DB *sql.DB }

func (l Logger) Record(ctx context.Context, tx *sql.Tx, actor, action, typ, object, result, request string, metadata map[string]any) error {
	b, e := json.Marshal(metadata)
	if e != nil {
		return e
	}
	_, e = tx.ExecContext(ctx, `INSERT INTO audit_logs(id,actor_id,action,object_type,object_id,result,request_id,metadata,created_at) VALUES(?,?,?,?,?,?,?,?,?)`, platform.ID(), actor, action, typ, object, result, request, string(b), time.Now().UTC().Format(time.RFC3339Nano))
	return e
}
