package api

import "net/http"

// registerPauseRoutes wires up the job pause/resume endpoints onto mux.
// All mutating routes require POST; the listing route accepts GET.
func (s *Server) registerPauseRoutes(mux *http.ServeMux) {
	mux.Handle("/jobs/pause",
		s.authMiddleware(
			loggingMiddleware(
				metricsMiddleware(
					methodGuard(http.MethodPost,
						http.HandlerFunc(s.handlePauseJob),
					),
				),
			),
		),
	)

	mux.Handle("/jobs/resume",
		s.authMiddleware(
			loggingMiddleware(
				metricsMiddleware(
					methodGuard(http.MethodPost,
						http.HandlerFunc(s.handleResumeJob),
					),
				),
			),
		),
	)

	mux.Handle("/jobs/paused",
		s.authMiddleware(
			loggingMiddleware(
				metricsMiddleware(
					methodGuard(http.MethodGet,
						http.HandlerFunc(s.handlePausedJobs),
					),
				),
			),
		),
	)
}
