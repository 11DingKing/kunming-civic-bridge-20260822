package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/auth"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"time"
)

type AuthService struct {
	DB         *sql.DB
	Identities repository.IdentityStore
	Sessions   auth.SessionService
	Policy     auth.Policy
	Now        func() time.Time
}

func (a AuthService) Login(ctx context.Context, phone, password string) (string, domain.User, error) {
	u, e := a.Identities.Authenticate(ctx, phone, password)
	if e != nil {
		return "", u, e
	}
	id, e := a.Sessions.Issue(ctx, u.ID)
	return id, u, e
}
func (a AuthService) Logout(ctx context.Context, session string) error {
	return a.Sessions.Revoke(ctx, session)
}
func (a AuthService) Authorize(role domain.Role, action string) error {
	return a.Policy.Check(role, action)
}
func (a AuthService) ExpireAt() time.Time { return a.Now().Add(a.Sessions.TTL) }
