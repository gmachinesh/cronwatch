package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type silenceEntry struct {
	JobName  string    `json:"job_name"`
	Until    time.Time `json:"until"`
	Reason   string    `json:"reason,omitempty"`
}

var (
	silenceMu    sync.RWMutex
	silenceStore = map[string]silenceEntry{}
)

// IsSilenced returns true if the given job has an active silence window.
func IsSilenced(jobName string) bool {
	silenceMu.RLock()
	defer silenceMu.RUnlock()
	e, ok := silenceStore[jobName]
	if !ok {
		return false
	}
	if time.Now().After(e.Until) {
		return false
	}
	return true
}

func handleSilenceJob(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")
	if job == "" {
		http.Error(w, "missing job parameter", http.StatusBadRequest)
		return
	}

	var body struct {
		DurationMinutes int    `json:"duration_minutes"`
		Reason          string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.DurationMinutes <= 0 {
		http.Error(w, "duration_minutes must be a positive integer", http.StatusBadRequest)
		return
	}

	until := time.Now().Add(time.Duration(body.DurationMinutes) * time.Minute)
	silenceMu.Lock()
	silenceStore[job] = silenceEntry{JobName: job, Until: until, Reason: body.Reason}
	silenceMu.Unlock()

	writeJSON(w, http.StatusOK, map[string]interface{}{"job": job, "silenced_until": until})
}

func handleUnsilenceJob(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")
	if job == "" {
		http.Error(w, "missing job parameter", http.StatusBadRequest)
		return
	}
	silenceMu.Lock()
	delete(silenceStore, job)
	silenceMu.Unlock()
	writeJSON(w, http.StatusOK, map[string]string{"job": job, "status": "unsilenced"})
}

func handleSilencedJobs(w http.ResponseWriter, r *http.Request) {
	silenceMu.RLock()
	defer silenceMu.RUnlock()
	now := time.Now()
	active := []silenceEntry{}
	for _, e := range silenceStore {
		if now.Before(e.Until) {
			active = append(active, e)
		}
	}
	writeJSON(w, http.StatusOK, active)
}
