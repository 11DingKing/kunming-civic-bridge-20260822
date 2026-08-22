package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/service"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Svc service.Service
	Log *slog.Logger
}

func New(s service.Service, l *slog.Logger) *Server { return &Server{Svc: s, Log: l} }
func (s *Server) Routes() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { write(w, 200, map[string]any{"ok": true}) })
	m.HandleFunc("/readyz", s.ready)
	m.HandleFunc("/api/v1/suggestions", s.suggestions)
	m.HandleFunc("/api/v1/suggestions/", s.suggestion)
	m.HandleFunc("/api/v1/reviews", s.reviews)
	m.HandleFunc("/api/v1/assignments", s.assignments)
	m.HandleFunc("/api/v1/responses", s.responses)
	m.HandleFunc("/api/v1/feedback", s.feedback)
	return requestID(recoverer(m))
}
func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if e := s.Svc.DB.PingContext(r.Context()); e != nil {
		write(w, 503, map[string]any{"ok": false})
	} else {
		write(w, 200, map[string]any{"ok": true})
	}
}

type submitRequest struct{ AuthorID, CampaignID, IntakePointID, Title, Body, Scope string }

func (s *Server) suggestions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		write(w, 405, map[string]string{"error": "method_not_allowed"})
		return
	}
	var in submitRequest
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		write(w, 400, map[string]string{"error": "invalid_json"})
		return
	}
	x, e := s.Svc.Submit(r.Context(), in.AuthorID, in.CampaignID, in.IntakePointID, in.Title, in.Body, in.Scope, r.Header.Get("X-Request-ID"))
	if e != nil {
		problem(w, e)
		return
	}
	write(w, 201, x)
}
func (s *Server) suggestion(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/suggestions/")
	if id == "" {
		write(w, 404, nil)
		return
	}
	if r.Method == "GET" {
		x, e := s.Svc.Store.GetSuggestion(r.Context(), id)
		if e != nil {
			problem(w, e)
			return
		}
		write(w, 200, x)
		return
	}
	if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/submit") {
		x, e := s.Svc.Send(r.Context(), strings.TrimSuffix(id, "/submit"), r.Header.Get("X-Actor-ID"), r.Header.Get("X-Request-ID"))
		if e != nil {
			problem(w, e)
			return
		}
		write(w, 200, x)
		return
	}
	write(w, 405, nil)
}

type reviewRequest struct{ SuggestionID, ActorID, Decision, Note string }

func (s *Server) reviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		write(w, 405, nil)
		return
	}
	var in reviewRequest
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		problem(w, domain.ErrInvalid)
		return
	}
	x, e := s.Svc.Review(r.Context(), in.SuggestionID, in.ActorID, in.Decision, in.Note, r.Header.Get("X-Request-ID"))
	if e != nil {
		problem(w, e)
		return
	}
	write(w, 200, x)
}

type assignmentRequest struct{ SuggestionID, ActorID, DepartmentID string }

func (s *Server) assignments(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		write(w, 405, nil)
		return
	}
	var in assignmentRequest
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		problem(w, domain.ErrInvalid)
		return
	}
	if e := s.Svc.Assign(r.Context(), in.SuggestionID, in.ActorID, in.DepartmentID, r.Header.Get("X-Request-ID")); e != nil {
		problem(w, e)
		return
	}
	write(w, 202, map[string]string{"status": "assigned"})
}

type responseRequest struct{ SuggestionID, ActorID, Body string }

func (s *Server) responses(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		write(w, 405, nil)
		return
	}
	var in responseRequest
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		problem(w, domain.ErrInvalid)
		return
	}
	x, e := s.Svc.DraftResponse(r.Context(), in.SuggestionID, in.ActorID, in.Body, r.Header.Get("X-Request-ID"))
	if e != nil {
		problem(w, e)
		return
	}
	write(w, 201, x)
}

type feedbackRequest struct {
	SuggestionID, AuthorID, Comment string
	Rating                          int
}

func (s *Server) feedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		write(w, 405, nil)
		return
	}
	var in feedbackRequest
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		problem(w, domain.ErrInvalid)
		return
	}
	if e := s.Svc.AddFeedback(r.Context(), in.SuggestionID, in.AuthorID, in.Rating, in.Comment, r.Header.Get("X-Request-ID")); e != nil {
		problem(w, e)
		return
	}
	write(w, 201, map[string]string{"status": "recorded"})
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}
func problem(w http.ResponseWriter, e error) {
	status := 500
	code := "internal_error"
	switch {
	case errors.Is(e, domain.ErrNotFound):
		status = 404
		code = "not_found"
	case errors.Is(e, domain.ErrInvalid):
		status = 400
		code = "invalid_request"
	case errors.Is(e, domain.ErrConflict):
		status = 409
		code = "conflict"
	case errors.Is(e, domain.ErrForbidden):
		status = 403
		code = "forbidden"
	}
	write(w, status, map[string]string{"error": code, "message": e.Error()})
}
func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = time.Now().UTC().Format("20060102150405.000000000")
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				write(w, 500, map[string]string{"error": "internal_error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

var _ context.Context
