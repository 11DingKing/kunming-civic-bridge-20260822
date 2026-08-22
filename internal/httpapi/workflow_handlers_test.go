package httpapi

import (
	"github.com/11DingKing/kunming-civic-bridge/internal/service"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWorkflowEndpointsRejectWrongMethods(t *testing.T) {
	s := New(service.Service{}, nil)
	for _, path := range []string{"/api/v1/reviews", "/api/v1/assignments", "/api/v1/responses", "/api/v1/feedback"} {
		r := httptest.NewRecorder()
		q := httptest.NewRequest("GET", path, strings.NewReader("{}"))
		s.Routes().ServeHTTP(r, q)
		if r.Code != 405 {
			t.Fatalf("%s status=%d", path, r.Code)
		}
	}
}
func TestSuggestionEndpointRejectsMalformedBody(t *testing.T) {
	s := New(service.Service{}, nil)
	r := httptest.NewRecorder()
	q := httptest.NewRequest("POST", "/api/v1/suggestions", strings.NewReader("{"))
	s.Routes().ServeHTTP(r, q)
	if r.Code != 400 {
		t.Fatalf("status=%d", r.Code)
	}
}
