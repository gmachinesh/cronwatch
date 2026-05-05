package api

import (
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

type statusEntry struct {
	Job       string `json:"job"`
	LastRunAt string `json:"last_run_at,omitempty"`
	LastResult string `json:"last_result,omitempty"`
	Drifted   bool   `json:"drifted"`
}

type historyRequest struct {
	Job   string
	Limit int
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	states := s.mon.States()
	var entries []statusEntry
	for job, st := range states {
		e := statusEntry{
			Job:     job,
			Drifted: st.Drifted,
		}
		if !st.LastRun.IsZero() {
			e.LastRunAt = st.LastRun.Format("2006-01-02T15:04:05Z07:00")
		}
		if st.LastError != nil {
			e.LastResult = st.LastError.Error()
		} else if !st.LastRun.IsZero() {
			e.LastResult = "success"
		}
		entries = append(entries, e)
	}
	writeJSON(w, http.StatusOK, entries)
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	job := r.URL.Query().Get("job")
	records, err := s.hist.Recent(job, 50)
	if err != nil {
		http.Error(w, "failed to read history", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, records)
}
