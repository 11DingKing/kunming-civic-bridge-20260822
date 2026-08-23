package httpapi

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeBodyRejectsUnknownFields(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{"title":"x","extra":1}`))
	var v struct {
		Title string `json:"title"`
	}
	if e := decodeBody(r, &v); e == nil {
		t.Fatal("unknown field accepted")
	}
}
func TestRequireHeader(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	if e := requireHeader(r, "X-Actor-ID"); e == nil {
		t.Fatal("missing header accepted")
	}
	r.Header.Set("X-Actor-ID", "u")
	if e := requireHeader(r, "X-Actor-ID"); e != nil {
		t.Fatal(e)
	}
}
