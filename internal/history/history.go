package history

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry represents a single cron job execution record.
type Entry struct {
	JobName   string        `json:"job_name"`
	StartedAt time.Time     `json:"started_at"`
	Duration  time.Duration `json:"duration_ns"`
	Success   bool          `json:"success"`
	Output    string        `json:"output,omitempty"`
	Error     string        `json:"error,omitempty"`
}

// History manages persistent job execution history stored as newline-delimited JSON.
type History struct {
	mu   sync.Mutex
	path string
}

// New opens (or creates) a history file at the given path.
func New(path string) (*History, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("history: open %s: %w", path, err)
	}
	f.Close()
	return &History{path: path}, nil
}

// Append writes a new entry to the history file.
func (h *History) Append(e Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	f, err := os.OpenFile(h.path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("history: append open: %w", err)
	}
	defer f.Close()

	line, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}

// Recent returns the last n entries, optionally filtered by jobName (empty = all jobs).
func (h *History) Recent(n int, jobName string) ([]Entry, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	all, err := h.readAll()
	if err != nil {
		return nil, err
	}

	var filtered []Entry
	for _, e := range all {
		if jobName == "" || e.JobName == jobName {
			filtered = append(filtered, e)
		}
	}

	if n > 0 && len(filtered) > n {
		filtered = filtered[len(filtered)-n:]
	}
	return filtered, nil
}

// Prune removes entries older than the given age from the history file.
// It rewrites the file in place, keeping only entries whose StartedAt
// timestamp is within maxAge of the current time.
func (h *History) Prune(maxAge time.Duration) (int, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	all, err := h.readAll()
	if err != nil {
		return 0, err
	}

	cutoff := time.Now().Add(-maxAge)
	var keep []Entry
	for _, e := range all {
		if e.StartedAt.After(cutoff) {
			keep = append(keep, e)
		}
	}

	removed := len(all) - len(keep)
	if removed == 0 {
		return 0, nil
	}
	return removed, h.writeAll(keep)
}

// readAll reads every entry from the file. Caller must hold h.mu.
func (h *History) readAll() ([]Entry, error) {
	data, err := os.ReadFile(h.path)
	if err != nil {
		return nil, fmt.Errorf("history: read: %w", err)
	}
	var entries []Entry
	for _, raw := range splitLines(data) {
		if len(raw) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(raw, &e); err != nil {
			continue // skip malformed lines
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// writeAll rewrites the entire history file with the given entries.
func (h *History) writeAll(entries []Entry) error {
	f, err := os.OpenFile(h.path, os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("history: writeAll open: %w", err)
	}
	defer f.Close()
	for _, e := range entries {
		line, err := json.Marshal(e)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(f, "%s\n", line); err != nil {
			return err
		}
	}
	return nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
