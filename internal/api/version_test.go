package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
)

func TestHandleVersion_ReturnsOK(t *testing.T) {
	srv := makeTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	srv.handleVersion(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandleVersion_Fields(t *testing.T) {
	Version = "1.2.3"
	Commit = "abc1234"
	BuiltAt = "2024-01-01T00:00:00Z"
	t.Cleanup(func() {
		Version = "dev"
		Commit = "none"
		BuiltAt = ""
	})

	srv := makeTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	srv.handleVersion(rec, req)

	var info BuildInfo
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if info.Version != "1.2.3" {
		t.Errorf("expected version 1.2.3, got %s", info.Version)
	}
	if info.Commit != "abc1234" {
		t.Errorf("expected commit abc1234, got %s", info.Commit)
	}
	if info.BuiltAt != "2024-01-01T00:00:00Z" {
		t.Errorf("expected built_at 2024-01-01T00:00:00Z, got %s", info.BuiltAt)
	}
	if info.GoVersion != runtime.Version() {
		t.Errorf("expected go version %s, got %s", runtime.Version(), info.GoVersion)
	}
	if info.OS != runtime.GOOS {
		t.Errorf("expected os %s, got %s", runtime.GOOS, info.OS)
	}
	if info.Arch != runtime.GOARCH {
		t.Errorf("expected arch %s, got %s", runtime.GOARCH, info.Arch)
	}
}

func TestHandleVersion_DefaultBuiltAt(t *testing.T) {
	BuiltAt = ""
	t.Cleanup(func() { BuiltAt = "" })

	srv := makeTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	srv.handleVersion(rec, req)

	var info BuildInfo
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if info.BuiltAt == "" {
		t.Error("expected non-empty built_at when BuiltAt var is unset")
	}
}
