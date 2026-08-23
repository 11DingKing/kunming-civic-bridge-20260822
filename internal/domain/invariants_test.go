package domain

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSuggestionValidationBoundaries(t *testing.T) {
	if err := ValidateSuggestion("", "body", "scope"); !errors.Is(err, ErrInvalid) {
		t.Fatal("empty title must be invalid")
	}
	if err := ValidateSuggestion("title", "", "scope"); !errors.Is(err, ErrInvalid) {
		t.Fatal("empty body must be invalid")
	}
	if err := ValidateSuggestion("title", "body", ""); !errors.Is(err, ErrInvalid) {
		t.Fatal("empty scope must be invalid")
	}
	if err := ValidateSuggestion("title", "body", "district"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateRating(0); !errors.Is(err, ErrInvalid) {
		t.Fatal("zero rating must be invalid")
	}
	if err := ValidateRating(6); !errors.Is(err, ErrInvalid) {
		t.Fatal("six rating must be invalid")
	}
	for i := 1; i <= 5; i++ {
		if err := ValidateRating(i); err != nil {
			t.Fatalf("rating %d: %v", i, err)
		}
	}
}

func TestStatusTransitionMatrix(t *testing.T) {
	valid := [][2]SuggestionStatus{{StatusDraft, StatusSubmitted}, {StatusSubmitted, StatusTriaged}, {StatusSubmitted, StatusRejected}, {StatusTriaged, StatusAssigned}, {StatusAssigned, StatusInProgress}, {StatusInProgress, StatusResponded}, {StatusResponded, StatusClosed}, {StatusResponded, StatusReopened}, {StatusEscalated, StatusAssigned}}
	for _, pair := range valid {
		if !pair[0].CanTransition(pair[1]) {
			t.Fatalf("%s to %s should be valid", pair[0], pair[1])
		}
	}
	invalid := [][2]SuggestionStatus{{StatusDraft, StatusClosed}, {StatusRejected, StatusAssigned}, {StatusClosed, StatusReopened}, {StatusSubmitted, StatusClosed}}
	for _, pair := range invalid {
		if pair[0].CanTransition(pair[1]) {
			t.Fatalf("%s to %s should be invalid", pair[0], pair[1])
		}
	}
}

func TestCampaignBoundaries(t *testing.T) {
	start := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	end := start.Add(14 * 24 * time.Hour)
	c := Campaign{Status: "open", StartsAt: start, EndsAt: end}
	for _, tc := range []struct {
		name string
		at   time.Time
		ok   bool
	}{{"before", start.Add(-time.Second), false}, {"start", start, true}, {"inside", start.Add(7 * 24 * time.Hour), true}, {"end", end, false}, {"after", end.Add(time.Second), false}} {
		t.Run(tc.name, func(t *testing.T) {
			if got := c.OpenAt(tc.at) == nil; got != tc.ok {
				t.Fatalf("open=%v want=%v", got, tc.ok)
			}
		})
	}
}

func TestContextCancellationIsObservable(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("cancel did not propagate")
	}
}
