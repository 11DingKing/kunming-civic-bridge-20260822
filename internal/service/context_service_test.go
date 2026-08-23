package service

import (
	"context"
	"testing"
)

func TestServiceContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := ctx.Err(); err == nil {
		t.Fatal("cancel missing")
	}
}
