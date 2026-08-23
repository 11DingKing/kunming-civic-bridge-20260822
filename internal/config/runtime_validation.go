package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type RuntimeValidation struct {
	Port                       string
	DBPath                     string
	SessionTTL, WorkerInterval time.Duration
}

func (v RuntimeValidation) Validate() error {
	if v.Port == "" {
		return fmt.Errorf("port required")
	}
	if v.DBPath == "" || filepath.IsAbs(v.DBPath) && filepath.Dir(v.DBPath) == string(filepath.Separator) {
		return fmt.Errorf("database path unsafe")
	}
	if v.SessionTTL < time.Minute {
		return fmt.Errorf("session ttl too short")
	}
	if v.WorkerInterval <= 0 {
		return fmt.Errorf("worker interval invalid")
	}
	return nil
}
func ExternalEndpoint(raw string) (*url.URL, error) {
	u, e := url.Parse(raw)
	if e != nil || u.Scheme != "https" || u.Host == "" {
		return nil, fmt.Errorf("endpoint must be https")
	}
	return u, nil
}
func EnvSnapshot() map[string]string {
	out := map[string]string{}
	for _, k := range []string{"PORT", "DB_PATH", "SESSION_TTL", "WORKER_INTERVAL"} {
		out[k] = os.Getenv(k)
	}
	return out
}
