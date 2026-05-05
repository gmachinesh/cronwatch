package api

import "net/http"

// registerRoutes wires all API endpoints onto the given mux, wrapped with
// the logging middleware and per-route method guards.
func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.Handle("/healthz", loggingMiddleware(
		methodGuard(http.MethodGet, s.handleHealthz),
	))

	mux.Handle("/status", loggingMiddleware(
		methodGuard(http.MethodGet, s.handleStatus),
	))

	mux.Handle("/history", loggingMiddleware(
		methodGuard(http.MethodGet, s.handleHistory),
	))
}
