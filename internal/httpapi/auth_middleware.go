package httpapi

import (
	"context"
	"github.com/11DingKing/kunming-civic-bridge/internal/auth"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"net/http"
	"strings"
)

type principalKey struct{}
type Principal struct {
	UserID string
	Roles  []domain.Role
	Scope  string
}

func WithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, p)
}
func PrincipalFrom(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(Principal)
	return p, ok
}

type AuthMiddleware struct {
	Sessions auth.SessionService
	Users    func(context.Context, string) (Principal, error)
}

func (m AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
			next.ServeHTTP(w, r)
			return
		}
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			problem(w, domain.ErrForbidden)
			return
		}
		id := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		user, e := m.Sessions.Resolve(r.Context(), id)
		if e != nil {
			problem(w, e)
			return
		}
		p, e := m.Users(r.Context(), user)
		if e != nil {
			problem(w, e)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), p)))
	})
}
func RequireRole(role domain.Role, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, ok := PrincipalFrom(r.Context())
		if !ok {
			problem(w, domain.ErrForbidden)
			return
		}
		for _, candidate := range p.Roles {
			if candidate == role || candidate == domain.RoleSupervisor {
				next.ServeHTTP(w, r)
				return
			}
		}
		problem(w, domain.ErrForbidden)
	})
}
