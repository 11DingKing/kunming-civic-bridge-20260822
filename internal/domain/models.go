package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrConflict  = errors.New("conflict")
	ErrForbidden = errors.New("forbidden")
	ErrInvalid   = errors.New("invalid request")
	ErrExpired   = errors.New("expired")
	ErrClosed    = errors.New("closed")
)

type Role string

const (
	RoleCitizen    Role = "citizen"
	RoleReviewer   Role = "reviewer"
	RoleOperator   Role = "operator"
	RoleSupervisor Role = "supervisor"
)

func (r Role) Valid() bool {
	return r == RoleCitizen || r == RoleReviewer || r == RoleOperator || r == RoleSupervisor
}

type SuggestionStatus string

const (
	StatusDraft      SuggestionStatus = "draft"
	StatusSubmitted  SuggestionStatus = "submitted"
	StatusTriaged    SuggestionStatus = "triaged"
	StatusAssigned   SuggestionStatus = "assigned"
	StatusInProgress SuggestionStatus = "in_progress"
	StatusResponded  SuggestionStatus = "responded"
	StatusClosed     SuggestionStatus = "closed"
	StatusRejected   SuggestionStatus = "rejected"
	StatusEscalated  SuggestionStatus = "escalated"
	StatusReopened   SuggestionStatus = "reopened"
)

func (s SuggestionStatus) CanTransition(to SuggestionStatus) bool {
	transitions := map[SuggestionStatus]map[SuggestionStatus]bool{
		StatusDraft: {StatusSubmitted: true}, StatusSubmitted: {StatusTriaged: true, StatusRejected: true}, StatusTriaged: {StatusAssigned: true, StatusRejected: true}, StatusAssigned: {StatusInProgress: true, StatusEscalated: true}, StatusInProgress: {StatusResponded: true, StatusEscalated: true}, StatusResponded: {StatusClosed: true, StatusReopened: true}, StatusReopened: {StatusAssigned: true, StatusClosed: true}, StatusEscalated: {StatusAssigned: true, StatusInProgress: true}, StatusClosed: {}, StatusRejected: {},
	}
	return transitions[s][to]
}
func (s SuggestionStatus) String() string { return string(s) }

type User struct {
	ID, Name, Phone, PasswordHash, Scope string
	Roles                                []Role
	CreatedAt                            time.Time
}
type Campaign struct {
	ID, Name, Status            string
	StartsAt, EndsAt, CreatedAt time.Time
}

func (c Campaign) OpenAt(now time.Time) error {
	if c.Status != "open" || now.Before(c.StartsAt) || !now.Before(c.EndsAt) {
		return ErrClosed
	}
	return nil
}

type IntakePoint struct {
	ID, Name, Scope string
	Active          bool
	CreatedAt       time.Time
}
type Suggestion struct {
	ID, CampaignID, IntakePointID, AuthorID, Title, Body, Scope string
	Status                                                      SuggestionStatus
	Version                                                     int
	DueAt                                                       *time.Time
	CreatedAt, UpdatedAt                                        time.Time
}
type Review struct {
	ID, SuggestionID, ReviewerID, Decision, Note string
	Version                                      int
	CreatedAt                                    time.Time
}
type Department struct {
	ID, Name, Scope string
	Active          bool
	CreatedAt       time.Time
}
type Assignment struct {
	ID, SuggestionID, DepartmentID, AssigneeID, Status, LeaseToken string
	LeaseUntil                                                     *time.Time
	Version                                                        int
	AssignedAt                                                     time.Time
	CompletedAt                                                    *time.Time
}
type Response struct {
	ID, SuggestionID, AuthorID, Body, Status string
	ReviewedBy                               string
	PublishedAt                              *time.Time
	CreatedAt                                time.Time
}
type Feedback struct {
	ID, SuggestionID, AuthorID, Comment string
	Rating                              int
	CreatedAt                           time.Time
}
type Job struct {
	ID, Kind, Payload, Status, LastError string
	Attempts                             int
	RunAfter                             time.Time
	LeaseUntil                           *time.Time
	CreatedAt, UpdatedAt                 time.Time
}

func ValidateSuggestion(title, body, scope string) error {
	if strings.TrimSpace(title) == "" || len([]rune(title)) > 120 {
		return fmt.Errorf("%w: title", ErrInvalid)
	}
	if strings.TrimSpace(body) == "" || len([]rune(body)) > 10000 {
		return fmt.Errorf("%w: body", ErrInvalid)
	}
	if strings.TrimSpace(scope) == "" {
		return fmt.Errorf("%w: scope", ErrInvalid)
	}
	return nil
}
func ValidateRating(r int) error {
	if r < 1 || r > 5 {
		return fmt.Errorf("%w: rating", ErrInvalid)
	}
	return nil
}
