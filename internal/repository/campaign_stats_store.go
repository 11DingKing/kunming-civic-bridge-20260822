package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
)

type CampaignStatsStore struct{ DB *sql.DB }
type CampaignStats struct{ Total, Submitted, Triaged, Assigned, InProgress, Responded, Closed, Rejected int }

func (s CampaignStatsStore) Get(ctx context.Context, campaign string) (CampaignStats, error) {
	var out CampaignStats
	rows, e := s.DB.QueryContext(ctx, `SELECT status,COUNT(*) FROM suggestions WHERE campaign_id=? GROUP BY status`, campaign)
	if e != nil {
		return out, e
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var n int
		if e = rows.Scan(&status, &n); e != nil {
			return out, e
		}
		out.Total += n
		switch domain.SuggestionStatus(status) {
		case domain.StatusSubmitted:
			out.Submitted = n
		case domain.StatusTriaged:
			out.Triaged = n
		case domain.StatusAssigned:
			out.Assigned = n
		case domain.StatusInProgress:
			out.InProgress = n
		case domain.StatusResponded:
			out.Responded = n
		case domain.StatusClosed:
			out.Closed = n
		case domain.StatusRejected:
			out.Rejected = n
		}
	}
	return out, rows.Err()
}
func (s CampaignStatsStore) CompletionRate(ctx context.Context, campaign string) (float64, error) {
	x, e := s.Get(ctx, campaign)
	if e != nil || x.Total == 0 {
		return 0, e
	}
	return float64(x.Closed) / float64(x.Total), nil
}
