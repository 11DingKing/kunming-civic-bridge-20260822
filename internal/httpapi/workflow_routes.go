package httpapi

import (
	"encoding/json"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/service"
	"net/http"
)

type WorkflowHTTP struct{ Service service.Service }

func (w WorkflowHTTP) Triage(suggestion string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			write(rw, 405, nil)
			return
		}
		var req struct {
			ActorID, Category, Priority, Reason string
			Sensitive                           bool
		}
		if decodeBody(r, &req) != nil {
			problem(rw, domain.ErrInvalid)
			return
		}
		if req.ActorID == "" || req.Category == "" || req.Priority == "" || req.Reason == "" || suggestion == "" {
			problem(rw, domain.ErrInvalid)
			return
		}
		write(rw, 202, map[string]string{"status": "queued"})
	})
}
func campaignJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
