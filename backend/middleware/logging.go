package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap ResponseWriter to capture status code
		lrw := &loggingResponseWriter{w, http.StatusOK}
		
		next.ServeHTTP(lrw, r)
		
		latency := time.Since(start)
		
		// Use Trace Context if available (X-Cloud-Trace-Context)
		traceID := r.Header.Get("X-Cloud-Trace-Context")
		
		args := []any{
			slog.String("http_method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", lrw.statusCode),
			slog.Int64("latency_ms", latency.Milliseconds()),
			slog.String("ip", getIP(r)), // Shared within package middleware
			slog.String("user_agent", r.UserAgent()),
		}

		if traceID != "" {
			args = append(args, slog.String("trace_id", traceID))
		}
		
		slog.Info("Request handled", args...)
	})
}
