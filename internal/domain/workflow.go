package domain

import (
	"fmt"
	"strings"
	"time"
)

type IntakeMode string

const (
	ModeOnline  IntakeMode = "online"
	ModePoint   IntakeMode = "point"
	ModeInvited IntakeMode = "invited"
)

type ReviewDecision string

const (
	DecisionAccept ReviewDecision = "accept"
	DecisionReject ReviewDecision = "reject"
	DecisionReturn ReviewDecision = "return"
)

type AssignmentStatus string

const (
	AssignmentAssigned AssignmentStatus = "assigned"
	AssignmentClaimed  AssignmentStatus = "claimed"
	AssignmentDone     AssignmentStatus = "done"
	AssignmentDeclined AssignmentStatus = "declined"
)

type ResponseStatus string

const (
	ResponseDraft     ResponseStatus = "draft"
	ResponsePending   ResponseStatus = "pending_review"
	ResponseApproved  ResponseStatus = "approved"
	ResponsePublished ResponseStatus = "published"
	ResponseRecalled  ResponseStatus = "recalled"
)

type SuggestionDraft struct {
	CampaignID, AuthorID, IntakePointID, Title, Body, Scope string
	Mode                                                    IntakeMode
	DueAt                                                   *time.Time
}

func (d SuggestionDraft) Validate() error {
	if d.CampaignID == "" || d.AuthorID == "" {
		return fmt.Errorf("%w: identity", ErrInvalid)
	}
	if d.Mode != ModeOnline && d.Mode != ModePoint && d.Mode != ModeInvited {
		return fmt.Errorf("%w: mode", ErrInvalid)
	}
	return ValidateSuggestion(d.Title, d.Body, d.Scope)
}

func CanReview(actor Role) bool  { return actor == RoleReviewer || actor == RoleSupervisor }
func CanAssign(actor Role) bool  { return actor == RoleReviewer || actor == RoleSupervisor }
func CanOperate(actor Role) bool { return actor == RoleOperator || actor == RoleSupervisor }
func CanApprove(actor Role) bool { return actor == RoleSupervisor }

func ValidateReview(decision ReviewDecision, note string) error {
	if decision != DecisionAccept && decision != DecisionReject && decision != DecisionReturn {
		return fmt.Errorf("%w: decision", ErrInvalid)
	}
	if strings.TrimSpace(note) == "" {
		return fmt.Errorf("%w: note", ErrInvalid)
	}
	return nil
}
func ValidateResponse(body string) error {
	body = strings.TrimSpace(body)
	if len([]rune(body)) < 20 {
		return fmt.Errorf("%w: response too short", ErrInvalid)
	}
	if len([]rune(body)) > 20000 {
		return fmt.Errorf("%w: response too long", ErrInvalid)
	}
	return nil
}
func ResponseTransition(from, to ResponseStatus) bool {
	return map[ResponseStatus]map[ResponseStatus]bool{ResponseDraft: {ResponsePending: true}, ResponsePending: {ResponseApproved: true, ResponseDraft: true}, ResponseApproved: {ResponsePublished: true, ResponseDraft: true}, ResponsePublished: {ResponseRecalled: true}}[from][to]
}
