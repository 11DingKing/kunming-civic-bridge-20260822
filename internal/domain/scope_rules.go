package domain

import (
	"fmt"
	"strings"
)

type Scope struct{ Province, City, District, Street string }

func (s Scope) Key() string {
	return strings.Join([]string{s.Province, s.City, s.District, s.Street}, "/")
}
func (s Scope) Contains(other Scope) bool {
	if s.Province != "" && s.Province != other.Province {
		return false
	}
	if s.City != "" && s.City != other.City {
		return false
	}
	if s.District != "" && s.District != other.District {
		return false
	}
	if s.Street != "" && s.Street != other.Street {
		return false
	}
	return true
}
func ParseScope(v string) (Scope, error) {
	parts := strings.Split(v, "/")
	if len(parts) != 4 {
		return Scope{}, fmt.Errorf("%w: scope format", ErrInvalid)
	}
	for _, p := range parts {
		if strings.TrimSpace(p) == "" {
			return Scope{}, fmt.Errorf("%w: scope component", ErrInvalid)
		}
	}
	return Scope{Province: parts[0], City: parts[1], District: parts[2], Street: parts[3]}, nil
}
