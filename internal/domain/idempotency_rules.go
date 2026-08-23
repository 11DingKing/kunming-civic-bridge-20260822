package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type IdempotencyRecord struct {
	Key, UserID, RequestHash, ResponseJSON string
	ExpiresAt                              time.Time
}

func HashRequest(body []byte) string { sum := sha256.Sum256(body); return hex.EncodeToString(sum[:]) }
func (r IdempotencyRecord) Usable(user, hash string, now time.Time) error {
	if r.UserID != user {
		return fmt.Errorf("%w: idempotency owner", ErrForbidden)
	}
	if r.RequestHash != hash {
		return fmt.Errorf("%w: request changed", ErrConflict)
	}
	if !now.Before(r.ExpiresAt) {
		return ErrExpired
	}
	return nil
}
