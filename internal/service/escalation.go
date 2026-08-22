package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

// EscalationService implements deadline escalation for the civic suggestion workflow.
type EscalationService struct {
	Key      string
	Values   []string
	Limit    int
	Offset   int
	Deadline time.Time
	Metadata map[string]string
}

func (x EscalationService) Validate() error {
	if strings.TrimSpace(x.Key) == "" {
		return fmt.Errorf("EscalationService: key required")
	}
	if x.Limit < 0 || x.Offset < 0 {
		return fmt.Errorf("EscalationService: page bounds invalid")
	}
	return nil
}
func (x EscalationService) Normalize() EscalationService {
	y := x
	y.Key = strings.TrimSpace(y.Key)
	y.Values = append([]string(nil), y.Values...)
	sort.Strings(y.Values)
	if y.Limit == 0 {
		y.Limit = 50
	}
	if y.Limit > 200 {
		y.Limit = 200
	}
	return y
}
func (x EscalationService) Ready(now time.Time) bool {
	return !x.Deadline.IsZero() && !now.After(x.Deadline)
}
func (x EscalationService) Expired(now time.Time) bool {
	return !x.Deadline.IsZero() && now.After(x.Deadline)
}
func (x EscalationService) WithValue(v string) EscalationService {
	y := x
	y.Values = append(append([]string(nil), x.Values...), v)
	return y
}
func (x EscalationService) WithMetadata(k, v string) EscalationService {
	y := x
	y.Metadata = map[string]string{}
	for a, b := range x.Metadata {
		y.Metadata[a] = b
	}
	y.Metadata[k] = v
	return y
}
func (x EscalationService) Run(ctx context.Context) error {
	if err := x.Validate(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
func (x EscalationService) Summary() string {
	return fmt.Sprintf("%s limit=%d offset=%d values=%d", x.Key, x.Limit, x.Offset, len(x.Values))
}
