package main

import (
	"context"
	"errors"
	"github.com/11DingKing/kunming-civic-bridge/internal/audit"
	"github.com/11DingKing/kunming-civic-bridge/internal/clock"
	"github.com/11DingKing/kunming-civic-bridge/internal/config"
	"github.com/11DingKing/kunming-civic-bridge/internal/httpapi"
	"github.com/11DingKing/kunming-civic-bridge/internal/notify"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"github.com/11DingKing/kunming-civic-bridge/internal/repository"
	"github.com/11DingKing/kunming-civic-bridge/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	cfg := config.Load()
	db, e := platform.OpenDatabase(ctx, cfg.DBPath)
	if e != nil {
		panic(e)
	}
	defer db.Close()
	if e = platform.Migrate(ctx, db); e != nil {
		panic(e)
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	svc := service.Service{DB: db, Store: repository.Store{DB: db}, Audit: audit.Logger{DB: db}, Notify: notify.Queue{DB: db}, Clock: clock.Real{}}
	srv := &http.Server{Addr: ":" + cfg.Port, Handler: httpapi.New(svc, log).Routes(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		<-ctx.Done()
		shut, stop := context.WithTimeout(context.Background(), 10*time.Second)
		defer stop()
		_ = srv.Shutdown(shut)
	}()
	log.Info("server_started", "addr", srv.Addr)
	if e = srv.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
		log.Error("server_failed", "error", e)
		os.Exit(1)
	}
}
