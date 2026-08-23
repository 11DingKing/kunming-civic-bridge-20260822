package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"time"
)

type IdempotencyService struct {
	DB    *sql.DB
	Store repository.IdempotencyStore
	TTL   time.Duration
}

func (s IdempotencyService) Replay(ctx context.Context, key, user string, body []byte, now time.Time) (string, bool, error) {
	record, e := s.Store.Find(ctx, key)
	if e == nil {
		if e = record.Usable(user, domain.HashRequest(body), now); e != nil {
			return "", false, e
		}
		return record.ResponseJSON, true, nil
	}
	if e != sql.ErrNoRows {
		return "", false, e
	}
	return "", false, nil
}
func (s IdempotencyService) Save(ctx context.Context, key, user string, body []byte, response string, now time.Time) error {
	return s.Store.Put(ctx, domain.IdempotencyRecord{Key: key, UserID: user, RequestHash: domain.HashRequest(body), ResponseJSON: response, ExpiresAt: now.Add(s.TTL)})
}
