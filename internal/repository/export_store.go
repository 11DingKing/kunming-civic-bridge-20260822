package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"time"
)

type ExportStore struct{ DB *sql.DB }

func (s ExportStore) Audit(ctx context.Context, req domain.AuditExport) ([]AuditRecord, error) {
	if e := req.Validate(); e != nil {
		return nil, e
	}
	records, e := (&AuditQueryStore{DB: s.DB}).ForObject(ctx, req.ObjectType, req.ObjectID)
	if e != nil {
		return nil, e
	}
	out := records[:0]
	for _, record := range records {
		if req.Includes(record.CreatedAt) {
			out = append(out, record)
		}
	}
	return out, nil
}
func (s ExportStore) HasActivity(ctx context.Context, id string, from, to time.Time) (bool, error) {
	var n int
	e := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_logs WHERE object_id=? AND created_at>=? AND created_at<?`, id, from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano)).Scan(&n)
	return n > 0, e
}
