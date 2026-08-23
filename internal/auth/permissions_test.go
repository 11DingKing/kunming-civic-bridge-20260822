package auth

import (
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"testing"
)

func TestRolePolicy(t *testing.T) {
	p := Policy{}
	if e := p.Check(domain.RoleCitizen, "submit"); e != nil {
		t.Fatal(e)
	}
	if e := p.Check(domain.RoleCitizen, "approve"); e == nil {
		t.Fatal("citizen approved")
	}
	if e := p.Check(domain.RoleSupervisor, "archive"); e != nil {
		t.Fatal(e)
	}
	for _, role := range []domain.Role{domain.RoleReviewer, domain.RoleOperator, domain.RoleSupervisor} {
		if e := p.Check(role, "read"); e != nil {
			t.Fatalf("%s read: %v", role, e)
		}
	}
}
