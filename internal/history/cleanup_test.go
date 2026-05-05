package history

import (
	"testing"
	"time"
)

func TestCleanup_MaxAge(t *testing.T) {
	h, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	now := time.Now()
	old := Entry{JobName: "old-job", StartedAt: now.Add(-48 * time.Hour), Success: true}
	recent := Entry{JobName: "new-job", StartedAt: now.Add(-1 * time.Hour), Success: true}

	for _, e := range []Entry{old, recent} {
		if err := h.Append(e); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}

	removed, err := h.Cleanup(CleanupOptions{MaxAge: 24 * time.Hour})
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}

	entries, err := h.Recent(10, "")
	if err != nil {
		t.Fatalf("Recent: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].JobName != "new-job" {
		t.Errorf("expected new-job, got %s", entries[0].JobName)
	}
}

func TestCleanup_MaxEntries(t *testing.T) {
	h, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	base := time.Now().Add(-10 * time.Minute)
	for i := 0; i < 5; i++ {
		e := Entry{
			JobName:   "job",
			StartedAt: base.Add(time.Duration(i) * time.Minute),
			Success:   true,
		}
		if err := h.Append(e); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}

	removed, err := h.Cleanup(CleanupOptions{MaxEntries: 3})
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}

	entries, err := h.Recent(10, "")
	if err != nil {
		t.Fatalf("Recent: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestCleanup_NoOptions(t *testing.T) {
	h, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for i := 0; i < 3; i++ {
		_ = h.Append(Entry{JobName: "job", StartedAt: time.Now(), Success: true})
	}
	removed, err := h.Cleanup(CleanupOptions{})
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if removed != 0 {
		t.Errorf("expected 0 removed with no options, got %d", removed)
	}
}
