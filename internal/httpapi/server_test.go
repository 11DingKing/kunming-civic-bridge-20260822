package httpapi

import (
	"github.com/11DingKing/kunming-civic-bridge/internal/service"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	s := New(service.Service{}, nil)
	r := httptest.NewRecorder()
	q := httptest.NewRequest("GET", "/healthz", nil)
	s.Routes().ServeHTTP(r, q)
	if r.Code != 200 {
		t.Fatalf("status %d", r.Code)
	}
}
