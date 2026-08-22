package domain

import (
	"context"
	"testing"
	"time"
)

func TestCampaignPolicyValidationMatrix(t *testing.T) {
	cases := []struct {
		name    string
		value   CampaignPolicy
		wantErr bool
	}{{"missing", CampaignPolicy{}, true}, {"negative", CampaignPolicy{Name: "x", Timeout: -time.Second}, true}, {"disabled", CampaignPolicy{Name: "x", Timeout: time.Second}, false}, {"ready", CampaignPolicy{Name: "x", Enabled: true, Timeout: time.Second}, false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.Validate(); (got != nil) != tc.wantErr {
				t.Fatalf("validate=%v wantErr=%v", got, tc.wantErr)
			}
		})
	}
	x := CampaignPolicy{Name: "initial", Enabled: true, Timeout: time.Second}
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

func TestIntakePolicyValidationMatrix(t *testing.T) {
	cases := []struct {
		name    string
		value   IntakePolicy
		wantErr bool
	}{{"missing", IntakePolicy{}, true}, {"negative", IntakePolicy{Name: "x", Timeout: -time.Second}, true}, {"disabled", IntakePolicy{Name: "x", Timeout: time.Second}, false}, {"ready", IntakePolicy{Name: "x", Enabled: true, Timeout: time.Second}, false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.Validate(); (got != nil) != tc.wantErr {
				t.Fatalf("validate=%v wantErr=%v", got, tc.wantErr)
			}
		})
	}
	x := IntakePolicy{Name: "initial", Enabled: true, Timeout: time.Second}
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
