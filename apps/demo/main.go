package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// --- Telemetry ---------------------------------------------------------
	shutdownTracing, err := initTracing(ctx, logger)
	if err != nil {
		logger.Error("tracing init failed", "event", "tracing_init_failed", "error", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownTracing(shutdownCtx)
	}()

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	metrics := newMetrics(registry)

	// --- Database (optional at boot; readiness reports it) -----------------
	store, err := newStore(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Warn("database not ready at startup", "event", "db_unavailable", "error", err)
	}
	if store != nil {
		defer store.Close()
	}

	// --- HTTP --------------------------------------------------------------
	mux := http.NewServeMux()
	registerHandlers(mux, store, metrics, logger)

	handler := instrument(mux, metrics, logger)

	srv := &http.Server{
		Addr:              ":" + envOr("PORT", "8080"),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("server listening", "event", "server_start", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "event", "server_error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down", "event", "server_shutdown")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

// registerHandlers wires every route. Kept in one place so the metric route
// labels stay in sync with the actual paths.
func registerHandlers(mux *http.ServeMux, store *store, metrics *metrics, logger *slog.Logger) {
	mux.HandleFunc("GET /healthz/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /healthz/ready", func(w http.ResponseWriter, r *http.Request) {
		if store == nil || store.Ping(r.Context()) != nil {
			http.Error(w, "database not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	mux.HandleFunc("GET /api/items", handleListItems(store, logger))
	mux.HandleFunc("POST /api/items", handleCreateItem(store, logger))
	mux.HandleFunc("GET /api/slow", handleSlow(logger))
	mux.HandleFunc("GET /api/error", handleError(logger))

	mux.Handle("/metrics", promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{}))
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
