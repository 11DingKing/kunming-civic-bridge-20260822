package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

type AttachmentService struct {
	DB    *sql.DB
	Store repository.AttachmentStore
}

func (s AttachmentService) Add(ctx context.Context, suggestion, name, media, storage string, size int64) (string, error) {
	return s.Store.Create(ctx, suggestion, name, media, storage, size)
}
func (s AttachmentService) Archive(ctx context.Context, id string) error {
	return s.Store.Archive(ctx, id)
}
func (s AttachmentService) List(ctx context.Context, suggestion string) ([]repository.AttachmentRecord, error) {
	return s.Store.ForSuggestion(ctx, suggestion)
}
func (s AttachmentService) ValidateDownload(x repository.AttachmentRecord) error {
	if x.Archived {
		return domain.ErrForbidden
	}
	return nil
}
