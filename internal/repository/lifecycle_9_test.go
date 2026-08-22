package repository

import (
	"context"
	"testing"
	"time"
)

func TestAssignmentRepositoryValidationMatrix(t *testing.T) {
	cases := []struct {
		name    string
		value   AssignmentRepository
		wantErr bool
	}{{"missing", AssignmentRepository{}, true}, {"negative", AssignmentRepository{Name: "x", Timeout: -time.Second}, true}, {"disabled", AssignmentRepository{Name: "x", Timeout: time.Second}, false}, {"ready", AssignmentRepository{Name: "x", Enabled: true, Timeout: time.Second}, false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.Validate(); (got != nil) != tc.wantErr {
				t.Fatalf("validate=%v wantErr=%v", got, tc.wantErr)
			}
		})
	}
	x := AssignmentRepository{Name: "initial", Enabled: true, Timeout: time.Second}
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

func TestResponseRepositoryValidationMatrix(t *testing.T) {
	cases := []struct {
		name    string
		value   ResponseRepository
		wantErr bool
	}{{"missing", ResponseRepository{}, true}, {"negative", ResponseRepository{Name: "x", Timeout: -time.Second}, true}, {"disabled", ResponseRepository{Name: "x", Timeout: time.Second}, false}, {"ready", ResponseRepository{Name: "x", Enabled: true, Timeout: time.Second}, false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.Validate(); (got != nil) != tc.wantErr {
				t.Fatalf("validate=%v wantErr=%v", got, tc.wantErr)
			}
		})
	}
	x := ResponseRepository{Name: "initial", Enabled: true, Timeout: time.Second}
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
