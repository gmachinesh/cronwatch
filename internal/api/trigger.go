package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/user/cronwatch/internal/scheduler"
)

type triggerResponse struct {
	Job       string    `json:"job"`
	Triggered bool      `json:"triggered"`
	At        time.Time `json:"at"`
	Message   string    `json:"message,omitempty"`
}

// handleTriggerJob manually triggers a cron job by name outside its schedule.
func (s *Server) handleTriggerJob(w http.ResponseWriter, r *http.Request) {
	jobName := r.URL.Query().Get("job")
	if jobName == "" {
		http.Error(w, "missing required query parameter: job", http.StatusBadRequest)
		return
	}

	if IsPaused(jobName) {
		writeJSON(w, http.StatusConflict, triggerResponse{
			Job:       jobName,
			Triggered: false,
			At:        time.Now().UTC(),
			Message:   "job is paused and cannot be triggered",
		})
		return
	}

	if err := s.scheduler.TriggerNow(jobName); err != nil {
		if err == scheduler.ErrJobNotFound {
			http.Error(w, "job not found: "+jobName, http.StatusNotFound)
			return
		}
		http.Error(w, "failed to trigger job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, triggerResponse{
		Job:       jobName,
		Triggered: true,
		At:        time.Now().UTC(),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
