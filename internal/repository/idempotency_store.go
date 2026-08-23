package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"time"
)

type IdempotencyStore struct{ DB *sql.DB }

func (s IdempotencyStore) Find(ctx context.Context, key string) (domain.IdempotencyRecord, error) {
	var r domain.IdempotencyRecord
	var created, expires string
	e := s.DB.QueryRowContext(ctx, `SELECT key,user_id,request_hash,response_json,expires_at FROM idempotency_keys WHERE key=?`, key).Scan(&r.Key, &r.UserID, &r.RequestHash, &r.ResponseJSON, &expires)
	if e != nil {
		return r, e
	}
	r.ExpiresAt, _ = time.Parse(time.RFC3339Nano, expires)
	_ = created
	return r, nil
}
func (s IdempotencyStore) Put(ctx context.Context, r domain.IdempotencyRecord) error {
	_, e := s.DB.ExecContext(ctx, `INSERT INTO idempotency_keys(key,user_id,request_hash,response_json,created_at,expires_at) VALUES(?,?,?,?,?,?)`, r.Key, r.UserID, r.RequestHash, r.ResponseJSON, time.Now().UTC().Format(time.RFC3339Nano), r.ExpiresAt.Format(time.RFC3339Nano))
	return e
}
func (s IdempotencyStore) Purge(ctx context.Context, now time.Time) (int64, error) {
	res, e := s.DB.ExecContext(ctx, `DELETE FROM idempotency_keys WHERE expires_at<=?`, now.Format(time.RFC3339Nano))
	if e != nil {
		return 0, e
	}
	return res.RowsAffected()
}
