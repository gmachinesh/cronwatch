package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cronwatch/internal/history"
	"github.com/cronwatch/internal/monitor"
)

// Server exposes a lightweight HTTP API for cronwatch status.
type Server struct {
	addr    string
	mon     *monitor.Monitor
	hist    *history.History
	httpSrv *http.Server
}

// New creates a new API Server.
func New(addr string, mon *monitor.Monitor, hist *history.History) *Server {
	s := &Server{addr: addr, mon: mon, hist: hist}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/status", s.handleStatus)
	mux.HandleFunc("/history", s.handleHistory)
	s.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s
}

// Start begins listening. Blocks until the server stops.
func (s *Server) Start() error {
	return s.httpSrv.ListenAndServe()
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop() error {
	return s.httpSrv.Close()
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
