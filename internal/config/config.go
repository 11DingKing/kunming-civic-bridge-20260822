package config

import (
	"os"
	"time"
)

type Config struct {
	Port                       string
	DBPath                     string
	SessionTTL, WorkerInterval time.Duration
}

func Load() Config {
	return Config{Port: get("PORT", "8080"), DBPath: get("DB_PATH", "./civic-bridge.db"), SessionTTL: duration("SESSION_TTL", 24*time.Hour), WorkerInterval: duration("WORKER_INTERVAL", 2*time.Second)}
}
func get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func duration(k string, d time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if x, e := time.ParseDuration(v); e == nil {
			return x
		}
	}
	return d
}
