package httpapi

import (
	"context"
	"testing"
	"time"
)

func TestResponseHandlersValidationMatrix(t *testing.T) {
	cases := []struct {
		name    string
		value   ResponseHandlers
		wantErr bool
	}{{"missing", ResponseHandlers{}, true}, {"negative", ResponseHandlers{Name: "x", Timeout: -time.Second}, true}, {"disabled", ResponseHandlers{Name: "x", Timeout: time.Second}, false}, {"ready", ResponseHandlers{Name: "x", Enabled: true, Timeout: time.Second}, false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.Validate(); (got != nil) != tc.wantErr {
				t.Fatalf("validate=%v wantErr=%v", got, tc.wantErr)
			}
		})
	}
	x := ResponseHandlers{Name: "initial", Enabled: true, Timeout: time.Second}
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

func TestFeedbackHandlersValidationMatrix(t *testing.T) {
	cases := []struct {
		name    string
		value   FeedbackHandlers
		wantErr bool
	}{{"missing", FeedbackHandlers{}, true}, {"negative", FeedbackHandlers{Name: "x", Timeout: -time.Second}, true}, {"disabled", FeedbackHandlers{Name: "x", Timeout: time.Second}, false}, {"ready", FeedbackHandlers{Name: "x", Enabled: true, Timeout: time.Second}, false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.Validate(); (got != nil) != tc.wantErr {
				t.Fatalf("validate=%v wantErr=%v", got, tc.wantErr)
			}
		})
	}
	x := FeedbackHandlers{Name: "initial", Enabled: true, Timeout: time.Second}
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
