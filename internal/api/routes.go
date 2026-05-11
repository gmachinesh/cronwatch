package api

import "net/http"

// registerRoutes wires all HTTP handlers onto the server's mux.
func (s *Server) registerRoutes() {
	mux := s.mux

	mux.Handle("/healthz", loggingMiddleware(metricsMiddleware(methodGuard(http.MethodGet, http.HandlerFunc(s.handleHealthz)))))
	mux.Handle("/api/status", loggingMiddleware(metricsMiddleware(apiKeyAuth(s.cfg, methodGuard(http.MethodGet, http.HandlerFunc(s.handleStatus))))))
	mux.Handle("/api/history", loggingMiddleware(metricsMiddleware(apiKeyAuth(s.cfg, methodGuard(http.MethodGet, http.HandlerFunc(s.handleHistory))))))
	mux.Handle("/api/jobs", loggingMiddleware(metricsMiddleware(apiKeyAuth(s.cfg, methodGuard(http.MethodGet, http.HandlerFunc(s.handleJobs))))))
	mux.Handle("/api/version", loggingMiddleware(metricsMiddleware(methodGuard(http.MethodGet, http.HandlerFunc(s.handleVersion)))))
	mux.Handle("/api/metrics", loggingMiddleware(metricsMiddleware(methodGuard(http.MethodGet, http.HandlerFunc(s.handleMetrics)))))
}
