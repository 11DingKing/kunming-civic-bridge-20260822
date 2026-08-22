package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

type AttachmentReportService struct {
	DB    *sql.DB
	Store repository.AttachmentStore
}

func (s AttachmentReportService) Count(ctx context.Context, id string) (int, error) {
	items, e := s.Store.ForSuggestion(ctx, id)
	return len(items), e
}
