package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type suppressEntry struct {
	JobName     string    `json:"job_name"`
	SuppressedAt time.Time `json:"suppressed_at"`
	Until       time.Time `json:"until"`
}

var (
	suppressMu    sync.RWMutex
	suppressStore = map[string]suppressEntry{}
)

// IsSuppressed returns true if the job currently has alert suppression active.
func IsSuppressed(jobName string) bool {
	suppressMu.RLock()
	defer suppressMu.RUnlock()
	e, ok := suppressStore[jobName]
	if !ok {
		return false
	}
	return time.Now().Before(e.Until)
}

func handleSuppressJob(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")
	if job == "" {
		http.Error(w, "missing job parameter", http.StatusBadRequest)
		return
	}

	var body struct {
		DurationMinutes int `json:"duration_minutes"`
	}
	body.DurationMinutes = 60 // default
	if err := json.NewDecoder(r.Body).Decode(&body); err == nil && body.DurationMinutes <= 0 {
		body.DurationMinutes = 60
	}

	now := time.Now()
	entry := suppressEntry{
		JobName:      job,
		SuppressedAt: now,
		Until:        now.Add(time.Duration(body.DurationMinutes) * time.Minute),
	}

	suppressMu.Lock()
	suppressStore[job] = entry
	suppressMu.Unlock()

	writeJSON(w, http.StatusOK, entry)
}

func handleUnsuppressJob(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")
	if job == "" {
		http.Error(w, "missing job parameter", http.StatusBadRequest)
		return
	}

	suppressMu.Lock()
	delete(suppressStore, job)
	suppressMu.Unlock()

	writeJSON(w, http.StatusOK, map[string]string{"status": "unsuppressed", "job": job})
}

func handleSuppressedJobs(w http.ResponseWriter, r *http.Request) {
	suppressMu.RLock()
	defer suppressMu.RUnlock()

	now := time.Now()
	active := []suppressEntry{}
	for _, e := range suppressStore {
		if now.Before(e.Until) {
			active = append(active, e)
		}
	}
	writeJSON(w, http.StatusOK, active)
}
