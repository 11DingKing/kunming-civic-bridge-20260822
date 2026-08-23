package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type AssignmentStore struct{ DB *sql.DB }

func (a AssignmentStore) Create(ctx context.Context, suggestion, department string, due time.Time) error {
	if _, e := a.DB.ExecContext(ctx, `INSERT INTO assignments(id,suggestion_id,department_id,status,version,assigned_at) VALUES(?,?,?,'assigned',0,?)`, platform.ID(), suggestion, department, time.Now().UTC().Format(time.RFC3339Nano)); e != nil {
		return e
	}
	_ = due
	return nil
}
func (a AssignmentStore) Current(ctx context.Context, suggestion string) (domain.AssignmentRules, error) {
	var r domain.AssignmentRules
	var status, assigned string
	var version int
	e := a.DB.QueryRowContext(ctx, `SELECT status,COALESCE(assignee_id,''),version FROM assignments WHERE suggestion_id=? ORDER BY assigned_at DESC LIMIT 1`, suggestion).Scan(&status, &assigned, &version)
	if e != nil {
		return r, e
	}
	r.Status = domain.AssignmentStatus(status)
	r.AssigneeID = assigned
	r.Version = version
	return r, nil
}
func (a AssignmentStore) ReleaseExpired(ctx context.Context, now time.Time) error {
	_, e := a.DB.ExecContext(ctx, `UPDATE assignments SET status='assigned',assignee_id=NULL,lease_token=NULL,lease_until=NULL,version=version+1 WHERE status='claimed' AND lease_until<=?`, now.Format(time.RFC3339Nano))
	return e
}
