package auth

import (
	"fmt"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
)

type Policy struct{}

func (Policy) Check(role domain.Role, action string) error {
	allowed := map[domain.Role]map[string]bool{domain.RoleCitizen: {"submit": true, "feedback": true}, domain.RoleReviewer: {"review": true, "assign": true, "read": true}, domain.RoleOperator: {"claim": true, "work": true, "respond": true, "read": true}, domain.RoleSupervisor: {"review": true, "assign": true, "claim": true, "work": true, "respond": true, "approve": true, "archive": true, "read": true}}
	if !allowed[role][action] {
		return fmt.Errorf("%w: action %s", domain.ErrForbidden, action)
	}
	return nil
}
