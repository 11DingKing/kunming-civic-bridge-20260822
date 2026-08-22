package domain

import (
	"fmt"
	"strings"
)

type Visibility string

const (
	VisibilityInternal Visibility = "internal"
	VisibilityPublic   Visibility = "public"
	VisibilityRedacted Visibility = "redacted"
)

type ResponseVisibility struct {
	Status        ResponseStatus
	Visibility    Visibility
	Sensitive     bool
	RedactionNote string
}

func (v ResponseVisibility) Validate() error {
	if v.Status != ResponsePublished && v.Visibility == VisibilityPublic {
		return fmt.Errorf("%w: unpublished response public", ErrConflict)
	}
	if v.Visibility == VisibilityRedacted && strings.TrimSpace(v.RedactionNote) == "" {
		return fmt.Errorf("%w: redaction note", ErrInvalid)
	}
	return nil
}
func (v ResponseVisibility) Audience() string {
	if v.Sensitive {
		return "author"
	}
	if v.Visibility == VisibilityPublic {
		return "public"
	}
	return "staff"
}
