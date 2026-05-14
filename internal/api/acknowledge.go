package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type acknowledgement struct {
	JobName    string    `json:"job_name"`
	AckedAt    time.Time `json:"acked_at"`
	AckedBy    string    `json:"acked_by"`
	Note       string    `json:"note"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
}

var (
	ackMu    sync.RWMutex
	ackStore = map[string]acknowledgement{}
)

// IsAcknowledged returns true if the given job has an active acknowledgement.
func IsAcknowledged(jobName string) bool {
	ackMu.RLock()
	defer ackMu.RUnlock()
	ack, ok := ackStore[jobName]
	if !ok {
		return false
	}
	if !ack.ExpiresAt.IsZero() && time.Now().After(ack.ExpiresAt) {
		return false
	}
	return true
}

func handleAcknowledgeJob(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")
	if job == "" {
		http.Error(w, "missing job parameter", http.StatusBadRequest)
		return
	}

	var body struct {
		AckedBy  string `json:"acked_by"`
		Note     string `json:"note"`
		Duration string `json:"duration"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	ack := acknowledgement{
		JobName: job,
		AckedAt: time.Now(),
		AckedBy: body.AckedBy,
		Note:    body.Note,
	}
	if body.Duration != "" {
		d, err := time.ParseDuration(body.Duration)
		if err != nil {
			http.Error(w, "invalid duration: "+err.Error(), http.StatusBadRequest)
			return
		}
		ack.ExpiresAt = time.Now().Add(d)
	}

	ackMu.Lock()
	ackStore[job] = ack
	ackMu.Unlock()

	writeJSON(w, http.StatusOK, map[string]any{"acknowledged": job, "acked_at": ack.AckedAt})
}

func handleUnacknowledgeJob(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")
	if job == "" {
		http.Error(w, "missing job parameter", http.StatusBadRequest)
		return
	}

	ackMu.Lock()
	delete(ackStore, job)
	ackMu.Unlock()

	writeJSON(w, http.StatusOK, map[string]string{"unacknowledged": job})
}

func handleAcknowledgedJobs(w http.ResponseWriter, r *http.Request) {
	ackMu.RLock()
	defer ackMu.RUnlock()

	now := time.Now()
	active := []acknowledgement{}
	for _, ack := range ackStore {
		if ack.ExpiresAt.IsZero() || now.Before(ack.ExpiresAt) {
			active = append(active, ack)
		}
	}
	writeJSON(w, http.StatusOK, active)
}
