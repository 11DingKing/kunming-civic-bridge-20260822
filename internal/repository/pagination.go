package repository

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

// PaginationRepository implements paged queries for the civic suggestion workflow.
type PaginationRepository struct {
	Key      string
	Values   []string
	Limit    int
	Offset   int
	Deadline time.Time
	Metadata map[string]string
}

func (x PaginationRepository) Validate() error {
	if strings.TrimSpace(x.Key) == "" {
		return fmt.Errorf("PaginationRepository: key required")
	}
	if x.Limit < 0 || x.Offset < 0 {
		return fmt.Errorf("PaginationRepository: page bounds invalid")
	}
	return nil
}
func (x PaginationRepository) Normalize() PaginationRepository {
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
func (x PaginationRepository) Ready(now time.Time) bool {
	return !x.Deadline.IsZero() && !now.After(x.Deadline)
}
func (x PaginationRepository) Expired(now time.Time) bool {
	return !x.Deadline.IsZero() && now.After(x.Deadline)
}
func (x PaginationRepository) WithValue(v string) PaginationRepository {
	y := x
	y.Values = append(append([]string(nil), x.Values...), v)
	return y
}
func (x PaginationRepository) WithMetadata(k, v string) PaginationRepository {
	y := x
	y.Metadata = map[string]string{}
	for a, b := range x.Metadata {
		y.Metadata[a] = b
	}
	y.Metadata[k] = v
	return y
}
func (x PaginationRepository) Run(ctx context.Context) error {
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
func (x PaginationRepository) Summary() string {
	return fmt.Sprintf("%s limit=%d offset=%d values=%d", x.Key, x.Limit, x.Offset, len(x.Values))
}
