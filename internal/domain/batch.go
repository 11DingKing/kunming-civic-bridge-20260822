package domain

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

// BatchResult implements partial outcomes for the civic suggestion workflow.
type BatchResult struct {
	Key      string
	Values   []string
	Limit    int
	Offset   int
	Deadline time.Time
	Metadata map[string]string
}

func (x BatchResult) Validate() error {
	if strings.TrimSpace(x.Key) == "" {
		return fmt.Errorf("BatchResult: key required")
	}
	if x.Limit < 0 || x.Offset < 0 {
		return fmt.Errorf("BatchResult: page bounds invalid")
	}
	return nil
}
func (x BatchResult) Normalize() BatchResult {
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
func (x BatchResult) Ready(now time.Time) bool { return !x.Deadline.IsZero() && !now.After(x.Deadline) }
func (x BatchResult) Expired(now time.Time) bool {
	return !x.Deadline.IsZero() && now.After(x.Deadline)
}
func (x BatchResult) WithValue(v string) BatchResult {
	y := x
	y.Values = append(append([]string(nil), x.Values...), v)
	return y
}
func (x BatchResult) WithMetadata(k, v string) BatchResult {
	y := x
	y.Metadata = map[string]string{}
	for a, b := range x.Metadata {
		y.Metadata[a] = b
	}
	y.Metadata[k] = v
	return y
}
func (x BatchResult) Run(ctx context.Context) error {
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
func (x BatchResult) Summary() string {
	return fmt.Sprintf("%s limit=%d offset=%d values=%d", x.Key, x.Limit, x.Offset, len(x.Values))
}
