package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetMetrics() {
	globalMetrics.totalRequests.Store(0)
	globalMetrics.totalErrors.Store(0)
	globalMetrics.totalDurationMs.Store(0)
}

func TestMetricsMiddleware_CountsRequests(t *testing.T) {
	resetMetrics()
	handler := metricsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	for i := 0; i < 3; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	}
	snap := snapshotMetrics()
	if snap.TotalRequests != 3 {
		t.Errorf("expected 3 requests, got %d", snap.TotalRequests)
	}
	if snap.TotalErrors != 0 {
		t.Errorf("expected 0 errors, got %d", snap.TotalErrors)
	}
}

func TestMetricsMiddleware_CountsErrors(t *testing.T) {
	resetMetrics()
	handler := metricsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	snap := snapshotMetrics()
	if snap.TotalErrors != 1 {
		t.Errorf("expected 1 error, got %d", snap.TotalErrors)
	}
}

func TestMetricsMiddleware_AvgDuration(t *testing.T) {
	resetMetrics()
	handler := metricsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	snap := snapshotMetrics()
	if snap.TotalRequests != 1 {
		t.Errorf("expected 1 request, got %d", snap.TotalRequests)
	}
	// avg duration should be >= 0
	if snap.AvgDurationMs < 0 {
		t.Errorf("expected non-negative avg duration, got %d", snap.AvgDurationMs)
	}
}

func TestSnapshotMetrics_ZeroOnInit(t *testing.T) {
	resetMetrics()
	snap := snapshotMetrics()
	if snap.TotalRequests != 0 || snap.TotalErrors != 0 || snap.AvgDurationMs != 0 {
		t.Errorf("expected zero snapshot, got %+v", snap)
	}
}
