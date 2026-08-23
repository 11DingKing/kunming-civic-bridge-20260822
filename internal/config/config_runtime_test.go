package config

import (
	"os"
	"testing"
)

func TestEnvSnapshot(t *testing.T) {
	os.Setenv("PORT", "9000")
	defer os.Unsetenv("PORT")
	if got := EnvSnapshot()["PORT"]; got != "9000" {
		t.Fatalf("port=%s", got)
	}
}
