package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

// ArchiveService implements archive preconditions for the civic suggestion workflow.
type ArchiveService struct {
	Key      string
	Values   []string
	Limit    int
	Offset   int
	Deadline time.Time
	Metadata map[string]string
}

func (x ArchiveService) Validate() error {
	if strings.TrimSpace(x.Key) == "" {
		return fmt.Errorf("ArchiveService: key required")
	}
	if x.Limit < 0 || x.Offset < 0 {
		return fmt.Errorf("ArchiveService: page bounds invalid")
	}
	return nil
}
func (x ArchiveService) Normalize() ArchiveService {
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
func (x ArchiveService) Ready(now time.Time) bool {
	return !x.Deadline.IsZero() && !now.After(x.Deadline)
}
func (x ArchiveService) Expired(now time.Time) bool {
	return !x.Deadline.IsZero() && now.After(x.Deadline)
}
func (x ArchiveService) WithValue(v string) ArchiveService {
	y := x
	y.Values = append(append([]string(nil), x.Values...), v)
	return y
}
func (x ArchiveService) WithMetadata(k, v string) ArchiveService {
	y := x
	y.Metadata = map[string]string{}
	for a, b := range x.Metadata {
		y.Metadata[a] = b
	}
	y.Metadata[k] = v
	return y
}
func (x ArchiveService) Run(ctx context.Context) error {
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
func (x ArchiveService) Summary() string {
	return fmt.Sprintf("%s limit=%d offset=%d values=%d", x.Key, x.Limit, x.Offset, len(x.Values))
}
