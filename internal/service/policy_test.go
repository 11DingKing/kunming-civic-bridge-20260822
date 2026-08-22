package service

import (
	"context"
	"testing"
	"time"
)

func TestAuthenticationServicePolicy(t *testing.T) {
	x := AuthenticationService{Key: "authenticationservice", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
	if err := x.Validate(); err != nil {
		t.Fatal(err)
	}
	y := x.Normalize()
	if y.Limit != 50 || len(y.Values) != 2 || y.Values[0] != "a" {
		t.Fatalf("normalize=%+v", y)
	}
	if !y.Ready(time.Now()) || y.Expired(time.Now()) {
		t.Fatal("deadline state")
	}
	z := y.WithValue("v").WithMetadata("actor", "tester")
	if len(z.Values) != 3 || z.Metadata["actor"] != "tester" {
		t.Fatal("copy")
	}
	if err := z.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if z.Summary() == "" {
		t.Fatal("summary")
	}
}

func TestEscalationServicePolicy(t *testing.T) {
	x := EscalationService{Key: "escalationservice", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
	if err := x.Validate(); err != nil {
		t.Fatal(err)
	}
	y := x.Normalize()
	if y.Limit != 50 || len(y.Values) != 2 || y.Values[0] != "a" {
		t.Fatalf("normalize=%+v", y)
	}
	if !y.Ready(time.Now()) || y.Expired(time.Now()) {
		t.Fatal("deadline state")
	}
	z := y.WithValue("v").WithMetadata("actor", "tester")
	if len(z.Values) != 3 || z.Metadata["actor"] != "tester" {
		t.Fatal("copy")
	}
	if err := z.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if z.Summary() == "" {
		t.Fatal("summary")
	}
}

func TestNotificationServicePolicy(t *testing.T) {
	x := NotificationService{Key: "notificationservice", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
	if err := x.Validate(); err != nil {
		t.Fatal(err)
	}
	y := x.Normalize()
	if y.Limit != 50 || len(y.Values) != 2 || y.Values[0] != "a" {
		t.Fatalf("normalize=%+v", y)
	}
	if !y.Ready(time.Now()) || y.Expired(time.Now()) {
		t.Fatal("deadline state")
	}
	z := y.WithValue("v").WithMetadata("actor", "tester")
	if len(z.Values) != 3 || z.Metadata["actor"] != "tester" {
		t.Fatal("copy")
	}
	if err := z.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if z.Summary() == "" {
		t.Fatal("summary")
	}
}

func TestArchiveServicePolicy(t *testing.T) {
	x := ArchiveService{Key: "archiveservice", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
	if err := x.Validate(); err != nil {
		t.Fatal(err)
	}
	y := x.Normalize()
	if y.Limit != 50 || len(y.Values) != 2 || y.Values[0] != "a" {
		t.Fatalf("normalize=%+v", y)
	}
	if !y.Ready(time.Now()) || y.Expired(time.Now()) {
		t.Fatal("deadline state")
	}
	z := y.WithValue("v").WithMetadata("actor", "tester")
	if len(z.Values) != 3 || z.Metadata["actor"] != "tester" {
		t.Fatal("copy")
	}
	if err := z.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if z.Summary() == "" {
		t.Fatal("summary")
	}
}

func TestBatchServicePolicy(t *testing.T) {
	x := BatchService{Key: "batchservice", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
	if err := x.Validate(); err != nil {
		t.Fatal(err)
	}
	y := x.Normalize()
	if y.Limit != 50 || len(y.Values) != 2 || y.Values[0] != "a" {
		t.Fatalf("normalize=%+v", y)
	}
	if !y.Ready(time.Now()) || y.Expired(time.Now()) {
		t.Fatal("deadline state")
	}
	z := y.WithValue("v").WithMetadata("actor", "tester")
	if len(z.Values) != 3 || z.Metadata["actor"] != "tester" {
		t.Fatal("copy")
	}
	if err := z.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if z.Summary() == "" {
		t.Fatal("summary")
	}
}
