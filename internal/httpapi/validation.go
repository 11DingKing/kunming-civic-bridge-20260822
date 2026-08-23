package httpapi

import (
	"encoding/json"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"net/http"
	"strings"
)

type ValidationError struct{ Field, Message string }

func decodeBody(r *http.Request, target any) error {
	if r.Body == nil {
		return domain.ErrInvalid
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if e := dec.Decode(target); e != nil {
		return e
	}
	return nil
}
func requireHeader(r *http.Request, name string) error {
	if strings.TrimSpace(r.Header.Get(name)) == "" {
		return domain.ErrInvalid
	}
	return nil
}
func writeError(w http.ResponseWriter, status int, code, message string) {
	write(w, status, map[string]string{"code": code, "message": message})
}
