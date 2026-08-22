package repository

import (
	"context"
	"testing"
	"time"
)

func TestAuditRepositoryValidationMatrix(t *testing.T) {
	cases := []struct {
		name    string
		value   AuditRepository
		wantErr bool
	}{{"missing", AuditRepository{}, true}, {"negative", AuditRepository{Name: "x", Timeout: -time.Second}, true}, {"disabled", AuditRepository{Name: "x", Timeout: time.Second}, false}, {"ready", AuditRepository{Name: "x", Enabled: true, Timeout: time.Second}, false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.Validate(); (got != nil) != tc.wantErr {
				t.Fatalf("validate=%v wantErr=%v", got, tc.wantErr)
			}
		})
	}
	x := AuditRepository{Name: "initial", Enabled: true, Timeout: time.Second}
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

func TestUserRepositoryValidationMatrix(t *testing.T) {
	cases := []struct {
		name    string
		value   UserRepository
		wantErr bool
	}{{"missing", UserRepository{}, true}, {"negative", UserRepository{Name: "x", Timeout: -time.Second}, true}, {"disabled", UserRepository{Name: "x", Timeout: time.Second}, false}, {"ready", UserRepository{Name: "x", Enabled: true, Timeout: time.Second}, false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.Validate(); (got != nil) != tc.wantErr {
				t.Fatalf("validate=%v wantErr=%v", got, tc.wantErr)
			}
		})
	}
	x := UserRepository{Name: "initial", Enabled: true, Timeout: time.Second}
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
