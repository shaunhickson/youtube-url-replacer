# Design: Structured Observability

## Overview
Replace unstructured text logs with JSON-structured logs compatible with Google Cloud Run (and Cloud Logging). Use the standard library `log/slog`.

## Goals
- **JSON Output:** Machine-readable logs.
- **Severity Levels:** `DEBUG`, `INFO`, `WARN`, `ERROR` mapped to GCP `severity` field.
- **Context:** Include `trace_id` (from Cloud Trace headers) and `component` fields.
- **Request Logging:** Middleware to log every HTTP request with status, latency, and path.

## Implementation Details

### 1. Logger Configuration
Create `backend/logger/logger.go`:
- Use `slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))`.
- Add a hook or `ReplaceAttr` function to map `level` key to `severity` (GCP convention).
- Map `msg` to `message` (GCP convention).

### 2. Middleware
Create `backend/middleware/logging.go`:
- Wrap `http.Handler`.
- Extract `X-Cloud-Trace-Context` header.
- Log at start (Debug) and end (Info) of request.
- Fields: `http_method`, `path`, `status`, `latency`, `ip`, `user_agent`.

### 3. Application Updates
- Replace `log.Printf` with `slog.Info`, `slog.Error`.
- Ensure `main.go` initializes the logger early.

## Verification
- Run locally and verify JSON output.
- Check that error logs include stack traces or error details.
