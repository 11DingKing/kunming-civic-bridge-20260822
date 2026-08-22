package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type CampaignStore struct{ DB *sql.DB }

func (c CampaignStore) Create(ctx context.Context, name string, start, end time.Time, scope string) string {
	id := platform.ID()
	_, _ = c.DB.ExecContext(ctx, `INSERT INTO campaigns(id,name,starts_at,ends_at,status,created_at) VALUES(?,?,?,?,?,?)`, id, name, start.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano), domain.CampaignDraft, time.Now().UTC().Format(time.RFC3339Nano))
	return id
}
func (c CampaignStore) SetStatus(ctx context.Context, id string, status domain.CampaignStatus) error {
	_, e := c.DB.ExecContext(ctx, `UPDATE campaigns SET status=? WHERE id=?`, status, id)
	return e
}
func (c CampaignStore) Window(ctx context.Context, id string) (domain.CampaignRules, error) {
	var r domain.CampaignRules
	var start, end, status string
	e := c.DB.QueryRowContext(ctx, `SELECT starts_at,ends_at,status FROM campaigns WHERE id=?`, id).Scan(&start, &end, &status)
	if e != nil {
		return r, e
	}
	r.StartsAt, _ = time.Parse(time.RFC3339Nano, start)
	r.EndsAt, _ = time.Parse(time.RFC3339Nano, end)
	r.Status = domain.CampaignStatus(status)
	r.Timezone = "Asia/Shanghai"
	r.MaxSubmissions = 10000
	return r, nil
}
