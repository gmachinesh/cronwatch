package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func resetSilenceStore() {
	silenceMu.Lock()
	silenceStore = map[string]silenceEntry{}
	silenceMu.Unlock()
}

func TestHandleSilenceJob_MissingParam(t *testing.T) {
	defer resetSilenceStore()
	body := `{"duration_minutes":10}`
	req := httptest.NewRequest(http.MethodPost, "/silence", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handleSilenceJob(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleSilenceJob_InvalidDuration(t *testing.T) {
	defer resetSilenceStore()
	body := `{"duration_minutes":0}`
	req := httptest.NewRequest(http.MethodPost, "/silence?job=backup", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handleSilenceJob(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleSilenceJob_Success(t *testing.T) {
	defer resetSilenceStore()
	body := `{"duration_minutes":30,"reason":"maintenance"}`
	req := httptest.NewRequest(http.MethodPost, "/silence?job=backup", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handleSilenceJob(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !IsSilenced("backup") {
		t.Fatal("expected job to be silenced")
	}
}

func TestHandleUnsilenceJob_Success(t *testing.T) {
	defer resetSilenceStore()
	silenceMu.Lock()
	silenceStore["backup"] = silenceEntry{JobName: "backup", Until: time.Now().Add(time.Hour)}
	silenceMu.Unlock()

	req := httptest.NewRequest(http.MethodPost, "/unsilence?job=backup", nil)
	w := httptest.NewRecorder()
	handleUnsilenceJob(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if IsSilenced("backup") {
		t.Fatal("expected job to be unsilenced")
	}
}

func TestHandleSilencedJobs_List(t *testing.T) {
	defer resetSilenceStore()
	silenceMu.Lock()
	silenceStore["nightly"] = silenceEntry{JobName: "nightly", Until: time.Now().Add(time.Hour)}
	silenceStore["expired"] = silenceEntry{JobName: "expired", Until: time.Now().Add(-time.Minute)}
	silenceMu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/silenced", nil)
	w := httptest.NewRecorder()
	handleSilencedJobs(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result []silenceEntry
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(result) != 1 || result[0].JobName != "nightly" {
		t.Fatalf("expected only active silence, got %+v", result)
	}
}

func TestIsSilenced_Expired(t *testing.T) {
	defer resetSilenceStore()
	silenceMu.Lock()
	silenceStore["old"] = silenceEntry{JobName: "old", Until: time.Now().Add(-time.Second)}
	silenceMu.Unlock()
	if IsSilenced("old") {
		t.Fatal("expected expired silence to return false")
	}
}
