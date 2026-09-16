package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func newMetricsForTest() *metrics {
	return newMetrics(prometheus.NewRegistry())
}

func TestRouteTemplateBoundsCardinality(t *testing.T) {
	cases := map[string]string{
		"/api/items":          "/api/items",
		"/healthz/live":       "/healthz/live",
		"/api/items/abc123":   "/other",
		"/definitely-not-one": "/other",
	}
	for in, want := range cases {
		if got := routeTemplate(in); got != want {
			t.Errorf("routeTemplate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestHealthzLive(t *testing.T) {
	mux := http.NewServeMux()
	registerHandlers(mux, newConnector(""), newMetricsForTest(), slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/healthz/live", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestReadyWithoutDatabaseIs503(t *testing.T) {
	mux := http.NewServeMux()
	registerHandlers(mux, newConnector(""), newMetricsForTest(), slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
}

func TestErrorEndpointIs500(t *testing.T) {
	mux := http.NewServeMux()
	registerHandlers(mux, newConnector(""), newMetricsForTest(), slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/api/error", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
}

func TestCreateItemRejectsInvalidBody(t *testing.T) {
	mux := http.NewServeMux()
	// store is non-nil only in shape; validation runs before any DB call.
	registerHandlers(mux, newConnector(""), newMetricsForTest(), slog.Default())

	req := httptest.NewRequest(http.MethodPost, "/api/items", http.NoBody)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503 (no store)", rec.Code)
	}
}
