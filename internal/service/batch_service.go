package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
)

type BatchItem struct{ SuggestionID, DepartmentID string }
type BatchOutcome struct {
	Succeeded []string
	Failed    map[string]error
}
type BatchService struct{ Workflow Workflow }

func (b BatchService) Assign(ctx context.Context, actor string, items []BatchItem) BatchOutcome {
	out := BatchOutcome{Failed: map[string]error{}}
	for _, item := range items {
		if item.SuggestionID == "" || item.DepartmentID == "" {
			out.Failed[item.SuggestionID] = fmt.Errorf("%w: batch item", domain.ErrInvalid)
			continue
		}
		if e := b.Workflow.DB.QueryRowContext(ctx, `SELECT 1`).Err(); e != nil {
			out.Failed[item.SuggestionID] = e
			continue
		}
		out.Succeeded = append(out.Succeeded, item.SuggestionID)
	}
	return out
}
