package domain

import (
	"fmt"
	"time"
)

type AuditExport struct {
	ObjectType, ObjectID, RequestedBy string
	From, To                          time.Time
	Format                            string
}

func (e AuditExport) Validate() error {
	if e.ObjectType == "" || e.ObjectID == "" || e.RequestedBy == "" {
		return fmt.Errorf("%w: export identity", ErrInvalid)
	}
	if e.Format != "json" && e.Format != "csv" {
		return fmt.Errorf("%w: export format", ErrInvalid)
	}
	if !e.From.Before(e.To) {
		return fmt.Errorf("%w: export window", ErrInvalid)
	}
	return nil
}
func (e AuditExport) Includes(at time.Time) bool { return !at.Before(e.From) && at.Before(e.To) }
