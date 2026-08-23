package domain

import (
	"fmt"
	"strings"
)

type SuggestionQuery struct {
	Scope, Status, Keyword string
	Limit, Offset          int
	IncludeClosed          bool
}

func (q SuggestionQuery) Normalize() SuggestionQuery {
	q.Scope = strings.TrimSpace(q.Scope)
	q.Status = strings.TrimSpace(q.Status)
	q.Keyword = strings.TrimSpace(q.Keyword)
	if q.Limit < 1 {
		q.Limit = 20
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	if q.Offset < 0 {
		q.Offset = 0
	}
	return q
}
func (q SuggestionQuery) Validate() error {
	if q.Limit < 1 || q.Limit > 100 {
		return fmt.Errorf("%w: query limit", ErrInvalid)
	}
	if q.Offset < 0 {
		return fmt.Errorf("%w: query offset", ErrInvalid)
	}
	if q.Status != "" {
		valid := map[SuggestionStatus]bool{StatusDraft: true, StatusSubmitted: true, StatusTriaged: true, StatusAssigned: true, StatusInProgress: true, StatusResponded: true, StatusClosed: true, StatusRejected: true, StatusEscalated: true, StatusReopened: true}
		if !valid[SuggestionStatus(q.Status)] {
			return fmt.Errorf("%w: query status", ErrInvalid)
		}
	}
	return nil
}
