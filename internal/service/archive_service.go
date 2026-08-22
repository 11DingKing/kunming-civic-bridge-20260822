package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
)

type ArchiveService struct{ DB *sql.DB }

func (a ArchiveService) Eligible(ctx context.Context, id string) error {
	var status string
	if e := a.DB.QueryRowContext(ctx, `SELECT status FROM suggestions WHERE id=?`, id).Scan(&status); e != nil {
		return e
	}
	if status != "closed" {
		return fmt.Errorf("%w: status %s", domain.ErrConflict, status)
	}
	var pending int
	if e := a.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE payload LIKE '%'||?||'%' AND status IN ('pending','running')`, id).Scan(&pending); e != nil {
		return e
	}
	if pending > 0 {
		return fmt.Errorf("%w: pending jobs", domain.ErrConflict)
	}
	return nil
}
