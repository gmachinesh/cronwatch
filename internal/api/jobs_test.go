package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/config"
	"github.com/cronwatch/cronwatch/internal/monitor"
)

func TestHandleJobs_Empty(t *testing.T) {
	srv := makeTestServer(t)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	srv.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var result []JobSummary
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty list, got %d items", len(result))
	}
}

func TestHandleJobs_WithJobs(t *testing.T) {
	srv := makeTestServer(t)
	srv.cfg.Jobs = []config.Job{
		{Name: "backup", Schedule: "@daily", Command: "tar -czf /tmp/b.tar.gz /data"},
		{Name: "cleanup", Schedule: "@hourly", Command: "rm -rf /tmp/old"},
	}

	now := time.Now()
	srv.monitor.RecordSuccess("backup", now, 1200*time.Millisecond)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	srv.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var result []JobSummary
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(result))
	}

	backup := result[0]
	if backup.Name != "backup" {
		t.Errorf("expected name=backup, got %q", backup.Name)
	}
	if backup.LastStatus != monitor.StatusSuccess {
		t.Errorf("expected last_status=success, got %q", backup.LastStatus)
	}
	if backup.LastRun == nil {
		t.Error("expected last_run to be set")
	}

	cleanup := result[1]
	if cleanup.LastStatus != "" {
		t.Errorf("expected empty last_status for unseen job, got %q", cleanup.LastStatus)
	}
}

func TestHandleJobs_MethodNotAllowed(t *testing.T) {
	srv := makeTestServer(t)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/jobs", nil)
	srv.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
