package api

import (
	"net/http"
	"sync/atomic"
	"time"
)

// requestMetrics tracks basic HTTP request counters for the API server.
type requestMetrics struct {
	totalRequests  atomic.Int64
	totalErrors    atomic.Int64
	totalDurationMs atomic.Int64
}

var globalMetrics = &requestMetrics{}

// metricsMiddleware records request count, error count, and cumulative duration.
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		elapsed := time.Since(start).Milliseconds()
		globalMetrics.totalRequests.Add(1)
		globalMetrics.totalDurationMs.Add(elapsed)
		if rw.status >= 500 {
			globalMetrics.totalErrors.Add(1)
		}
	})
}

// statusRecorder wraps ResponseWriter to capture the HTTP status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// MetricsSnapshot holds a point-in-time snapshot of API metrics.
type MetricsSnapshot struct {
	TotalRequests   int64 `json:"total_requests"`
	TotalErrors     int64 `json:"total_errors"`
	AvgDurationMs   int64 `json:"avg_duration_ms"`
}

func snapshotMetrics() MetricsSnapshot {
	reqs := globalMetrics.totalRequests.Load()
	errs := globalMetrics.totalErrors.Load()
	dur := globalMetrics.totalDurationMs.Load()
	var avg int64
	if reqs > 0 {
		avg = dur / reqs
	}
	return MetricsSnapshot{
		TotalRequests: reqs,
		TotalErrors:   errs,
		AvgDurationMs: avg,
	}
}
