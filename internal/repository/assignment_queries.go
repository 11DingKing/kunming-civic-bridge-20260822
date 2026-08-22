package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type AssignmentQueryStore struct{ DB *sql.DB }
type AssignmentRecord struct {
	ID, SuggestionID, DepartmentID, AssigneeID, Status, LeaseToken string
	LeaseUntil                                                     *time.Time
	Version                                                        int
	AssignedAt                                                     time.Time
}

func (a AssignmentQueryStore) Create(ctx context.Context, suggestion, department string, due *time.Time) (string, error) {
	if suggestion == "" || department == "" {
		return "", domain.ErrInvalid
	}
	id := platform.ID()
	now := time.Now().UTC()
	_, e := a.DB.ExecContext(ctx, `INSERT INTO assignments(id,suggestion_id,department_id,status,version,assigned_at) VALUES(?,?,?,'assigned',0,?)`, id, suggestion, department, now.Format(time.RFC3339Nano))
	return id, e
}
func (a AssignmentQueryStore) Get(ctx context.Context, id string) (AssignmentRecord, error) {
	var x AssignmentRecord
	var lease sql.NullString
	var assigned string
	e := a.DB.QueryRowContext(ctx, `SELECT id,suggestion_id,department_id,COALESCE(assignee_id,''),status,COALESCE(lease_token,''),lease_until,version,assigned_at FROM assignments WHERE id=?`, id).Scan(&x.ID, &x.SuggestionID, &x.DepartmentID, &x.AssigneeID, &x.Status, &x.LeaseToken, &lease, &x.Version, &assigned)
	if e != nil {
		return x, e
	}
	if lease.Valid {
		v, _ := time.Parse(time.RFC3339Nano, lease.String)
		x.LeaseUntil = &v
	}
	x.AssignedAt, _ = time.Parse(time.RFC3339Nano, assigned)
	return x, nil
}
func (a AssignmentQueryStore) Claim(ctx context.Context, id, owner, token string, until time.Time, version int) error {
	res, e := a.DB.ExecContext(ctx, `UPDATE assignments SET status='claimed',assignee_id=?,lease_token=?,lease_until=?,version=version+1 WHERE id=? AND status='assigned' AND version=?`, owner, token, until.Format(time.RFC3339Nano), id, version)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (a AssignmentQueryStore) Complete(ctx context.Context, id, owner string, now time.Time) error {
	res, e := a.DB.ExecContext(ctx, `UPDATE assignments SET status='done',completed_at=?,lease_token=NULL,lease_until=NULL,version=version+1 WHERE id=? AND status='claimed' AND assignee_id=? AND lease_until>?`, now.Format(time.RFC3339Nano), id, owner, now.Format(time.RFC3339Nano))
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
