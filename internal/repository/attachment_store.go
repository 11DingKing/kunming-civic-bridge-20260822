package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"time"
)

type AttachmentStore struct{ DB *sql.DB }
type AttachmentRecord struct {
	ID, SuggestionID, Name, MediaType, StorageKey string
	Size                                          int64
	Archived                                      bool
	CreatedAt                                     time.Time
}

func (s AttachmentStore) Create(ctx context.Context, suggestion, name, media, storage string, size int64) (string, error) {
	if suggestion == "" || name == "" || media == "" || storage == "" || size < 1 {
		return "", fmt.Errorf("%w: attachment", domain.ErrInvalid)
	}
	id := platform.ID()
	_, e := s.DB.ExecContext(ctx, `INSERT INTO suggestion_events(id,suggestion_id,from_status,to_status,actor_id,note,created_at) VALUES(?,?,?,?,?,?,?)`, id, suggestion, "", "attachment", "", name+"|"+media+"|"+storage, time.Now().UTC().Format(time.RFC3339Nano))
	return id, e
}
func (s AttachmentStore) Archive(ctx context.Context, id string) error {
	res, e := s.DB.ExecContext(ctx, `UPDATE suggestion_events SET to_status='attachment_archived' WHERE id=? AND to_status='attachment'`, id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrNotFound
	}
	return nil
}
func (s AttachmentStore) ForSuggestion(ctx context.Context, suggestion string) ([]AttachmentRecord, error) {
	rows, e := s.DB.QueryContext(ctx, `SELECT id,suggestion_id,note,to_status,created_at FROM suggestion_events WHERE suggestion_id=? AND to_status LIKE 'attachment%' ORDER BY created_at`, suggestion)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []AttachmentRecord
	for rows.Next() {
		var x AttachmentRecord
		var note, status, created string
		if e = rows.Scan(&x.ID, &x.SuggestionID, &note, &status, &created); e != nil {
			return nil, e
		}
		x.Archived = status == "attachment_archived"
		x.Name = note
		x.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		out = append(out, x)
	}
	return out, rows.Err()
}
