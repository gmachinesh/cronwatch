package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetPauseStore() {
	globalPauseStore.mu.Lock()
	globalPauseStore.paused = make(map[string]bool)
	globalPauseStore.mu.Unlock()
}

func TestHandlePauseJob_MissingParam(t *testing.T) {
	defer resetPauseStore()
	s := makeTestServer(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/jobs/pause", nil)
	s.handlePauseJob(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestHandlePauseJob_Success(t *testing.T) {
	defer resetPauseStore()
	s := makeTestServer(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/jobs/pause?job=backup", nil)
	s.handlePauseJob(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if !IsPaused("backup") {
		t.Fatal("expected backup to be paused")
	}
}

func TestHandleResumeJob_Success(t *testing.T) {
	defer resetPauseStore()
	globalPauseStore.paused["backup"] = true
	s := makeTestServer(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/jobs/resume?job=backup", nil)
	s.handleResumeJob(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if IsPaused("backup") {
		t.Fatal("expected backup to be resumed")
	}
}

func TestHandlePausedJobs_List(t *testing.T) {
	defer resetPauseStore()
	globalPauseStore.paused["cleanup"] = true
	globalPauseStore.paused["report"] = true
	s := makeTestServer(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/jobs/paused", nil)
	s.handlePausedJobs(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string][]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(resp["paused"]) != 2 {
		t.Fatalf("expected 2 paused jobs, got %d", len(resp["paused"]))
	}
}

func TestIsPaused_FalseByDefault(t *testing.T) {
	defer resetPauseStore()
	if IsPaused("nonexistent") {
		t.Fatal("expected nonexistent job to not be paused")
	}
}
