package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
)

type DepartmentService struct {
	DB    *sql.DB
	Store repository.DepartmentStore
}

func (s DepartmentService) Eligible(ctx context.Context, scope string) ([]repository.DepartmentRecord, error) {
	return s.Store.ActiveForScope(ctx, scope)
}
func (s DepartmentService) Disable(ctx context.Context, id string) error {
	return s.Store.Disable(ctx, id)
}
