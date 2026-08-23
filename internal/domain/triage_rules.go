package domain

import (
	"fmt"
	"strings"
)

type TriageResult struct {
	Category, Priority, Reason string
	Sensitive                  bool
}

func (t TriageResult) Validate() error {
	if strings.TrimSpace(t.Category) == "" || strings.TrimSpace(t.Priority) == "" || strings.TrimSpace(t.Reason) == "" {
		return fmt.Errorf("%w: triage result", ErrInvalid)
	}
	if t.Priority != "low" && t.Priority != "normal" && t.Priority != "high" && t.Priority != "urgent" {
		return fmt.Errorf("%w: triage priority", ErrInvalid)
	}
	return nil
}
func (t TriageResult) RequiresSupervisor() bool { return t.Sensitive || t.Priority == "urgent" }
func (t TriageResult) PublicCategory() string {
	if t.Sensitive {
		return "民生事项"
	}
	return t.Category
}
