package domain

import (
	"fmt"
	"strings"
)

type AuditEntry struct {
	ActorID, Action, ObjectType, ObjectID, Result, RequestID string
	Metadata                                                 map[string]string
}

func (a AuditEntry) Validate() error {
	for k, v := range map[string]string{"actor": a.ActorID, "action": a.Action, "object_type": a.ObjectType, "object_id": a.ObjectID, "result": a.Result, "request_id": a.RequestID} {
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("%w: audit %s", ErrInvalid, k)
		}
	}
	return nil
}
func (a AuditEntry) Success() bool { return a.Result == "ok" }
