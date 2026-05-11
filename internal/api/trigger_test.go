package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/cronwatch/internal/scheduler"
)

// mockScheduler satisfies the scheduler interface used by Server.
type mockScheduler struct {
	triggerErr error
	triggered  []string
}

func (m *mockScheduler) TriggerNow(job string) error {
	if m.triggerErr != nil {
		return m.triggerErr
	}
	m.triggered = append(m.triggered, job)
	return nil
}

func TestHandleTriggerJob_MissingParam(t *testing.T) {
	s := makeTestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/api/trigger", nil)
	rec := httptest.NewRecorder()
	s.handleTriggerJob(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleTriggerJob_JobNotFound(t *testing.T) {
	s := makeTestServer(t)
	s.scheduler = &mockScheduler{triggerErr: scheduler.ErrJobNotFound}
	req := httptest.NewRequest(http.MethodPost, "/api/trigger?job=missing", nil)
	rec := httptest.NewRecorder()
	s.handleTriggerJob(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandleTriggerJob_Paused(t *testing.T) {
	resetPauseStore()
	s := makeTestServer(t)
	pauseStore["backup"] = true
	req := httptest.NewRequest(http.MethodPost, "/api/trigger?job=backup", nil)
	rec := httptest.NewRecorder()
	s.handleTriggerJob(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", rec.Code)
	}
	var resp triggerResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Triggered {
		t.Error("expected Triggered=false for paused job")
	}
}

func TestHandleTriggerJob_Success(t *testing.T) {
	resetPauseStore()
	ms := &mockScheduler{}
	s := makeTestServer(t)
	s.scheduler = ms
	req := httptest.NewRequest(http.MethodPost, "/api/trigger?job=cleanup", nil)
	rec := httptest.NewRecorder()
	s.handleTriggerJob(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp triggerResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if !resp.Triggered {
		t.Error("expected Triggered=true")
	}
	if resp.Job != "cleanup" {
		t.Errorf("expected job=cleanup, got %s", resp.Job)
	}
}

func TestHandleTriggerJob_InternalError(t *testing.T) {
	resetPauseStore()
	s := makeTestServer(t)
	s.scheduler = &mockScheduler{triggerErr: errors.New("exec pool full")}
	req := httptest.NewRequest(http.MethodPost, "/api/trigger?job=sync", nil)
	rec := httptest.NewRecorder()
	s.handleTriggerJob(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}
