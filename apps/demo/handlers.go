package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func handleListItems(conn *connector, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := conn.get(r.Context())
		if s == nil {
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		items, err := s.ListItems(r.Context())
		if err != nil {
			logger.Error("list items failed", "event", "items_list_failed", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}

func handleCreateItem(conn *connector, logger *slog.Logger) http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s := conn.get(r.Context())
		if s == nil {
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		var req request
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&req); err != nil || req.Name == "" {
			http.Error(w, "invalid body: name is required", http.StatusBadRequest)
			return
		}
		it, err := s.CreateItem(r.Context(), req.Name)
		if err != nil {
			logger.Error("create item failed", "event", "item_create_failed", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusCreated, it)
	}
}

// handleSlow sleeps for ?ms= (capped) to exercise latency SLOs and dashboards.
func handleSlow(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ms, _ := strconv.Atoi(r.URL.Query().Get("ms"))
		if ms < 0 {
			ms = 0
		}
		if ms > 5000 {
			ms = 5000
		}
		select {
		case <-time.After(time.Duration(ms) * time.Millisecond):
		case <-r.Context().Done():
			return
		}
		logger.Info("slow request served", "event", "slow_request", "ms", ms)
		writeJSON(w, http.StatusOK, map[string]int{"sleptMs": ms})
	}
}

// handleError always fails, for error-rate SLO and alerting demos.
func handleError(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Error("intentional error", "event", "intentional_error")
		http.Error(w, "intentional error", http.StatusInternalServerError)
	}
}
