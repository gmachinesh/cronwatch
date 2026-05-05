package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestNew_EmptyFile(t *testing.T) {
	s, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Recent("job1", 10); len(got) != 0 {
		t.Errorf("expected 0 entries, got %d", len(got))
	}
}

func TestAppend_AndRecent(t *testing.T) {
	path := tempPath(t)
	s, _ := New(path)

	entry := Entry{
		JobName:   "backup",
		StartedAt: time.Now(),
		Duration:  2 * time.Second,
		Success:   true,
	}
	if err := s.Append(entry); err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	results := s.Recent("backup", 5)
	if len(results) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(results))
	}
	if results[0].JobName != "backup" {
		t.Errorf("unexpected job name: %s", results[0].JobName)
	}
}

func TestRecent_LimitsResults(t *testing.T) {
	s, _ := New(tempPath(t))
	for i := 0; i < 10; i++ {
		_ = s.Append(Entry{JobName: "job", Success: true, StartedAt: time.Now()})
	}
	results := s.Recent("job", 3)
	if len(results) != 3 {
		t.Errorf("expected 3 entries, got %d", len(results))
	}
}

func TestRecent_FiltersByJobName(t *testing.T) {
	s, _ := New(tempPath(t))
	_ = s.Append(Entry{JobName: "jobA", Success: true, StartedAt: time.Now()})
	_ = s.Append(Entry{JobName: "jobB", Success: false, StartedAt: time.Now()})

	results := s.Recent("jobA", 10)
	if len(results) != 1 || results[0].JobName != "jobA" {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestPersistence_ReloadFromDisk(t *testing.T) {
	path := tempPath(t)
	s1, _ := New(path)
	_ = s1.Append(Entry{JobName: "persist", Success: true, StartedAt: time.Now()})

	s2, err := New(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	results := s2.Recent("persist", 5)
	if len(results) != 1 {
		t.Errorf("expected 1 persisted entry, got %d", len(results))
	}
}

func TestNew_InvalidJSON(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := New(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
