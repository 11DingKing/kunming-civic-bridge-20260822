package domain

import (
	"fmt"
	"strings"
	"time"
)

type FeedbackRules struct {
	Rating    int
	Comment   string
	CreatedAt time.Time
	CanReopen bool
}

func (f FeedbackRules) Validate() error {
	if e := ValidateRating(f.Rating); e != nil {
		return e
	}
	if strings.TrimSpace(f.Comment) == "" || len([]rune(f.Comment)) > 2000 {
		return fmt.Errorf("%w: feedback comment", ErrInvalid)
	}
	if f.CreatedAt.IsZero() {
		return fmt.Errorf("%w: feedback time", ErrInvalid)
	}
	return nil
}
func (f FeedbackRules) RequiresFollowUp() bool { return f.Rating <= 2 || f.CanReopen }
func (f FeedbackRules) Public() bool           { return f.Rating >= 3 }
