package domain

import (
	"fmt"
	"time"
)

type AssignmentEvent struct {
	AssignmentID, ActorID, From, To, Reason string
	At                                      time.Time
}

func (e AssignmentEvent) Validate() error {
	if e.AssignmentID == "" || e.ActorID == "" || e.To == "" {
		return fmt.Errorf("%w: assignment event", ErrInvalid)
	}
	if e.At.IsZero() {
		return fmt.Errorf("%w: assignment event time", ErrInvalid)
	}
	return nil
}
func (e AssignmentEvent) Terminal() bool {
	return e.To == string(AssignmentDone) || e.To == string(AssignmentDeclined)
}
func (e AssignmentEvent) LeaseRequired() bool { return e.To == string(AssignmentClaimed) }
func (e AssignmentEvent) Summary() string     { return e.From + "->" + e.To + " by " + e.ActorID }
