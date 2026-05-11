package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

// pauseStore holds the set of paused job names.
type pauseStore struct {
	mu     sync.RWMutex
	paused map[string]bool
}

var globalPauseStore = &pauseStore{
	paused: make(map[string]bool),
}

// IsPaused reports whether the named job is currently paused.
func IsPaused(name string) bool {
	globalPauseStore.mu.RLock()
	defer globalPauseStore.mu.RUnlock()
	return globalPauseStore.paused[name]
}

func (s *Server) handlePauseJob(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("job")
	if name == "" {
		http.Error(w, `{"error":"missing job parameter"}`, http.StatusBadRequest)
		return
	}

	globalPauseStore.mu.Lock()
	globalPauseStore.paused[name] = true
	globalPauseStore.mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]string{"status": "paused", "job": name})
}

func (s *Server) handleResumeJob(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("job")
	if name == "" {
		http.Error(w, `{"error":"missing job parameter"}`, http.StatusBadRequest)
		return
	}

	globalPauseStore.mu.Lock()
	delete(globalPauseStore.paused, name)
	globalPauseStore.mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]string{"status": "resumed", "job": name})
}

func (s *Server) handlePausedJobs(w http.ResponseWriter, r *http.Request) {
	globalPauseStore.mu.RLock()
	list := make([]string, 0, len(globalPauseStore.paused))
	for name := range globalPauseStore.paused {
		list = append(list, name)
	}
	globalPauseStore.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"paused": list})
}
