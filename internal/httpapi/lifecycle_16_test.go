package httpapi

import (
	"context"
	"testing"
	"time"
)

func TestReviewHandlersValidationMatrix(t *testing.T) {
	cases := []struct {
		name    string
		value   ReviewHandlers
		wantErr bool
	}{{"missing", ReviewHandlers{}, true}, {"negative", ReviewHandlers{Name: "x", Timeout: -time.Second}, true}, {"disabled", ReviewHandlers{Name: "x", Timeout: time.Second}, false}, {"ready", ReviewHandlers{Name: "x", Enabled: true, Timeout: time.Second}, false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.Validate(); (got != nil) != tc.wantErr {
				t.Fatalf("validate=%v wantErr=%v", got, tc.wantErr)
			}
		})
	}
	x := ReviewHandlers{Name: "initial", Enabled: true, Timeout: time.Second}
	if !x.Ready() {
		t.Fatal("ready expected")
	}
	if x.Status() != "ready" {
		t.Fatal("status")
	}
	if x.Describe() == "" {
		t.Fatal("description")
	}
	y := x.WithName("changed").WithEnabled(false).WithTimeout(2 * time.Second)
	if y.Name != "changed" || y.Enabled || y.Timeout != 2*time.Second {
		t.Fatalf("copy=%+v", y)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := x.Start(ctx); err == nil {
		t.Fatal("cancel should propagate")
	}
	if err := x.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestAssignmentHandlersValidationMatrix(t *testing.T) {
	cases := []struct {
		name    string
		value   AssignmentHandlers
		wantErr bool
	}{{"missing", AssignmentHandlers{}, true}, {"negative", AssignmentHandlers{Name: "x", Timeout: -time.Second}, true}, {"disabled", AssignmentHandlers{Name: "x", Timeout: time.Second}, false}, {"ready", AssignmentHandlers{Name: "x", Enabled: true, Timeout: time.Second}, false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.Validate(); (got != nil) != tc.wantErr {
				t.Fatalf("validate=%v wantErr=%v", got, tc.wantErr)
			}
		})
	}
	x := AssignmentHandlers{Name: "initial", Enabled: true, Timeout: time.Second}
	if !x.Ready() {
		t.Fatal("ready expected")
	}
	if x.Status() != "ready" {
		t.Fatal("status")
	}
	if x.Describe() == "" {
		t.Fatal("description")
	}
	y := x.WithName("changed").WithEnabled(false).WithTimeout(2 * time.Second)
	if y.Name != "changed" || y.Enabled || y.Timeout != 2*time.Second {
		t.Fatalf("copy=%+v", y)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := x.Start(ctx); err == nil {
		t.Fatal("cancel should propagate")
	}
	if err := x.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
}
