package config

import (
	"testing"
	"time"
)

func TestRuntimeValidation(t *testing.T) {
	v := RuntimeValidation{Port: "8080", DBPath: "./data.db", SessionTTL: time.Hour, WorkerInterval: time.Second}
	if e := v.Validate(); e != nil {
		t.Fatal(e)
	}
	v.DBPath = "/data.db"
	if e := v.Validate(); e == nil {
		t.Fatal("root database path accepted")
	}
	if _, e := ExternalEndpoint("http://example.com"); e == nil {
		t.Fatal("http endpoint accepted")
	}
	if u, e := ExternalEndpoint("https://example.com/api"); e != nil || u.Host != "example.com" {
		t.Fatalf("endpoint=%v err=%v", u, e)
	}
}
