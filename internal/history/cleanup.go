package history

import (
	"time"
)

// CleanupOptions configures the cleanup behaviour.
type CleanupOptions struct {
	// MaxAge removes entries older than this duration. Zero means no age limit.
	MaxAge time.Duration
	// MaxEntries keeps at most this many entries (newest first). Zero means no limit.
	MaxEntries int
}

// Cleanup removes old or excess entries from the history file and rewrites it.
// It applies MaxAge first, then MaxEntries.
func (h *History) Cleanup(opts CleanupOptions) (removed int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	entries, err := h.readAll()
	if err != nil {
		return 0, err
	}

	orig := len(entries)

	if opts.MaxAge > 0 {
		cutoff := time.Now().Add(-opts.MaxAge)
		filtered := entries[:0]
		for _, e := range entries {
			if e.StartedAt.After(cutoff) {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	if opts.MaxEntries > 0 && len(entries) > opts.MaxEntries {
		// entries are stored oldest-first; keep the newest N
		entries = entries[len(entries)-opts.MaxEntries:]
	}

	if err := h.writeAll(entries); err != nil {
		return 0, err
	}

	return orig - len(entries), nil
}
