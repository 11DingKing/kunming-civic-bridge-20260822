package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"time"
)

type EscalationService struct{ DB *sql.DB }

func (e EscalationService) FindDue(ctx context.Context, now time.Time) ([]domain.Escalation, error) {
	rows, err := e.DB.QueryContext(ctx, `SELECT id,due_at FROM suggestions WHERE due_at IS NOT NULL AND due_at<=? AND status IN ('assigned','in_progress')`, now.Format(time.RFC3339Nano))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Escalation
	for rows.Next() {
		var id, due string
		if err = rows.Scan(&id, &due); err != nil {
			return nil, err
		}
		at, _ := time.Parse(time.RFC3339Nano, due)
		out = append(out, domain.Escalation{SuggestionID: id, DueAt: at, Level: 1, Reason: "deadline"})
	}
	return out, rows.Err()
}
func (e EscalationService) Mark(ctx context.Context, id string) error {
	_, err := e.DB.ExecContext(ctx, `UPDATE suggestions SET status='escalated',version=version+1,updated_at=? WHERE id=? AND status IN ('assigned','in_progress')`, time.Now().UTC().Format(time.RFC3339Nano), id)
	return err
}
