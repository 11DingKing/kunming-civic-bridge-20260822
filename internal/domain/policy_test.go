package domain

import (
	"context"
	"testing"
	"time"
)

func TestFilterPolicy(t *testing.T) {
	x := Filter{Key: "filter", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
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

func TestPagePolicy(t *testing.T) {
	x := Page{Key: "page", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
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

func TestDeadlinePolicy(t *testing.T) {
	x := Deadline{Key: "deadline", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
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

func TestPermissionPolicy(t *testing.T) {
	x := Permission{Key: "permission", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
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

func TestTransitionPolicy(t *testing.T) {
	x := Transition{Key: "transition", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
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

func TestNotificationPolicy(t *testing.T) {
	x := Notification{Key: "notification", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
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

func TestAttachmentPolicy(t *testing.T) {
	x := Attachment{Key: "attachment", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
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

func TestBatchResultPolicy(t *testing.T) {
	x := BatchResult{Key: "batchresult", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
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

func TestTenantScopePolicy(t *testing.T) {
	x := TenantScope{Key: "tenantscope", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
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

func TestRecoveryPlanPolicy(t *testing.T) {
	x := RecoveryPlan{Key: "recoveryplan", Values: []string{"z", "a"}, Limit: 0, Deadline: time.Now().Add(time.Hour), Metadata: map[string]string{"scope": "x"}}
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
