package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"time"
)

type IntakeService struct {
	DB        *sql.DB
	Campaigns repository.CampaignQueryStore
	Points    repository.IntakePointStore
	Now       func() time.Time
}

func (i IntakeService) Submit(ctx context.Context, d domain.SuggestionDraft, request string) (domain.Suggestion, error) {
	if e := d.Validate(); e != nil {
		return domain.Suggestion{}, e
	}
	window, e := i.Campaigns.Get(ctx, d.CampaignID)
	if e != nil {
		return domain.Suggestion{}, e
	}
	rules := domain.CampaignRules{Status: domain.CampaignStatus(window.Status), StartsAt: window.StartsAt, EndsAt: window.EndsAt, Timezone: "Asia/Shanghai", MaxSubmissions: 10000}
	if e = rules.CanAccept(i.Now()); e != nil {
		return domain.Suggestion{}, e
	}
	point, e := i.Points.Get(ctx, d.IntakePointID)
	if e != nil {
		return domain.Suggestion{}, e
	}
	if !point.Active {
		return domain.Suggestion{}, fmt.Errorf("%w: intake point inactive", domain.ErrConflict)
	}
	now := i.Now()
	x := domain.Suggestion{ID: platform.ID(), CampaignID: d.CampaignID, IntakePointID: d.IntakePointID, AuthorID: d.AuthorID, Title: d.Title, Body: d.Body, Scope: d.Scope, Status: domain.StatusDraft, Version: 0, DueAt: d.DueAt, CreatedAt: now, UpdatedAt: now}
	tx, e := i.DB.BeginTx(ctx, nil)
	if e != nil {
		return x, e
	}
	defer tx.Rollback()
	_, e = tx.ExecContext(ctx, `INSERT INTO suggestions(id,campaign_id,intake_point_id,author_id,title,body,scope,status,version,due_at,created_at,updated_at) VALUES(?,?,?,?,?,?,?,'draft',0,?,?,?)`, x.ID, x.CampaignID, x.IntakePointID, x.AuthorID, x.Title, x.Body, x.Scope, x.DueAt, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if e != nil {
		return x, e
	}
	if e = tx.Commit(); e != nil {
		return x, e
	}
	_, e = i.DB.ExecContext(ctx, `INSERT INTO suggestion_events(id,suggestion_id,from_status,to_status,actor_id,note,created_at) VALUES(?,?,?,?,?,?,?)`, platform.ID(), x.ID, "", string(domain.StatusDraft), x.AuthorID, request, now.Format(time.RFC3339Nano))
	if e != nil {
		return x, e
	}
	return x, nil
}
