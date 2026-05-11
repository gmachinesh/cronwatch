package api

import (
	"net/http"
	"time"
)

// JobSummary represents a single cron job's current state for the API response.
type JobSummary struct {
	Name        string     `json:"name"`
	Schedule    string     `json:"schedule"`
	LastRun     *time.Time `json:"last_run,omitempty"`
	LastStatus  string     `json:"last_status,omitempty"`
	NextRun     *time.Time `json:"next_run,omitempty"`
	DriftSecs   float64    `json:"drift_seconds,omitempty"`
	Consecutive int        `json:"consecutive_failures,omitempty"`
}

// handleJobs returns a summary of all configured jobs and their current state.
func (s *Server) handleJobs(w http.ResponseWriter, r *http.Request) {
	jobs := s.cfg.Jobs
	summaries := make([]JobSummary, 0, len(jobs))

	for _, j := range jobs {
		sum := JobSummary{
			Name:     j.Name,
			Schedule: j.Schedule,
		}

		if state, ok := s.monitor.State(j.Name); ok {
			sum.LastRun = state.LastRun
			sum.LastStatus = state.LastStatus
			sum.NextRun = state.NextRun
			sum.DriftSecs = state.DriftSeconds
			sum.Consecutive = state.ConsecutiveFailures
		}

		summaries = append(summaries, sum)
	}

	writeJSON(w, http.StatusOK, summaries)
}
