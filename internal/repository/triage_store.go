package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"strings"
	"time"
)

type TriageStore struct{ DB *sql.DB }

func (s TriageStore) Save(ctx context.Context, suggestion, actor string, result domain.TriageResult) error {
	if e := result.Validate(); e != nil {
		return e
	}
	_, e := s.DB.ExecContext(ctx, `INSERT INTO reviews(id,suggestion_id,reviewer_id,decision,note,version,created_at) VALUES(?,?,?,?,?,?,?)`, platform.ID(), suggestion, actor, "triage", result.Category+":"+result.Priority, 0, time.Now().UTC().Format(time.RFC3339Nano))
	return e
}
func (s TriageStore) Latest(ctx context.Context, suggestion string) (domain.TriageResult, error) {
	var note string
	e := s.DB.QueryRowContext(ctx, `SELECT note FROM reviews WHERE suggestion_id=? AND decision='triage' ORDER BY created_at DESC LIMIT 1`, suggestion).Scan(&note)
	if e != nil {
		return domain.TriageResult{}, e
	}
	parts := strings.SplitN(note, ":", 2)
	if len(parts) != 2 {
		return domain.TriageResult{}, domain.ErrInvalid
	}
	return domain.TriageResult{Category: parts[0], Priority: parts[1], Reason: "stored"}, nil
}
