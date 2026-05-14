package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func resetAckStore() {
	ackMu.Lock()
	ackStore = map[string]acknowledgement{}
	ackMu.Unlock()
}

func TestHandleAcknowledgeJob_MissingParam(t *testing.T) {
	defer resetAckStore()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs/acknowledge", nil)
	rec := httptest.NewRecorder()
	handleAcknowledgeJob(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleAcknowledgeJob_Success(t *testing.T) {
	defer resetAckStore()
	body := `{"acked_by":"alice","note":"investigating"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs/acknowledge?job=backup", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	handleAcknowledgeJob(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !IsAcknowledged("backup") {
		t.Fatal("expected job to be acknowledged")
	}
}

func TestHandleAcknowledgeJob_WithExpiry(t *testing.T) {
	defer resetAckStore()
	body := `{"acked_by":"bob","duration":"1ms"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs/acknowledge?job=sync", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	handleAcknowledgeJob(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	time.Sleep(5 * time.Millisecond)
	if IsAcknowledged("sync") {
		t.Fatal("expected acknowledgement to have expired")
	}
}

func TestHandleAcknowledgeJob_InvalidDuration(t *testing.T) {
	defer resetAckStore()
	body := `{"duration":"notaduration"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs/acknowledge?job=sync", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	handleAcknowledgeJob(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleUnacknowledgeJob_Success(t *testing.T) {
	defer resetAckStore()
	ackMu.Lock()
	ackStore["backup"] = acknowledgement{JobName: "backup", AckedAt: time.Now()}
	ackMu.Unlock()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs/unacknowledge?job=backup", nil)
	rec := httptest.NewRecorder()
	handleUnacknowledgeJob(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if IsAcknowledged("backup") {
		t.Fatal("expected job to be unacknowledged")
	}
}

func TestHandleAcknowledgedJobs_List(t *testing.T) {
	defer resetAckStore()
	ackMu.Lock()
	ackStore["job1"] = acknowledgement{JobName: "job1", AckedAt: time.Now()}
	ackStore["job2"] = acknowledgement{JobName: "job2", AckedAt: time.Now()}
	ackMu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/acknowledged", nil)
	rec := httptest.NewRecorder()
	handleAcknowledgedJobs(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var result []acknowledgement
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 acknowledged jobs, got %d", len(result))
	}
}
