package domain

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

// RecoveryPlan implements restart recovery for the civic suggestion workflow.
type RecoveryPlan struct {
	Key      string
	Values   []string
	Limit    int
	Offset   int
	Deadline time.Time
	Metadata map[string]string
}

func (x RecoveryPlan) Validate() error {
	if strings.TrimSpace(x.Key) == "" {
		return fmt.Errorf("RecoveryPlan: key required")
	}
	if x.Limit < 0 || x.Offset < 0 {
		return fmt.Errorf("RecoveryPlan: page bounds invalid")
	}
	return nil
}
func (x RecoveryPlan) Normalize() RecoveryPlan {
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
func (x RecoveryPlan) Ready(now time.Time) bool {
	return !x.Deadline.IsZero() && !now.After(x.Deadline)
}
func (x RecoveryPlan) Expired(now time.Time) bool {
	return !x.Deadline.IsZero() && now.After(x.Deadline)
}
func (x RecoveryPlan) WithValue(v string) RecoveryPlan {
	y := x
	y.Values = append(append([]string(nil), x.Values...), v)
	return y
}
func (x RecoveryPlan) WithMetadata(k, v string) RecoveryPlan {
	y := x
	y.Metadata = map[string]string{}
	for a, b := range x.Metadata {
		y.Metadata[a] = b
	}
	y.Metadata[k] = v
	return y
}
func (x RecoveryPlan) Run(ctx context.Context) error {
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
func (x RecoveryPlan) Summary() string {
	return fmt.Sprintf("%s limit=%d offset=%d values=%d", x.Key, x.Limit, x.Offset, len(x.Values))
}
