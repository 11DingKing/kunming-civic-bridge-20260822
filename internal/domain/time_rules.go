package domain

import (
	"fmt"
	"time"
)

type BusinessCalendar struct {
	Location                 *time.Location
	WorkingStart, WorkingEnd int
}

func (c BusinessCalendar) Validate() error {
	if c.Location == nil {
		return fmt.Errorf("%w: calendar location", ErrInvalid)
	}
	if c.WorkingStart < 0 || c.WorkingEnd > 24 || c.WorkingStart >= c.WorkingEnd {
		return fmt.Errorf("%w: calendar hours", ErrInvalid)
	}
	return nil
}
func (c BusinessCalendar) InHours(t time.Time) bool {
	if c.Location == nil {
		return false
	}
	local := t.In(c.Location)
	h := local.Hour()
	return h >= c.WorkingStart && h < c.WorkingEnd && local.Weekday() != time.Saturday && local.Weekday() != time.Sunday
}
func (c BusinessCalendar) NextOpening(t time.Time) time.Time {
	for i := 0; i < 8; i++ {
		local := t.In(c.Location).AddDate(0, 0, i)
		if local.Weekday() != time.Saturday && local.Weekday() != time.Sunday {
			return time.Date(local.Year(), local.Month(), local.Day(), c.WorkingStart, 0, 0, 0, c.Location).UTC()
		}
	}
	return t
}
