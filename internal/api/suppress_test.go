package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func resetSuppressStore() {
	suppressMu.Lock()
	defer suppressMu.Unlock()
	suppressStore = map[string]suppressEntry{}
}

func TestHandleSuppressJob_MissingParam(t *testing.T) {
	defer resetSuppressStore()
	req := httptest.NewRequest(http.MethodPost, "/suppress", nil)
	w := httptest.NewRecorder()
	handleSuppressJob(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleSuppressJob_Success(t *testing.T) {
	defer resetSuppressStore()
	body := bytes.NewBufferString(`{"duration_minutes":30}`)
	req := httptest.NewRequest(http.MethodPost, "/suppress?job=backup", body)
	w := httptest.NewRecorder()
	handleSuppressJob(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp suppressEntry
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.JobName != "backup" {
		t.Errorf("expected job_name=backup, got %s", resp.JobName)
	}
	if !IsSuppressed("backup") {
		t.Error("expected job to be suppressed")
	}
}

func TestHandleUnsuppressJob_Success(t *testing.T) {
	defer resetSuppressStore()
	suppressMu.Lock()
	suppressStore["backup"] = suppressEntry{
		JobName: "backup",
		Until:   time.Now().Add(10 * time.Minute),
	}
	suppressMu.Unlock()

	req := httptest.NewRequest(http.MethodPost, "/unsuppress?job=backup", nil)
	w := httptest.NewRecorder()
	handleUnsuppressJob(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if IsSuppressed("backup") {
		t.Error("expected job to no longer be suppressed")
	}
}

func TestHandleSuppressedJobs_List(t *testing.T) {
	defer resetSuppressStore()
	suppressMu.Lock()
	suppressStore["nightly"] = suppressEntry{
		JobName: "nightly",
		Until:   time.Now().Add(5 * time.Minute),
	}
	suppressStore["expired"] = suppressEntry{
		JobName: "expired",
		Until:   time.Now().Add(-1 * time.Minute),
	}
	suppressMu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/suppressed", nil)
	w := httptest.NewRecorder()
	handleSuppressedJobs(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []suppressEntry
	if err := json.NewDecoder(w.Body).Decode(&list); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 active suppression, got %d", len(list))
	}
	if list[0].JobName != "nightly" {
		t.Errorf("expected nightly, got %s", list[0].JobName)
	}
}
