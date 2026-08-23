package service

import (
	"context"
	"testing"
	"time"

	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

func TestBaselineWorkflowReportsPreserveLifecycleData(t *testing.T) {
	db := workflowSeed(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stats := CampaignStatsService{DB: db, Store: repository.CampaignStatsStore{DB: db}}
	got, err := stats.Get(ctx, "campaign-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Total < got.Submitted+got.Triaged+got.Assigned+got.InProgress+got.Responded+got.Closed+got.Rejected {
		t.Fatalf("inconsistent status totals: %+v", got)
	}
	rate, err := stats.CompletionRate(ctx, "campaign-1")
	if err != nil {
		t.Fatal(err)
	}
	if rate < 0 || rate > 1 {
		t.Fatalf("completion rate out of range: %v", rate)
	}
}

func TestBaselineDomainTransitionsRejectInvalidLifecycleEdges(t *testing.T) {
	cases := []struct {
		from domain.SuggestionStatus
		to   domain.SuggestionStatus
	}{
		{domain.StatusDraft, domain.StatusSubmitted},
		{domain.StatusSubmitted, domain.StatusTriaged},
		{domain.StatusTriaged, domain.StatusAssigned},
		{domain.StatusAssigned, domain.StatusInProgress},
		{domain.StatusInProgress, domain.StatusResponded},
		{domain.StatusResponded, domain.StatusClosed},
	}
	for _, tc := range cases {
		if !tc.from.CanTransition(tc.to) {
			t.Fatalf("expected transition %s -> %s", tc.from, tc.to)
		}
	}
	if domain.StatusClosed.CanTransition(domain.StatusDraft) {
		t.Fatal("closed suggestions must not return to draft")
	}
}
