package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type CampaignQueryStore struct{ DB *sql.DB }
type CampaignRecord struct {
	ID, Name, Status            string
	StartsAt, EndsAt, CreatedAt time.Time
}

func (s CampaignQueryStore) Get(ctx context.Context, id string) (CampaignRecord, error) {
	var c CampaignRecord
	var starts, ends, created string
	e := s.DB.QueryRowContext(ctx, `SELECT id,name,status,starts_at,ends_at,created_at FROM campaigns WHERE id=?`, id).Scan(&c.ID, &c.Name, &c.Status, &starts, &ends, &created)
	if e != nil {
		return c, e
	}
	c.StartsAt, _ = time.Parse(time.RFC3339Nano, starts)
	c.EndsAt, _ = time.Parse(time.RFC3339Nano, ends)
	c.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return c, nil
}
func (s CampaignQueryStore) ListOpen(ctx context.Context, now time.Time, limit int) ([]CampaignRecord, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	rows, e := s.DB.QueryContext(ctx, `SELECT id,name,status,starts_at,ends_at,created_at FROM campaigns WHERE status='open' AND starts_at<=? AND ends_at>? ORDER BY starts_at LIMIT ?`, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []CampaignRecord
	for rows.Next() {
		var c CampaignRecord
		var starts, ends, created string
		if e = rows.Scan(&c.ID, &c.Name, &c.Status, &starts, &ends, &created); e != nil {
			return nil, e
		}
		c.StartsAt, _ = time.Parse(time.RFC3339Nano, starts)
		c.EndsAt, _ = time.Parse(time.RFC3339Nano, ends)
		c.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		out = append(out, c)
	}
	return out, rows.Err()
}
func (s CampaignQueryStore) Open(ctx context.Context, id string, now time.Time) error {
	r, e := s.Get(ctx, id)
	if e != nil {
		return e
	}
	rules := domain.CampaignRules{Status: domain.CampaignStatus(r.Status), StartsAt: r.StartsAt, EndsAt: r.EndsAt, Timezone: "Asia/Shanghai", MaxSubmissions: 10000}
	if e = rules.CanOpen(now); e != nil {
		return e
	}
	res, e := s.DB.ExecContext(ctx, `UPDATE campaigns SET status='open' WHERE id=? AND status IN ('draft','paused')`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (s CampaignQueryStore) New(ctx context.Context, name string, start, end time.Time) (string, error) {
	if name == "" || !start.Before(end) {
		return "", domain.ErrInvalid
	}
	id := platform.ID()
	_, e := s.DB.ExecContext(ctx, `INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES(?,?,?,?,?,?)`, id, name, start.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano), domain.CampaignDraft, time.Now().UTC().Format(time.RFC3339Nano))
	return id, e
}
