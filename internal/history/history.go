package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single cron job execution record.
type Entry struct {
	JobName   string        `json:"job_name"`
	StartedAt time.Time     `json:"started_at"`
	Duration  time.Duration `json:"duration"`
	Success   bool          `json:"success"`
	Output    string        `json:"output,omitempty"`
	Error     string        `json:"error,omitempty"`
}

// Store persists job execution history to a JSON file.
type Store struct {
	mu      sync.Mutex
	path    string
	entries []Entry
}

// New creates a new Store backed by the given file path.
// Existing entries are loaded if the file exists.
func New(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Append adds a new entry and flushes to disk.
func (s *Store) Append(e Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, e)
	return s.flush()
}

// Recent returns the last n entries for the given job name.
func (s *Store) Recent(jobName string, n int) []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []Entry
	for i := len(s.entries) - 1; i >= 0 && len(result) < n; i-- {
		if s.entries[i].JobName == jobName {
			result = append([]Entry{s.entries[i]}, result...)
		}
	}
	return result
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.entries)
}

func (s *Store) flush() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
