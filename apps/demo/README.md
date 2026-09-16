# demo service

The SLO-bearing workload of the homelab: a small Go + PostgreSQL service
instrumented so that the observability stack has something real to observe.

## Endpoints

| Method | Path | Purpose |
|---|---|---|
| GET | `/healthz/live` | liveness (never touches the DB) |
| GET | `/healthz/ready` | readiness (pings PostgreSQL; 503 when down) |
| GET | `/api/items` | list items |
| POST | `/api/items` | create an item (`{"name": "..."}`) |
| GET | `/api/slow?ms=N` | sleep N ms (≤5000) — latency SLO / dashboard demos |
| GET | `/api/error` | always 500 — error-rate SLO and alerting demos |
| GET | `/metrics` | Prometheus metrics |

## Telemetry

Per the instrumentation rules: metrics tell you *that* something is wrong, traces
tell you *where*, logs tell you *why*.

- **Logs:** structured JSON (`slog`) with `event`, `requestId`, `traceId`, `route`,
  `status`, `duration_ms`. No secrets or PII.
- **Metrics:** RED per route — `http_requests_total{route,method,status_class}` and
  `http_request_duration_seconds{route,method}` (histogram, bounded labels only).
- **Traces:** OpenTelemetry OTLP/gRPC to `OTEL_EXPORTER_OTLP_ENDPOINT`
  (default `alloy.observability.svc.cluster.local:4317` → Tempo). Errors are kept;
  the rest are sampled at 25%.

## Configuration

| Env var | Default | Purpose |
|---|---|---|
| `PORT` | `8080` | listen port |
| `DATABASE_URL` | — | PostgreSQL DSN (required) |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `alloy.observability:4317` | trace export |
| `APP_VERSION` | `dev` | reported as a resource attribute |

## Local

```bash
export PATH="/opt/homebrew/bin:$PATH"
go test ./...
go build ./...
docker build -t homelab-demo:dev .
```
