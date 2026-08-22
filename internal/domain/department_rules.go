package domain

import (
	"fmt"
	"strings"
)

type DepartmentScope struct {
	DepartmentID, Name, Scope string
	Active                    bool
}

func (d DepartmentScope) Validate() error {
	if d.DepartmentID == "" || strings.TrimSpace(d.Name) == "" || d.Scope == "" {
		return fmt.Errorf("%w: department", ErrInvalid)
	}
	return nil
}
func (d DepartmentScope) Eligible(suggestionScope string) bool {
	return d.Active && (d.Scope == suggestionScope || strings.HasPrefix(suggestionScope, d.Scope+"/"))
}
func (d DepartmentScope) CanReceive() bool { return d.Active && d.DepartmentID != "" }
