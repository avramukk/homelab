package main

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const serviceName = "demo"

// metrics holds the RED signals every endpoint and dependency reports.
// Labels stay bounded: route template and status class, never raw URLs or IDs.
type metrics struct {
	registry     *prometheus.Registry
	requests     *prometheus.CounterVec
	duration     *prometheus.HistogramVec
	dependencyUp *prometheus.GaugeVec
}

func newMetrics(reg *prometheus.Registry) *metrics {
	m := &metrics{
		registry: reg,
		requests: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests by route, method and status class.",
		}, []string{"route", "method", "status_class"}),
		duration: promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		}, []string{"route", "method"}),
		dependencyUp: promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "dependency_up",
			Help: "Dependency liveness as seen by the service (1 = up).",
		}, []string{"dependency"}),
	}
	return m
}

// statusRecorder captures the status code for metrics and logs.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// instrument adds request ID + structured logging + RED metrics + tracing.
func instrument(next http.Handler, m *metrics, logger *slog.Logger) http.Handler {
	traced := otelhttp.NewHandler(next, "http.server",
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			return r.Method + " " + routeTemplate(r.URL.Path)
		}),
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := r.Header.Get("x-request-id")
		if requestID == "" {
			requestID = newID()
		}
		w.Header().Set("x-request-id", requestID)

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, requestID)
		ctx, span := otel.Tracer(serviceName).Start(ctx, routeTemplate(r.URL.Path))
		defer span.End()

		traced.ServeHTTP(rec, r.WithContext(ctx))

		route := routeTemplate(r.URL.Path)
		statusClass := strconv.Itoa(rec.status/100) + "xx"
		m.requests.WithLabelValues(route, r.Method, statusClass).Inc()
		m.duration.WithLabelValues(route, r.Method).Observe(time.Since(start).Seconds())

		logger.Info("http_request",
			"event", "http_request",
			"requestId", requestID,
			"traceId", span.SpanContext().TraceID().String(),
			"route", route,
			"method", r.Method,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// routeTemplate maps a concrete path to a bounded route label so metric
// cardinality stays fixed.
func routeTemplate(path string) string {
	switch path {
	case "/healthz/live", "/healthz/ready", "/api/items", "/api/slow", "/api/error", "/metrics":
		return path
	default:
		return "/other"
	}
}

type ctxKeyRequestID struct{}

// initTracing configures an OTLP/gRPC exporter (Alloy receives it on :4317)
// and returns a shutdown function.
func initTracing(ctx context.Context, logger *slog.Logger) (func(context.Context) error, error) {
	endpoint := envOr("OTEL_EXPORTER_OTLP_ENDPOINT", "alloy.observability.svc.cluster.local:4317")
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return func(context.Context) error { return nil }, err
	}

	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(envOr("APP_VERSION", "dev")),
		),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		// Keep every error trace; sample the rest lightly.
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.25))),
	)
	otel.SetTracerProvider(tp)
	logger.Info("tracing enabled", "event", "tracing_init", "endpoint", endpoint)

	trace.SpanFromContext(ctx)
	return tp.Shutdown, nil
}
