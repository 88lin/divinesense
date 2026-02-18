# AI Observability (`ai/observability`)

The `observability` package provides full-chain observability support for AI modules, including logging, metrics monitoring, and distributed tracing.

## Module Components

### 1. `logging`
Structured logging wrapper based on `slog`.
- Unified log format (JSON/Text)
- Support for context injection (TraceID, SpanID)

### 2. `metrics`
Metrics collection system based on Prometheus.

**Core Metrics**:
- `ai_request_total`: Total request counter
- `ai_request_duration_seconds`: Request duration histogram
- `ai_token_usage`: Token usage counter (Input/Output)
- `ai_error_total`: Error counter

### 3. `tracing`
Distributed tracing based on OpenTelemetry.
- **`Tracer`**: Automatically creates Spans for AI requests, records key steps (LLM calls, RAG retrieval) latency.

## Architecture

```mermaid
graph TD
    App[AI Service] --> Log[Logging (slog)]
    App --> Trace[Tracing (OTEL)]
    App --> Met[Metrics (Prometheus)]

    Met --> Persist[Persister (Disk/DB)]
    Met --> Scrape[Prometheus Server]

    Trace --> Export[Jaeger/Tempo]
```
